package object_test_wrapper

import (
	"errors"
	"fmt"

	"github.com/khevse/codegen/internal/pkg/astpkg"
	"github.com/samber/lo"
)

type field struct {
	ObjectSpecName string
	FuncSpecName   string
	TypeName       string
	MockPackage    string
	MockTypeName   string
	Type           astpkg.Type
}

type methodSpec struct {
	Name    string
	Comment string
	Params  []field
	Results []field
}

type objectSpecField struct {
	Name         string
	TypeName     string
	MockPackage  string
	MockTypeName string
	Type         astpkg.Type
}

type objectSpec struct {
	Name               string
	Comment            string
	Methods            []methodSpec
	Fields             []objectSpecField
	BaseObjectTypeName string
}

func newObjectSpec(
	interfaceType argInterfaceType,
	objectType argObjectType,
	mockPackageName string,
	typeDecl *astpkg.TypeDecl,
	imports astpkg.ImportList,
) (*objectSpec, error) {
	castedType, ok := astpkg.CastToType[astpkg.InterfaceType](typeDecl.Type)
	if !ok {
		return nil, errors.New("type is not interface")
	}

	for _, item := range castedType.Methods {
		err := astpkg.InspectType(item.Type, func(t astpkg.Type) error {
			return astpkg.ReplaceImportAliasByImportPath(t, imports)
		})
		if err != nil {
			return nil, fmt.Errorf("replace method imports(%s): %T", item.Name, item.Type)
		}
	}

	objectPackage, ok := imports.GetByPath(objectType.Package)
	if !ok {
		return nil, fmt.Errorf("get object type package: %s", objectType.Package)
	}

	var baseObjectTypeName string
	if objectPackage.Alias == "" {
		baseObjectTypeName = fmt.Sprintf("*%s", objectType.TypeName)
	} else {
		baseObjectTypeName = fmt.Sprintf("*%s.%s", objectPackage.Alias, objectType.TypeName)
	}

	methodList := make([]methodSpec, 0, len(castedType.Methods))
	objectSpecFieldList := make([]objectSpecField, 0)

	for _, item := range castedType.Methods {
		if !astpkg.IsExported(item.Name) {
			continue
		}

		casedMethod, ok := item.Type.(*astpkg.FuncType)
		if !ok {
			return nil, fmt.Errorf("cast method type(%s): %T", item.Name, item.Type)
		}

		params, err := newFieldsList(item.Name, casedMethod.Params, mockPackageName, imports)
		if err != nil {
			return nil, fmt.Errorf("fields list from params: %w", err)
		}

		results, err := newFieldsList(item.Name, casedMethod.Results, mockPackageName, imports)
		if err != nil {
			return nil, fmt.Errorf("fields list from results: %w", err)
		}

		method := methodSpec{
			Name:    item.Name,
			Comment: fmt.Sprintf("%s .", item.Name),
			Params:  params,
			Results: results,
		}

		methodList = append(methodList, method)

		for _, item := range results {
			objectSpecFieldList = append(objectSpecFieldList, newObjectSpecField(item))
		}
	}

	objectSpecComment := fmt.Sprintf(
		"%s wrapper for type %s: %s",
		interfaceType.WrapperName, typeDecl.Name, typeDecl.Comment,
	)

	return &objectSpec{
		Name:               interfaceType.WrapperName,
		Comment:            objectSpecComment,
		Methods:            methodList,
		Fields:             objectSpecFieldList,
		BaseObjectTypeName: baseObjectTypeName,
	}, nil
}

func newFieldsList(
	methodName string,
	filedList []*astpkg.Field,
	mockPackageName string,
	imports astpkg.ImportList,
) ([]field, error) {
	imp, ok := imports.GetByPath(mockPackageName)
	if !ok {
		return nil, fmt.Errorf("not found mock package: %s", mockPackageName)
	}

	var mockPackageAlias string
	if imp.Path != "" {
		mockPackageAlias = imp.Alias
	}

	fieldList := make([]field, 0, len(filedList))
	for i, item := range filedList {
		objectSpecName := fmt.Sprintf("%sArg%d", methodName, i)
		mockTypeName := ""
		mockPackage := ""
		if castedType, ok := astpkg.CastToType[astpkg.Ident](item.Type); ok {
			if _, ok := astpkg.CastToType[astpkg.InterfaceType](castedType.Type); ok {
				objectSpecName = castedType.Name
				mockTypeName = castedType.Name + "Mock"
				mockPackage = mockPackageAlias
			}
		} else if castedType, ok := astpkg.CastToType[astpkg.SelectorExpr](item.Type); ok {
			if _, ok := astpkg.CastToType[astpkg.InterfaceType](castedType.Type); ok {
				objectSpecName = castedType.Name
				mockTypeName = castedType.Name + "Mock"
				mockPackage = mockPackageAlias
			}
		}

		f := field{
			ObjectSpecName: objectSpecName,
			MockTypeName:   mockTypeName,
			MockPackage:    mockPackage,
			FuncSpecName:   lo.Ternary(item.Name == "", "_", item.Name),
			TypeName:       item.Type.ExprString(),
			Type:           item.Type,
		}

		fieldList = append(fieldList, f)
	}

	return fieldList, nil
}

func newObjectSpecField(fieldDesc field) objectSpecField {
	return objectSpecField{
		Name:         fieldDesc.ObjectSpecName,
		TypeName:     fieldDesc.TypeName,
		MockPackage:  fieldDesc.MockPackage,
		MockTypeName: fieldDesc.MockTypeName,
		Type:         fieldDesc.Type,
	}
}

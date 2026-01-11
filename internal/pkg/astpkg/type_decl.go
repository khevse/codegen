package astpkg

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/samber/lo"
)

type TypeDecl struct {
	Package     string
	PackagePath string
	Name        string
	Comment     string
	Type        Type
}

func (t TypeDecl) String() string {
	if t.PackagePath == "" {
		return fmt.Sprintf("%s(%s)", t.Name, t.Type)
	}

	return fmt.Sprintf("%s.%s(%s)", t.PackagePath, t.Name, t.Type)
}

func (t TypeDecl) GetFieldsImports() ImportList {
	return lo.Uniq(t.Type.Imports())
}

type TypeDeclList []*TypeDecl

func NewTypeDeclList(pkg string, generalDecl *ast.GenDecl) TypeDeclList {
	imp := NewImportWithAlias(pkg)
	list := make(TypeDeclList, 0, len(generalDecl.Specs))
	for _, spec := range generalDecl.Specs {
		ts, ok := NewTypeDecl(imp, generalDecl, spec)
		if ok {
			list = append(list, ts)
		}
	}

	return list
}

func (l TypeDeclList) GetByName(name string) (*TypeDecl, bool) {
	return lo.Find(l, func(item *TypeDecl) bool {
		return item.Name == name
	})
}

func NewTypeDecl(imp Import, generalDecl *ast.GenDecl, spec ast.Spec) (*TypeDecl, bool) {
	castedSpec, isTypeSpec := spec.(*ast.TypeSpec)
	if !isTypeSpec {
		return nil, false
	}

	var specType Type
	switch castedType := castedSpec.Type.(type) {
	case *ast.StructType:
		specType = NewType(castedType)
	case *ast.InterfaceType:
		specType = NewType(castedType)
	case *ast.FuncType:
		specType = NewType(castedType)
	case *ast.ArrayType:
		specType = NewType(castedType)
	case *ast.MapType:
		specType = NewType(castedType)
	case *ast.Ident:
		specType = NewType(castedType)
	case *ast.ChanType:
		specType = NewType(castedType)
	case *ast.SelectorExpr:
		specType = NewType(castedType)
	default:
		panic(fmt.Sprintf("unknown type: %+[1]v(%[1]T)", castedType))
	}

	specName := castedSpec.Name.Name

	var specComment string
	if doc := castedSpec.Doc; doc != nil {
		specComment = doc.Text()
	} else if len(generalDecl.Specs) == 1 && generalDecl.Doc != nil {
		specComment = generalDecl.Doc.Text()
	}

	return &TypeDecl{
		Name:        specName,
		Comment:     strings.TrimSpace(specComment),
		Type:        specType,
		Package:     imp.Alias,
		PackagePath: imp.Path,
	}, true
}

func InspectTypeDeclTypes(typeDecl *TypeDecl, fn func(Type) error) error {
	return InspectType(
		typeDecl.Type,
		fn,
	)
}

func GetTypeDeclAllImportPath(typeDecl *TypeDecl) ([]string, error) {
	imports := make(map[string]struct{})

	err := InspectTypeDeclTypes(
		typeDecl,
		func(t Type) error {
			return InspectType(t, func(t Type) error {
				if casted, ok := t.(PackageGetterType); ok {
					if pkg := casted.GetPackagePath(); pkg != "" {
						imports[pkg] = struct{}{}
					}
				}

				return nil
			})
		},
	)
	if err != nil {
		return nil, fmt.Errorf("inspect type(name=%s): %w", typeDecl, err)
	}

	return lo.MapToSlice(
		imports,
		func(key string, _ struct{}) string { return key },
	), nil
}

package interface_creator

import (
	"errors"
	"fmt"

	"github.com/khevse/codegen/internal/pkg/astpkg"
	"github.com/samber/lo"
)

type field struct {
	Name     string
	TypeName string
}

type methodSpec struct {
	Name    string
	Comment string
	Params  []field
	Results []field
}

type objectSpec struct {
	Name    string
	Comment string
	Methods []methodSpec
}

func newObjectSpec(
	name string,
	typeDecl *astpkg.TypeDecl,
	methods []*astpkg.FuncDecl,
	imports astpkg.ImportList,
) (objectSpec, error) {
	if _, ok := typeDecl.Type.(*astpkg.StructType); !ok {
		return objectSpec{}, errors.New("type is not struct")
	}

	for _, decl := range methods {
		err := astpkg.InspectFuncDeclFields(decl, func(f *astpkg.Field) error {
			return astpkg.ReplaceImportAliasByImportPath(f.Type, imports)
		})
		if err != nil {
			return objectSpec{}, fmt.Errorf("replace imports(%s): %w", decl, err)
		}
	}

	methodList := make([]methodSpec, 0, len(methods))
	for _, item := range methods {
		if !astpkg.IsExported(item.Name) {
			continue
		}

		params := newFieldsList(item.Params)
		results := newFieldsList(item.Results)

		method := methodSpec{
			Name:    item.Name,
			Comment: item.Comment,
			Params:  params,
			Results: results,
		}

		methodList = append(methodList, method)
	}

	return objectSpec{
		Name:    name,
		Comment: fmt.Sprintf("%s interface for type %s: %s", name, typeDecl.Name, typeDecl.Comment),
		Methods: methodList,
	}, nil
}

func newFieldsList(src []*astpkg.Field) []field {
	fieldList := make([]field, 0, len(src))
	for _, item := range src {
		f := field{
			Name:     lo.Ternary(item.Name == "", "_", item.Name),
			TypeName: item.Type.ExprString(),
		}

		fieldList = append(fieldList, f)
	}

	return fieldList
}

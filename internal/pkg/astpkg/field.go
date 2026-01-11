package astpkg

import (
	"fmt"
	"go/ast"

	"github.com/samber/lo"
)

type Field struct {
	Name string
	Type Type
}

func NewFieldList(fieldList *ast.FieldList) []*Field {
	list := make([]*Field, 0, fieldList.NumFields())
	if fieldList != nil {
		for _, field := range fieldList.List {
			list = append(list, NewField(field)...)
		}
	}

	return list
}

func NewField(field *ast.Field) []*Field {
	if len(field.Names) == 0 {
		return []*Field{
			{
				Name: "",
				Type: NewType(field.Type),
			},
		}
	}

	list := make([]*Field, 0, len(field.Names))
	for _, nameIdent := range field.Names {
		list = append(
			list,
			&Field{
				Name: nameIdent.Name,
				Type: NewType(field.Type),
			},
		)
	}

	return list
}

func (f Field) String() string {
	name := lo.Ternary(f.Name == "", "_", f.Name)
	return fmt.Sprintf("%s %s", name, f.Type)
}

func InspectFields(fieldList []*Field, fn func(*Field) error) error {
	for _, field := range fieldList {
		if err := fn(field); err != nil {
			return fmt.Errorf("inspect field(name=%s): %w", field.Name, err)
		}
	}

	return nil
}

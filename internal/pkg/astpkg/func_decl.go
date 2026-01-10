package astpkg

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/samber/lo"
)

type FuncDecl struct {
	Receiver string
	Name     string
	Comment  string
	Params   []*Field
	Results  []*Field
}

func (t FuncDecl) String() string {
	if t.Receiver == "" {
		return fmt.Sprintf("func %s", t.Name)
	}

	return fmt.Sprintf("func (%s) %s", t.Receiver, t.Name)
}

type FuncDeclList []*FuncDecl

func (l FuncDeclList) GetByReceiverName(name string) FuncDeclList {
	return lo.Filter(l, func(item *FuncDecl, _ int) bool {
		return item.Receiver == name
	})
}

func NewFuncDecl(spec *ast.FuncDecl) *FuncDecl {
	var specComment string
	if doc := spec.Doc; doc != nil {
		specComment = strings.TrimSpace(doc.Text())
	}

	params := NewFieldList(spec.Type.Params)
	results := NewFieldList(spec.Type.Results)

	var recvName string
	if spec.Recv != nil && len(spec.Recv.List) == 1 {
		switch t := spec.Recv.List[0].Type.(type) {
		case *ast.Ident:
			recvName = t.Name
		case *ast.StarExpr:
			if id, ok := t.X.(*ast.Ident); ok {
				recvName = id.Name
			}
		}
	}

	return &FuncDecl{
		Receiver: recvName,
		Name:     spec.Name.Name,
		Comment:  specComment,
		Params:   params,
		Results:  results,
	}
}

func InspectFuncDeclFields(funcDecl *FuncDecl, fn func(*Field) error) error {
	if err := InspectFields(funcDecl.Params, fn); err != nil {
		return fmt.Errorf("inspect params: %w", err)
	}

	if err := InspectFields(funcDecl.Results, fn); err != nil {
		return fmt.Errorf("inspect results: %w", err)
	}

	return nil
}

func GetFuncDeclAllImportPath(funcDecl *FuncDecl) ([]string, error) {
	imports := make(map[string]struct{})

	err := InspectFuncDeclFields(
		funcDecl,
		func(f *Field) error {
			return InspectType(f.Type, func(t Type) error {
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
		return nil, fmt.Errorf("inspect func declaration(%s): %w", funcDecl, err)
	}

	return lo.MapToSlice(
		imports,
		func(key string, _ struct{}) string { return key },
	), nil
}

package astpkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestNewFuncDecl(t *testing.T) {
	t.Parallel()

	t.Run("function with params", func(t *testing.T) {
		funcDecl := newFuncDeclForTest(
			t,
			`package p;

			// test comment
			func test(_ string) {};`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "test comment",
				Params:   []*Field{{Name: "_", Type: &Ident{Name: "string"}}},
				Results:  []*Field{},
			},
			funcDecl,
		)
	})

	t.Run("function without params", func(t *testing.T) {
		funcDecl := newFuncDeclForTest(
			t,
			`package p;

			// test comment
			func test() {};`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "test comment",
				Params:   []*Field{},
				Results:  []*Field{},
			},
			funcDecl,
		)
	})

	t.Run("method of receiver", func(t *testing.T) {
		funcDecl := newFuncDeclForTest(
			t,
			`package p;

			// test comment
			func (s *Struct) test() {};`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "Struct",
				Name:     "test",
				Comment:  "test comment",
				Params:   []*Field{},
				Results:  []*Field{},
			},
			funcDecl,
		)
	})

	t.Run("method without receiver", func(t *testing.T) {
		funcDecl := newFuncDeclForTest(
			t,
			`package p;

			// test comment
			func (Struct) test() {};`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "Struct",
				Name:     "test",
				Comment:  "test comment",
				Params:   []*Field{},
				Results:  []*Field{},
			},
			funcDecl,
		)
	})
}

func TestGetFuncDeclAllImportPath(t *testing.T) {
	t.Parallel()

	pkg, err := ParsePackage("github.com/khevse/codegen/tests/mainpkg")
	require.NoError(t, err)

	t.Run("external imports", func(t *testing.T) {
		decl, exists := lo.Find(pkg.FuncDeclList, func(item *FuncDecl) bool {
			return item.Name == "SetAllFields"
		})
		require.True(t, exists)

		imports, err := GetFuncDeclAllImportPath(decl)
		require.NoError(t, err)
		require.Equal(
			t,
			[]string{"github.com/khevse/codegen/tests/mainpkg/childpkg"},
			imports,
		)
	})

	require.NoError(t, SetPackagePathForAllDecl(pkg))

	t.Run("SetAllFields", func(t *testing.T) {
		decl, exists := lo.Find(pkg.FuncDeclList, func(item *FuncDecl) bool {
			return item.Name == "SetAllFields"
		})
		require.True(t, exists)

		imports, err := GetFuncDeclAllImportPath(decl)
		slices.Sort(imports)

		require.NoError(t, err)
		require.Equal(
			t,
			[]string{
				"github.com/khevse/codegen/tests/mainpkg",
				"github.com/khevse/codegen/tests/mainpkg/childpkg",
			},
			imports,
		)
	})

	t.Run("SetFieldStruct", func(t *testing.T) {
		decl, exists := lo.Find(pkg.FuncDeclList, func(item *FuncDecl) bool {
			return item.Name == "SetFieldStruct"
		})
		require.True(t, exists)

		imports, err := GetFuncDeclAllImportPath(decl)
		require.NoError(t, err)
		require.Equal(
			t,
			[]string{"github.com/khevse/codegen/tests/mainpkg/childpkg"},
			imports,
		)
	})
}

func newFuncDeclForTest(t *testing.T, code string) *FuncDecl {
	f, err := parser.ParseFile(
		token.NewFileSet(),
		"",
		code,
		parser.DeclarationErrors|parser.ParseComments,
	)
	require.NoError(t, err)

	declList := lo.FilterMap(f.Decls, func(item ast.Decl, _ int) (*ast.FuncDecl, bool) {
		casted, ok := item.(*ast.FuncDecl)
		return casted, ok
	})

	list := make([]*FuncDecl, 0, len(declList))
	for _, decl := range declList {
		list = append(list, NewFuncDecl(decl))
	}
	require.Len(t, list, 1)
	return list[0]
}

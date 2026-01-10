package astpkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTypeDeclList(t *testing.T) {
	t.Parallel()

	getTypeDeclList := func(t *testing.T, code string) []*TypeDecl {
		f, err := parser.ParseFile(
			token.NewFileSet(),
			"",
			code,
			parser.DeclarationErrors|parser.ParseComments,
		)
		require.NoError(t, err)
		if len(f.Decls) > 0 {
			list := make([]*TypeDecl, 0, len(f.Decls))
			for _, decl := range f.Decls {
				list = append(list, NewTypeDeclList("./test", decl.(*ast.GenDecl))...)
			}
			return list
		}

		return nil
	}

	t.Run("without types", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;`,
		)
		require.Empty(t, typeDeclList)
	})

	t.Run("one struct with comment", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			// Struct comment
			type Struct struct{}`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "Struct",
					Comment:     "Struct comment",
					Type:        &StructType{Fields: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)

	})

	t.Run("one struct without comment", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;
			type Struct struct{}`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "Struct",
					Comment:     "",
					Type:        &StructType{Fields: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)

	})

	t.Run("many structs", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			// Struct1 comment
			type Struct1 struct{}

			// Struct2 comment
			type Struct2 struct{}`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "Struct1",
					Comment:     "Struct1 comment",
					Type:        &StructType{Fields: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
				{
					Name:        "Struct2",
					Comment:     "Struct2 comment",
					Type:        &StructType{Fields: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})

	t.Run("many types", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			// Struct comment
			type Struct struct{}

			// Interface comment
			type Interface interface{ }`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "Struct",
					Comment:     "Struct comment",
					Type:        &StructType{Fields: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
				{
					Name:        "Interface",
					Comment:     "Interface comment",
					Type:        &InterfaceType{Methods: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})

	t.Run("func types", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			// Func comment
			type Func func()`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "Func",
					Comment:     "Func comment",
					Type:        &FuncType{Params: []*Field{}, Results: []*Field{}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})

	t.Run("array type", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			type List []string`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "List",
					Comment:     "",
					Type:        &ArrayType{Type: &Ident{Name: "string"}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})

	t.Run("map type", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			type Map map[string]string`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "Map",
					Comment:     "",
					Type:        &MapType{Key: &Ident{Name: "string"}, Value: &Ident{Name: "string"}},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})

	t.Run("ident type", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			type T string`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "T",
					Comment:     "",
					Type:        &Ident{Name: "string"},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})

	t.Run("chan type", func(t *testing.T) {
		typeDeclList := getTypeDeclList(
			t,
			`package p;

			type C chan string`,
		)
		require.Equal(
			t,
			[]*TypeDecl{
				{
					Name:        "C",
					Comment:     "",
					Type:        &ChanType{Type: &Ident{Name: "string"}, Direction: ast.SEND | ast.RECV},
					Package:     "test",
					PackagePath: "./test",
				},
			},
			typeDeclList,
		)
	})
}

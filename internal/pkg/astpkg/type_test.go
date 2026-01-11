package astpkg

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestType(t *testing.T) {
	t.Parallel()

	t.Run("Ident", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val string) string { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params:   []*Field{{Name: "val", Type: &Ident{Name: "string"}}},
				Results:  []*Field{{Name: "", Type: &Ident{Name: "string"}}},
			},
			decl,
		)
		require.Equal(t, "val Ident(string)", decl.Params[0].String())
		require.Equal(t, "_ Ident(string)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("Ident with type", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; type Struct struct{}; func test(val Struct) Struct { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &Ident{
							Name: "Struct",
							Type: &StructType{Fields: []*Field{}},
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &Ident{
							Name: "Struct",
							Type: &StructType{Fields: []*Field{}},
						},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val Ident(Struct)", decl.Params[0].String())
		require.Equal(t, "_ Ident(Struct)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("StarExpr", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val *string) *string { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &StarExpr{Type: &Ident{Name: "string"}},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &StarExpr{Type: &Ident{Name: "string"}},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val StarExpr(*string)", decl.Params[0].String())
		require.Equal(t, "_ StarExpr(*string)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("ArrayType", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val []string) []string { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &ArrayType{Type: &Ident{Name: "string"}},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &ArrayType{Type: &Ident{Name: "string"}},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val ArrayType([]string)", decl.Params[0].String())
		require.Equal(t, "_ ArrayType([]string)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("MapType", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val map[string]string) map[string]string { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &MapType{
							Key:   &Ident{Name: "string"},
							Value: &Ident{Name: "string"},
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &MapType{
							Key:   &Ident{Name: "string"},
							Value: &Ident{Name: "string"},
						},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val MapType(map[string]string)", decl.Params[0].String())
		require.Equal(t, "_ MapType(map[string]string)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("SelectorExpr", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p;
			import p2 "./example"
			func test(val p2.Struct) p2.Struct { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &SelectorExpr{Package: "p2", Name: "Struct", Type: nil},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &SelectorExpr{Package: "p2", Name: "Struct", Type: nil},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val SelectorExpr(p2.Struct(<nil>))", decl.Params[0].String())
		require.Equal(t, "_ SelectorExpr(p2.Struct(<nil>))", decl.Results[0].String())
		require.Equal(t, ImportList{{Alias: "p2", Path: ""}}, decl.GetSignatureImports())
	})

	t.Run("Ellipsis", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val ...any) []any { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &EllipsisType{Type: &Ident{Name: "any"}},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &ArrayType{Type: &Ident{Name: "any"}},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val EllipsisType(...any)", decl.Params[0].String())
		require.Equal(t, "_ ArrayType([]any)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("FuncType", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val func() string) string { return val() };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &FuncType{
							Params:  []*Field{},
							Results: []*Field{{Name: "", Type: &Ident{Name: "string"}}},
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &Ident{Name: "string"},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val FuncType(func() (_ string))", decl.Params[0].String())
		require.Equal(t, "_ Ident(string)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("StructType", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val struct{}) struct{} { return val() };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &StructType{Fields: []*Field{}},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &StructType{Fields: []*Field{}},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val StructType(struct{})", decl.Params[0].String())
		require.Equal(t, "_ StructType(struct{})", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("InterfaceType", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val interface{}) interface{} { return val() };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &InterfaceType{Methods: []*Field{}},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &InterfaceType{Methods: []*Field{}},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val InterfaceType(interface{})", decl.Params[0].String())
		require.Equal(t, "_ InterfaceType(interface{})", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("ChanType", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p; func test(val chan string) chan string { return val() };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{

					{
						Name: "val",
						Type: &ChanType{
							Type:      &Ident{Name: "string"},
							Direction: ast.RECV | ast.SEND,
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &ChanType{
							Type:      &Ident{Name: "string"},
							Direction: ast.RECV | ast.SEND,
						},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val ChanType(chan string)", decl.Params[0].String())
		require.Equal(t, "_ ChanType(chan string)", decl.Results[0].String())
		require.Empty(t, decl.GetSignatureImports())
	})

	t.Run("IndexExpr", func(t *testing.T) {
		decl := newFuncDeclForTest(
			t,
			`package p;
			import p2 "./example"
			func test(val *p2.Item[p2.ID]) *p2.Item[p2.ID] { return val };`,
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &StarExpr{
							Type: &IndexExpr{
								Index: &SelectorExpr{
									Package:     "p2",
									PackagePath: "",
									Name:        "ID",
									Type:        nil,
								},
								X: &SelectorExpr{
									Package:     "p2",
									PackagePath: "",
									Name:        "Item",
									Type:        nil,
								},
							},
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &StarExpr{
							Type: &IndexExpr{
								Index: &SelectorExpr{
									Package:     "p2",
									PackagePath: "",
									Name:        "ID",
									Type:        nil,
								},
								X: &SelectorExpr{
									Package:     "p2",
									PackagePath: "",
									Name:        "Item",
									Type:        nil,
								},
							},
						},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val StarExpr(*p2.Item[p2.ID])", decl.Params[0].String())
		require.Equal(t, "_ StarExpr(*p2.Item[p2.ID])", decl.Results[0].String())
		require.Equal(t, ImportList{{Alias: "p2", Path: ""}}, decl.GetSignatureImports())
	})
}

func TestSetPackageInformation(t *testing.T) {
	newFuncDecl := func(t *testing.T, code string, imports ImportList) *FuncDecl {
		decl := newFuncDeclForTest(t, code)
		err := InspectFuncDeclFields(
			decl,
			func(f *Field) error {
				return InspectType(f.Type, func(t Type) error {
					return SetPackageInformation(t, imports)
				})
			},
		)
		require.NoError(t, err)
		return decl
	}

	t.Run("nil", func(t *testing.T) {
		err := SetPackageInformation(nil, ImportList{})
		require.NoError(t, err)
	})

	t.Run("Ident", func(t *testing.T) {
		decl := newFuncDecl(
			t,
			`package p; func test(val string) string { return val };`,
			ImportList{},
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &Ident{Name: "string"},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &Ident{Name: "string"},
					},
				},
			},
			decl,
		)
		require.Equal(t, "val Ident(string)", decl.Params[0].String())
		require.Equal(t, "_ Ident(string)", decl.Results[0].String())
	})

	t.Run("Ident with type", func(t *testing.T) {
		decl := newFuncDecl(
			t,
			`package p; type Struct struct{}; func test(val Struct) Struct { return val };`,
			ImportList{},
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &Ident{
							Package: "",
							Name:    "Struct",
							Type:    &StructType{Fields: []*Field{}},
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &Ident{
							Name: "Struct",
							Type: &StructType{Fields: []*Field{}},
						},
					},
				},
			},
			decl,
		)
	})

	t.Run("SelectorExpr with import", func(t *testing.T) {
		decl := newFuncDecl(
			t,
			`package p;
			import p2 "./example"
			func test(val p2.Struct) p2.Struct { return val };`,
			ImportList{{Alias: "p2", Path: "./example"}},
		)
		require.Equal(
			t,
			&FuncDecl{
				Receiver: "",
				Name:     "test",
				Comment:  "",
				Params: []*Field{
					{
						Name: "val",
						Type: &SelectorExpr{
							PackagePath: "./example",
							Package:     "p2",
							Name:        "Struct",
							Type:        nil,
						},
					},
				},
				Results: []*Field{
					{
						Name: "",
						Type: &SelectorExpr{
							PackagePath: "./example",
							Package:     "p2",
							Name:        "Struct",
							Type:        nil,
						},
					},
				},
			},
			decl,
		)
	})
}

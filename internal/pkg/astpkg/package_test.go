package astpkg

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestParsePackage(t *testing.T) {
	t.Parallel()

	t.Run("failed", func(t *testing.T) {
		pkg, err := ParsePackage("-")
		require.EqualError(
			t,
			err,
			`package error(name=-): [malformed import path "-": leading dash(position:)]`,
		)
		require.Nil(t, pkg)
	})

	t.Run("success", func(t *testing.T) {
		_, file, _, _ := runtime.Caller(0)
		wantDir, err := filepath.Abs(filepath.Join(filepath.Dir(file), "../../..", "tests/mainpkg"))
		require.NoError(t, err)

		want := &Package{
			Path: "github.com/khevse/codegen/tests/mainpkg",
			Dir:  wantDir,
			TypeDeclList: []*TypeDecl{
				{
					Name:    "StructWithMethods",
					Comment: "StructWithMethods comment",
					Type: &StructType{
						Fields: []*Field{
							{
								Name: "FieldStruct",
								Type: &SelectorExpr{
									PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
									Package:     "",
									Name:        "Struct",
								},
							},
							{
								Name: "FieldString",
								Type: &Ident{
									Name: "string",
								},
							},
						},
					},
					Package:     "mainpkg",
					PackagePath: "github.com/khevse/codegen/tests/mainpkg",
				},
			},
			FuncDeclList: []*FuncDecl{
				{
					Receiver: "StructWithMethods",
					Name:     "GetFieldString",
					Params:   []*Field{},
					Results:  []*Field{{Name: "", Type: &Ident{Name: "string"}}},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "GetFieldStruct",
					Comment:  "GetFieldStruct comment",
					Params:   []*Field{},
					Results: []*Field{
						{
							Name: "",
							Type: &SelectorExpr{
								PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
								Package:     "",
								Name:        "Struct",
							},
						},
					},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "SetFieldStruct",
					Params: []*Field{
						{
							Name: "val",
							Type: &SelectorExpr{
								PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
								Package:     "childpkgalias",
								Name:        "Struct",
							},
						},
					},
					Results: []*Field{},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "SetAllFields",
					Params: []*Field{
						{
							Name: "val",
							Type: &Ident{
								Name: "StructWithMethods",
								Type: &StructType{
									Fields: []*Field{
										{
											Name: "FieldStruct",
											Type: &SelectorExpr{
												PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
												Package:     "",
												Name:        "Struct",
											},
										},
										{
											Name: "FieldString",
											Type: &Ident{Name: "string"},
										},
									},
								},
							},
						},
					},
					Results: []*Field{},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "SetFieldStringFromInterface",
					Params: []*Field{
						{
							Name: "val",
							Type: &SelectorExpr{
								PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
								Package:     "",
								Name:        "Interface",
							},
						},
					},
					Results: []*Field{},
				},
			},
		}

		pkg, err := ParsePackage("github.com/khevse/codegen/tests/mainpkg")
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(
			want,
			func() *Package {
				pkg.TypeDeclList = lo.Filter(pkg.TypeDeclList, func(item *TypeDecl, _ int) bool {
					return item.Name == "StructWithMethods"
				})
				pkg.FuncDeclList = lo.Filter(pkg.FuncDeclList, func(item *FuncDecl, _ int) bool {
					return item.Receiver == "StructWithMethods"
				})
				return pkg
			}(),
			cmpopts.SortSlices(func(i, j *TypeDecl) bool {
				return strings.Compare(i.Name, j.Name) < 0
			}),
			cmpopts.SortSlices(func(i, j *FuncDecl) bool {
				return strings.Compare(i.Name, j.Name) < 0
			}),
			cmpopts.SortSlices(func(i, j *Field) bool {
				return strings.Compare(i.Type.String(), j.Type.String()) < 0
			}),
		))
	})
}

func TestGetPackageDir(t *testing.T) {
	t.Parallel()

	t.Run("invalid path", func(t *testing.T) {
		dir, err := GetPackagePath("-")
		require.EqualError(t, err, `load package:err: chdir -: no such file or directory: stderr: `)
		require.Empty(t, dir)
	})

	t.Run("invalid package", func(t *testing.T) {
		dir, err := GetPackagePath("./..")
		require.EqualError(t, err, `package error(name=.): [no Go files in /home/jack/work/go/src/codegen/internal/pkg(position:)]`)
		require.Empty(t, dir)
	})

	t.Run("success", func(t *testing.T) {
		_, file, _, _ := runtime.Caller(0)
		pkgDir, err := filepath.Abs(filepath.Join(filepath.Dir(file), "../../..", "tests/mainpkg"))
		require.NoError(t, err)

		dir, err := GetPackagePath(pkgDir)
		require.NoError(t, err)
		require.Equal(t, "github.com/khevse/codegen/tests/mainpkg", dir)
	})
}

func TestSetPackagePathForAllDecl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, file, _, _ := runtime.Caller(0)
		wantDir, err := filepath.Abs(filepath.Join(filepath.Dir(file), "../../..", "tests/mainpkg"))
		require.NoError(t, err)

		want := &Package{
			Path: "github.com/khevse/codegen/tests/mainpkg",
			Dir:  wantDir,
			TypeDeclList: []*TypeDecl{
				{
					Name:    "StructWithMethods",
					Comment: "StructWithMethods comment",
					Type: &StructType{
						Fields: []*Field{
							{
								Name: "FieldStruct",
								Type: &SelectorExpr{
									PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
									Package:     "",
									Name:        "Struct",
								},
							},
							{
								Name: "FieldString",
								Type: &Ident{
									Name: "string",
								},
							},
						},
					},
					Package:     "mainpkg",
					PackagePath: "github.com/khevse/codegen/tests/mainpkg",
				},
				{
					Name:    "Factory",
					Comment: "Factory comment",
					Type: &StructType{
						Fields: []*Field{},
					},
					Package:     "mainpkg",
					PackagePath: "github.com/khevse/codegen/tests/mainpkg",
				},
				{
					Name:    "IFactory",
					Comment: "IFactory .",
					Type: &InterfaceType{
						Methods: []*Field{
							{
								Name: "NewObject1",
								Type: &FuncType{
									Params: []*Field{},
									Results: []*Field{
										{
											Name: "",
											Type: &Ident{
												Package:     "mainpkg",
												PackagePath: "github.com/khevse/codegen/tests/mainpkg",
												Name:        "IObject1",
												Type: &InterfaceType{
													Methods: []*Field{
														{
															Name: "String",
															Type: &FuncType{
																Params: []*Field{},
																Results: []*Field{
																	{
																		Name: "",
																		Type: &Ident{
																			Package:     "",
																			PackagePath: "",
																			Name:        "string",
																			Type:        nil,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							{
								Name: "NewObject2",
								Type: &FuncType{
									Params: []*Field{},
									Results: []*Field{
										{
											Name: "",
											Type: &Ident{
												Package:     "mainpkg",
												PackagePath: "github.com/khevse/codegen/tests/mainpkg",
												Name:        "IObject2",
												Type: &InterfaceType{
													Methods: []*Field{
														{
															Name: "String",
															Type: &FuncType{
																Params: []*Field{},
																Results: []*Field{
																	{
																		Name: "",
																		Type: &Ident{
																			Package:     "",
																			PackagePath: "",
																			Name:        "string",
																			Type:        nil,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Package:     "mainpkg",
					PackagePath: "github.com/khevse/codegen/tests/mainpkg",
				},
			},
			FuncDeclList: []*FuncDecl{
				{
					Receiver: "StructWithMethods",
					Name:     "GetFieldString",
					Params:   []*Field{},
					Results:  []*Field{{Name: "", Type: &Ident{Name: "string"}}},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "GetFieldStruct",
					Comment:  "GetFieldStruct comment",
					Params:   []*Field{},
					Results: []*Field{
						{
							Name: "",
							Type: &SelectorExpr{
								PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
								Package:     "",
								Name:        "Struct",
							},
						},
					},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "SetFieldStruct",
					Params: []*Field{
						{
							Name: "val",
							Type: &SelectorExpr{
								PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
								Package:     "childpkgalias",
								Name:        "Struct",
							},
						},
					},
					Results: []*Field{},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "SetAllFields",
					Params: []*Field{
						{
							Name: "val",
							Type: &Ident{
								Name: "StructWithMethods",
								Type: &StructType{
									Fields: []*Field{
										{
											Name: "FieldStruct",
											Type: &SelectorExpr{
												PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
												Package:     "",
												Name:        "Struct",
											},
										},
										{
											Name: "FieldString",
											Type: &Ident{Name: "string"},
										},
									},
								},
								Package:     "mainpkg",
								PackagePath: "github.com/khevse/codegen/tests/mainpkg",
							},
						},
					},
					Results: []*Field{},
				},
				{
					Receiver: "StructWithMethods",
					Name:     "SetFieldStringFromInterface",
					Params: []*Field{
						{
							Name: "val",
							Type: &SelectorExpr{
								PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
								Package:     "",
								Name:        "Interface",
							},
						},
					},
					Results: []*Field{},
				},
			},
		}

		pkg, err := ParsePackage("github.com/khevse/codegen/tests/mainpkg")
		require.NoError(t, err)

		require.NoError(t, SetPackagePathForAllDecl(pkg))
		require.Empty(t, cmp.Diff(
			want,
			func() *Package {
				pkg.TypeDeclList = lo.Filter(pkg.TypeDeclList, func(item *TypeDecl, _ int) bool {
					return item.Name == "StructWithMethods" ||
						item.Name == "Factory" ||
						item.Name == "IFactory"
				})
				pkg.FuncDeclList = lo.Filter(pkg.FuncDeclList, func(item *FuncDecl, _ int) bool {
					return item.Receiver == "StructWithMethods"
				})
				return pkg
			}(),
			cmpopts.SortSlices(func(i, j *TypeDecl) bool {
				return strings.Compare(i.Name, j.Name) < 0
			}),
			cmpopts.SortSlices(func(i, j *FuncDecl) bool {
				return strings.Compare(i.Name, j.Name) < 0
			}),
			cmpopts.SortSlices(func(i, j *Field) bool {
				return strings.Compare(i.Type.String(), j.Type.String()) < 0
			}),
		))
	})
}

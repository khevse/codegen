package interface_creator

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/khevse/codegen/internal/pkg/astpkg"
	"github.com/stretchr/testify/require"
)

func TestSplitFromType(t *testing.T) {
	t.Parallel()

	t.Run("success simple form", func(t *testing.T) {
		res, err := parseFromType("github.com/package.TypeName")
		require.NoError(t, err)
		require.Equal(t,
			[]argFromType{{
				Package:    "github.com/package",
				SourceName: "TypeName",
				TargetName: "TypeName",
			}},
			res,
		)
	})

	t.Run("success with interface name", func(t *testing.T) {
		res, err := parseFromType("github.com/package.TypeName=ITypeName")
		require.NoError(t, err)
		require.Equal(t,
			[]argFromType{{
				Package:    "github.com/package",
				SourceName: "TypeName",
				TargetName: "ITypeName",
			}},
			res,
		)
	})

	t.Run("success with multiple types", func(t *testing.T) {
		res, err := parseFromType("github.com/package.TypeName1=ITypeName1, github.com/package.TypeName2=ITypeName2")
		require.NoError(t, err)
		require.Equal(t,
			[]argFromType{
				{
					Package:    "github.com/package",
					SourceName: "TypeName1",
					TargetName: "ITypeName1",
				},
				{
					Package:    "github.com/package",
					SourceName: "TypeName2",
					TargetName: "ITypeName2",
				},
			},
			res,
		)
	})

	t.Run("success with duplicates", func(t *testing.T) {
		res, err := parseFromType("github.com/package.TypeName1=ITypeName1,github.com/package.TypeName1=ITypeName1")
		require.NoError(t, err)
		require.Equal(t,
			[]argFromType{
				{
					Package:    "github.com/package",
					SourceName: "TypeName1",
					TargetName: "ITypeName1",
				},
			},
			res,
		)
	})

	t.Run("failed", func(t *testing.T) {
		res, err := parseFromType("invalid name")
		require.EqualError(t, err, "invalid type: invalid name")
		require.Empty(t, res)
	})
}

func TestPrepareObjectSpecList(t *testing.T) {
	t.Parallel()

	t.Run("success to self package", func(t *testing.T) {
		args := commandArgs{
			fromType:   "github.com/khevse/codegen/tests/mainpkg.StructWithMethods",
			targetDir:  "./../../../tests/mainpkg",
			fileSuffix: "",
		}
		importList, list, err := prepareObjectSpecList(args)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(
			astpkg.ImportList{
				{Alias: "childpkg", Path: "github.com/khevse/codegen/tests/mainpkg/childpkg"},
			},
			importList,
			cmpopts.SortSlices(func(i, j astpkg.Import) bool {
				return i.Path < j.Path
			}),
		))
		require.Empty(t, cmp.Diff(
			[]objectSpec{
				{
					Name:    "StructWithMethods",
					Comment: "StructWithMethods interface for type StructWithMethods: StructWithMethods comment",
					Methods: []methodSpec{
						{
							Name:    "GetFieldStruct",
							Comment: "GetFieldStruct comment",
							Params:  []field{},
							Results: []field{
								{
									Name:     "_",
									TypeName: "childpkg.Struct",
									Type: &astpkg.SelectorExpr{
										Package:     "childpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
										Name:        "Struct",
										Type:        nil,
									},
								},
							},
						},
						{
							Name:    "GetFieldString",
							Comment: "",
							Params:  []field{},
							Results: []field{
								{
									Name:     "_",
									TypeName: "string",
									Type: &astpkg.Ident{
										Package:     "",
										PackagePath: "",
										Name:        "string",
										Type:        nil,
									},
								},
							},
						},
						{
							Name:    "SetFieldStruct",
							Comment: "",
							Params: []field{
								{
									Name:     "val",
									TypeName: "childpkg.Struct",
									Type: &astpkg.SelectorExpr{
										Package:     "childpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
										Name:        "Struct",
										Type:        nil,
									},
								},
							},
							Results: []field{},
						},
						{
							Name:    "SetAllFields",
							Comment: "",
							Params: []field{
								{
									Name:     "val",
									TypeName: "StructWithMethods",
									Type: &astpkg.Ident{
										Package:     "",
										PackagePath: "",
										Name:        "StructWithMethods",
										Type: &astpkg.StructType{
											Fields: []*astpkg.Field{
												{
													Name: "FieldStruct",
													Type: &astpkg.SelectorExpr{
														Package:     "childpkg",
														PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
														Name:        "Struct",
														Type:        nil,
													},
												},
												{
													Name: "FieldString",
													Type: &astpkg.Ident{
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
							Results: []field{},
						},
						{
							Name:    "SetFieldStringFromInterface",
							Comment: "",
							Params: []field{
								{
									Name:     "val",
									TypeName: "childpkg.Interface",
									Type: &astpkg.SelectorExpr{
										Package:     "childpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
										Name:        "Interface",
										Type:        nil,
									},
								},
							},
							Results: []field{},
						},
					},
				},
			},
			list,
			cmpopts.SortSlices(func(i, j methodSpec) bool {
				return i.Name < j.Name
			}),
		))
	})

	t.Run("success with another target package and another name", func(t *testing.T) {
		args := commandArgs{
			fromType:   "github.com/khevse/codegen/tests/mainpkg.StructWithMethods=StructWithMethodsOther",
			targetDir:  "./",
			fileSuffix: "",
		}
		importList, list, err := prepareObjectSpecList(args)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(
			astpkg.ImportList{
				{Alias: "mainpkg", Path: "github.com/khevse/codegen/tests/mainpkg"},
				{Alias: "childpkg", Path: "github.com/khevse/codegen/tests/mainpkg/childpkg"},
			},
			importList,
			cmpopts.SortSlices(func(i, j astpkg.Import) bool {
				return i.Path < j.Path
			}),
		))
		require.Empty(t, cmp.Diff(
			[]objectSpec{
				{
					Name:    "StructWithMethodsOther",
					Comment: "StructWithMethodsOther interface for type StructWithMethods: StructWithMethods comment",
					Methods: []methodSpec{
						{
							Name:    "GetFieldStruct",
							Comment: "GetFieldStruct comment",
							Params:  []field{},
							Results: []field{
								{
									Name:     "_",
									TypeName: "childpkg.Struct",
									Type: &astpkg.SelectorExpr{
										Package:     "childpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
										Name:        "Struct",
										Type:        nil,
									},
								},
							},
						},
						{
							Name:    "GetFieldString",
							Comment: "",
							Params:  []field{},
							Results: []field{
								{
									Name:     "_",
									TypeName: "string",
									Type: &astpkg.Ident{
										Package:     "",
										PackagePath: "",
										Name:        "string",
										Type:        nil,
									},
								},
							},
						},
						{
							Name:    "SetFieldStruct",
							Comment: "",
							Params: []field{
								{
									Name:     "val",
									TypeName: "childpkg.Struct",
									Type: &astpkg.SelectorExpr{
										Package:     "childpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
										Name:        "Struct",
										Type:        nil,
									},
								},
							},
							Results: []field{},
						},
						{
							Name:    "SetAllFields",
							Comment: "",
							Params: []field{
								{
									Name:     "val",
									TypeName: "mainpkg.StructWithMethods",
									Type: &astpkg.Ident{
										Package:     "mainpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg",
										Name:        "StructWithMethods",
										Type: &astpkg.StructType{
											Fields: []*astpkg.Field{
												{
													Name: "FieldStruct",
													Type: &astpkg.SelectorExpr{
														Package:     "childpkg",
														PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
														Name:        "Struct",
														Type:        nil,
													},
												},
												{
													Name: "FieldString",
													Type: &astpkg.Ident{
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
							Results: []field{},
						},
						{
							Name:    "SetFieldStringFromInterface",
							Comment: "",
							Params: []field{
								{
									Name:     "val",
									TypeName: "childpkg.Interface",
									Type: &astpkg.SelectorExpr{
										Package:     "childpkg",
										PackagePath: "github.com/khevse/codegen/tests/mainpkg/childpkg",
										Name:        "Interface",
										Type:        nil,
									},
								},
							},
							Results: []field{},
						},
					},
				},
			},
			list,
			cmpopts.SortSlices(func(i, j methodSpec) bool {
				return i.Name < j.Name
			}),
		))
	})
}

func TestExecute(t *testing.T) {
	args := commandArgs{
		fromType:   "github.com/khevse/codegen/tests/mainpkg.StructWithMethods=IStructWithMethods",
		targetDir:  "./",
		fileSuffix: "_generated",
	}
	require.NoError(t, (&Command{args: args}).Execute())

	const wantFile = "interfaces_generated.go"
	defer func() {
		require.NoError(t, os.Remove(wantFile))
	}()

	f, err := os.Open(wantFile)
	require.NoError(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)
	require.Equal(
		t,
		`// Code generated by http://github.com/khevse/codegen(version:; commit:; build:). DO NOT EDIT.

package interface_creator

import (
	mainpkg "github.com/khevse/codegen/tests/mainpkg"
	childpkg "github.com/khevse/codegen/tests/mainpkg/childpkg"
)

/* IStructWithMethods interface for type StructWithMethods: StructWithMethods comment */
type IStructWithMethods interface {
	GetFieldString() (_ string)
	/* GetFieldStruct comment */
	GetFieldStruct() (_ childpkg.Struct)
	SetAllFields(val mainpkg.StructWithMethods)
	SetFieldStringFromInterface(val childpkg.Interface)
	SetFieldStruct(val childpkg.Struct)
}
`,
		string(data),
	)
}

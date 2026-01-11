package object_test_wrapper

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/khevse/codegen/internal/pkg/astpkg"
	"github.com/stretchr/testify/require"
)

func TestPrepareObjectSpec(t *testing.T) {
	t.Parallel()

	t.Run("success to self package", func(t *testing.T) {
		args := commandArgs{
			interfaceType: "github.com/khevse/codegen/tests/mainpkg.IFactory=FactoryWrapper",
			targetDir:     "./../../../tests/mainpkg",
			fileSuffix:    "",
			objectType:    "github.com/khevse/codegen/tests/mainpkg.Factory",
			mockPackage:   "github.com/khevse/codegen/tests/mainpkg/mocks",
		}
		importList, spec, err := prepareObjectSpec(args)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(
			astpkg.ImportList{
				{Alias: "mocks", Path: "github.com/khevse/codegen/tests/mainpkg/mocks"},
			},
			importList,
			cmpopts.SortSlices(func(i, j astpkg.Import) bool {
				return i.Path < j.Path
			}),
		))
		require.Empty(t, cmp.Diff(
			&objectSpec{
				Name:    "FactoryWrapper",
				Comment: "FactoryWrapper wrapper for type IFactory: IFactory .",
				Methods: []methodSpec{
					{
						Name:    "NewObject1",
						Comment: "NewObject1 .",
						Params:  []field{},
						Results: []field{
							{
								FuncSpecName:   "_",
								TypeName:       "IObject1",
								ObjectSpecName: "IObject1",
								MockTypeName:   "IObject1Mock",
								MockPackage:    "mocks",
								Type: &astpkg.Ident{
									Package:     "",
									PackagePath: "",
									Name:        "IObject1",
									Type: &astpkg.InterfaceType{
										Methods: []*astpkg.Field{
											{
												Name: "String",
												Type: &astpkg.FuncType{
													Params: []*astpkg.Field{},
													Results: []*astpkg.Field{
														{
															Name: "",
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
								},
							},
						},
					},
					{
						Name:    "NewObject2",
						Comment: "NewObject2 .",
						Params:  []field{},
						Results: []field{
							{
								FuncSpecName:   "_",
								TypeName:       "IObject2",
								ObjectSpecName: "IObject2",
								MockTypeName:   "IObject2Mock",
								MockPackage:    "mocks",
								Type: &astpkg.Ident{
									Package:     "",
									PackagePath: "",
									Name:        "IObject2",
									Type: &astpkg.InterfaceType{
										Methods: []*astpkg.Field{
											{
												Name: "String",
												Type: &astpkg.FuncType{
													Params: []*astpkg.Field{},
													Results: []*astpkg.Field{
														{
															Name: "",
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
								},
							},
						},
					},
				},
				Fields: []objectSpecField{
					{
						Name:         "IObject1",
						TypeName:     "IObject1",
						MockPackage:  "mocks",
						MockTypeName: "IObject1Mock",
						Type: &astpkg.Ident{
							Package:     "",
							PackagePath: "",
							Name:        "IObject1",
							Type: &astpkg.InterfaceType{
								Methods: []*astpkg.Field{
									{
										Name: "String",
										Type: &astpkg.FuncType{
											Params: []*astpkg.Field{},
											Results: []*astpkg.Field{
												{
													Name: "",
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
						},
					},
					{
						Name:         "IObject2",
						TypeName:     "IObject2",
						MockPackage:  "mocks",
						MockTypeName: "IObject2Mock",
						Type: &astpkg.Ident{
							Package:     "",
							PackagePath: "",
							Name:        "IObject2",
							Type: &astpkg.InterfaceType{
								Methods: []*astpkg.Field{
									{
										Name: "String",
										Type: &astpkg.FuncType{
											Params: []*astpkg.Field{},
											Results: []*astpkg.Field{
												{
													Name: "",
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
						},
					},
				},
				BaseObjectTypeName: "*Factory",
			},
			spec,
			cmpopts.SortSlices(func(i, j methodSpec) bool {
				return i.Name < j.Name
			}),
		))
	})

	t.Run("success with another target package and another name", func(t *testing.T) {
		args := commandArgs{
			interfaceType: "github.com/khevse/codegen/tests/mainpkg.IFactory=FactoryWrapper",
			targetDir:     "./",
			fileSuffix:    "",
			objectType:    "github.com/khevse/codegen/tests/mainpkg.Factory",
			mockPackage:   "github.com/khevse/codegen/tests/mainpkg/mocks",
		}
		importList, spec, err := prepareObjectSpec(args)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(
			astpkg.ImportList{
				{Alias: "mainpkg", Path: "github.com/khevse/codegen/tests/mainpkg"},
				{Alias: "mocks", Path: "github.com/khevse/codegen/tests/mainpkg/mocks"},
			},
			importList,
			cmpopts.SortSlices(func(i, j astpkg.Import) bool {
				return i.Path < j.Path
			}),
		))
		require.Empty(t, cmp.Diff(
			&objectSpec{
				Name:    "FactoryWrapper",
				Comment: "FactoryWrapper wrapper for type IFactory: IFactory .",
				Methods: []methodSpec{
					{
						Name:    "NewObject1",
						Comment: "NewObject1 .",
						Params:  []field{},
						Results: []field{
							{
								FuncSpecName:   "_",
								TypeName:       "mainpkg.IObject1",
								ObjectSpecName: "IObject1",
								MockTypeName:   "IObject1Mock",
								MockPackage:    "mocks",
								Type: &astpkg.Ident{
									Package:     "mainpkg",
									PackagePath: "github.com/khevse/codegen/tests/mainpkg",
									Name:        "IObject1",
									Type: &astpkg.InterfaceType{
										Methods: []*astpkg.Field{
											{
												Name: "String",
												Type: &astpkg.FuncType{
													Params: []*astpkg.Field{},
													Results: []*astpkg.Field{
														{
															Name: "",
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
								},
							},
						},
					},
					{
						Name:    "NewObject2",
						Comment: "NewObject2 .",
						Params:  []field{},
						Results: []field{
							{
								FuncSpecName:   "_",
								TypeName:       "mainpkg.IObject2",
								ObjectSpecName: "IObject2",
								MockTypeName:   "IObject2Mock",
								MockPackage:    "mocks",
								Type: &astpkg.Ident{
									Package:     "mainpkg",
									PackagePath: "github.com/khevse/codegen/tests/mainpkg",
									Name:        "IObject2",
									Type: &astpkg.InterfaceType{
										Methods: []*astpkg.Field{
											{
												Name: "String",
												Type: &astpkg.FuncType{
													Params: []*astpkg.Field{},
													Results: []*astpkg.Field{
														{
															Name: "",
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
								},
							},
						},
					},
				},
				Fields: []objectSpecField{
					{
						Name:         "IObject1",
						TypeName:     "mainpkg.IObject1",
						MockPackage:  "mocks",
						MockTypeName: "IObject1Mock",
						Type: &astpkg.Ident{
							Package:     "mainpkg",
							PackagePath: "github.com/khevse/codegen/tests/mainpkg",
							Name:        "IObject1",
							Type: &astpkg.InterfaceType{
								Methods: []*astpkg.Field{
									{
										Name: "String",
										Type: &astpkg.FuncType{
											Params: []*astpkg.Field{},
											Results: []*astpkg.Field{
												{
													Name: "",
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
						},
					},
					{
						Name:         "IObject2",
						TypeName:     "mainpkg.IObject2",
						MockPackage:  "mocks",
						MockTypeName: "IObject2Mock",
						Type: &astpkg.Ident{
							Package:     "mainpkg",
							PackagePath: "github.com/khevse/codegen/tests/mainpkg",
							Name:        "IObject2",
							Type: &astpkg.InterfaceType{
								Methods: []*astpkg.Field{
									{
										Name: "String",
										Type: &astpkg.FuncType{
											Params: []*astpkg.Field{},
											Results: []*astpkg.Field{
												{
													Name: "",
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
						},
					},
				},
				BaseObjectTypeName: "*mainpkg.Factory",
			},
			spec,
			cmpopts.SortSlices(func(i, j methodSpec) bool {
				return i.Name < j.Name
			}),
		))
	})
}

func TestExecute(t *testing.T) {
	args := commandArgs{
		interfaceType: "github.com/khevse/codegen/tests/mainpkg.IFactory=FactoryWrapper",
		objectType:    "github.com/khevse/codegen/tests/mainpkg.Factory",
		targetDir:     "./",
		fileSuffix:    "_generated",
		mockPackage:   "github.com/khevse/codegen/tests/mainpkg/mocks",
	}
	require.NoError(t, (&Command{args: args}).Execute())

	const wantFile = "wrapper_generated.go"
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

package object_test_wrapper

import (
	minimock "github.com/gojuno/minimock/v3"
	mainpkg "github.com/khevse/codegen/tests/mainpkg"
	mocks "github.com/khevse/codegen/tests/mainpkg/mocks"
	"testing"
)

// FactoryWrapper mocks
type FactoryWrapperMocks struct {
	IObject1 *mocks.IObject1Mock
	IObject2 *mocks.IObject2Mock
}

// NewFactoryWrapperMocks return object FactoryWrapperMocks
func NewFactoryWrapperMocks(t *testing.T) *FactoryWrapperMocks {
	mc := minimock.NewController(t)

	return &FactoryWrapperMocks{
		IObject1: mocks.NewIObject1Mock(mc),
		IObject2: mocks.NewIObject2Mock(mc),
	}
}

/* FactoryWrapper wrapper for type IFactory: IFactory . */
type FactoryWrapper struct {
	mocks FactoryWrapperMocks
	base  *mainpkg.Factory
}

/* NewObject1 . */
func (w *FactoryWrapper) NewObject1() (_ mainpkg.IObject1) {
	existsMock := false ||
		w.mocks.IObject1 != nil
	if existsMock {
		return w.mocks.IObject1
	}

	return w.base.NewObject1()
}

/* NewObject2 . */
func (w *FactoryWrapper) NewObject2() (_ mainpkg.IObject2) {
	existsMock := false ||
		w.mocks.IObject2 != nil
	if existsMock {
		return w.mocks.IObject2
	}

	return w.base.NewObject2()
}

// FactoryWrapperBuilder wrapper builder
type FactoryWrapperBuilder struct {
	object FactoryWrapper
}

// SetBase set the base object with default behavior
func (b *FactoryWrapperBuilder) SetBase(val *mainpkg.Factory) *FactoryWrapperBuilder {
	b.object.base = val
	return b
}

// GetWrapper return wrapper object
func (b *FactoryWrapperBuilder) GetWrapper() *FactoryWrapper {
	return &b.object
}

// SetAllMocks set all mocks objects
func (b *FactoryWrapperBuilder) SetAllMocks(val *FactoryWrapperMocks) *FactoryWrapperBuilder {
	b.SetIObject1Mock(val)
	b.SetIObject2Mock(val)

	return b
}

// SetIObject1Mock set mock object
func (b *FactoryWrapperBuilder) SetIObject1Mock(val *FactoryWrapperMocks) *FactoryWrapperBuilder {
	b.object.mocks.IObject1 = val.IObject1
	return b
}

// SetIObject2Mock set mock object
func (b *FactoryWrapperBuilder) SetIObject2Mock(val *FactoryWrapperMocks) *FactoryWrapperBuilder {
	b.object.mocks.IObject2 = val.IObject2
	return b
}
`,
		string(data),
	)
}

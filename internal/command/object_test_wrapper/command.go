package object_test_wrapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/khevse/codegen/internal/pkg/astpkg"
	"github.com/khevse/codegen/internal/pkg/command"
	"github.com/samber/lo"
)

type commandArgs struct {
	interfaceType string
	objectType    string
	targetDir     string
	mockPackage   string
	fileSuffix    string
}

type Command struct {
	args commandArgs
}

func New() *Command {
	return new(Command)
}

func (c *Command) Name() string {
	return "object-test-wrapper"
}

func (c *Command) ShortName() string {
	return "w"
}

func (c *Command) InitFlags(flagSetter command.FlagSetter) error {
	const (
		flagInterfaceType = "interface-type"
		flagObjectType    = "object-type"
		flagTargetDir     = "target-dir"
		flagMockPackage   = "mock-package"
		flagFileSuffix    = "suffix"
	)

	flagSetter.Flags().StringVarP(
		&c.args.interfaceType,
		flagInterfaceType,
		"i",
		"",
		"interface type for mock generation. Examples: <package>.<InterfaceName>;<package>.<InterfaceName>=<MockName>",
	)
	flagSetter.Flags().StringVarP(
		&c.args.objectType,
		flagObjectType,
		"o",
		"",
		"object type which implement interface. Examples: <package>.<InterfaceName>",
	)
	flagSetter.Flags().StringVarP(
		&c.args.targetDir,
		flagTargetDir,
		"p",
		"",
		"target dir for result",
	)
	flagSetter.Flags().StringVarP(
		&c.args.mockPackage,
		flagMockPackage,
		"m",
		"",
		"mocks package",
	)
	flagSetter.Flags().StringVarP(
		&c.args.fileSuffix,
		flagFileSuffix,
		"",
		"",
		"result file suffix",
	)

	for _, flagName := range []string{flagInterfaceType, flagTargetDir, flagMockPackage} {
		if err := flagSetter.MarkFlagRequired(flagName); err != nil {
			return fmt.Errorf("mark flag as required(%s): %w", flagName, err)
		}
	}

	return nil
}

func (c *Command) Execute() error {
	targetDir, err := filepath.Abs(c.args.targetDir)
	if err != nil {
		return fmt.Errorf("get target dir full path: %w", err)
	}

	importList, objectSpec, err := prepareObjectSpec(c.args)
	if err != nil {
		return fmt.Errorf("prepare object specification: %w", err)
	}

	fileName := fmt.Sprintf("wrapper%s.go", c.args.fileSuffix)
	filePath := filepath.Join(targetDir, fileName)

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create file(%s): %w", filePath, err)
	}
	defer func() {
		f.Close()
	}()
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncate file: %w", err)
	}

	g := generator{
		Package:    filepath.Base(targetDir),
		Imports:    importList,
		ObjectSpec: *objectSpec,
	}

	if err := g.Generate(f); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}

func prepareObjectSpec(args commandArgs) (astpkg.ImportList, *objectSpec, error) {
	interfaceType, err := parseInterfaceType(args.interfaceType)
	if err != nil {
		return nil, nil, fmt.Errorf("parse interface type: %w", err)
	}

	objectType, err := parseObjectType(args.objectType)
	if err != nil {
		return nil, nil, fmt.Errorf("parse object type: %w", err)
	}

	pkg, err := astpkg.ParsePackage(interfaceType.Package)
	if err != nil {
		return nil, nil, fmt.Errorf("parse package(%s): %w", interfaceType.Package, err)
	}

	targetPackage, err := astpkg.GetPackagePath(args.targetDir)
	if err != nil {
		return nil, nil, fmt.Errorf("get target package: %w", err)
	}

	if err := astpkg.InitSelfPackageImports(targetPackage, pkg); err != nil {
		return nil, nil, fmt.Errorf("init self package imports: %w", err)
	}

	objectType.Package = lo.Ternary(
		objectType.Package == targetPackage,
		"",
		objectType.Package,
	)
	mockPackage := lo.Ternary(
		args.mockPackage == targetPackage,
		"",
		args.mockPackage,
	)
	baseImports := []string{"", mockPackage, objectType.Package}

	imports, err := astpkg.GetAllPackagesImports(baseImports, pkg)
	if err != nil {
		return nil, nil, fmt.Errorf("get all imports: %w", err)
	}

	typeDecl, ok := pkg.TypeDeclList.GetByName(interfaceType.TypeName)
	if !ok {
		return nil, nil, fmt.Errorf("not found type: %s", interfaceType.TypeName)
	}

	factoryDesc, err := newObjectSpec(interfaceType, objectType, mockPackage, typeDecl, imports)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"new factory description(%s): %w",
			interfaceType.WrapperName, err,
		)
	}

	clearImports := func() astpkg.ImportList {
		usedImports := make(map[string]struct{})
		addUsedImport := func(t astpkg.Type) {
			for _, item := range t.Imports() {
				usedImports[item.Alias] = struct{}{}
				usedImports[item.AliasFromPath()] = struct{}{}
			}
		}
		for _, item := range factoryDesc.Fields {
			addUsedImport(item.Type)
			usedImports[filepath.Base(item.MockPackage)] = struct{}{}
		}
		for _, method := range factoryDesc.Methods {
			for _, item := range method.Params {
				addUsedImport(item.Type)
			}
			for _, item := range method.Results {
				addUsedImport(item.Type)
			}
		}

		return lo.Filter(imports, func(item astpkg.Import, _ int) bool {
			_, used := usedImports[item.Alias]
			return item.Path != "" && item.Path != targetPackage && used
		})
	}

	return clearImports(), factoryDesc, nil
}

type argInterfaceType struct {
	Package     string
	TypeName    string
	WrapperName string
}

func parseInterfaceType(val string) (argInterfaceType, error) {
	val = strings.TrimSpace(val)

	fromTypeDelimiterIdx := strings.LastIndex(val, ".")
	if fromTypeDelimiterIdx == -1 {
		return argInterfaceType{}, fmt.Errorf("invalid type: %s", val)
	}

	packageName := val[:fromTypeDelimiterIdx]
	types := val[fromTypeDelimiterIdx+1:]

	typesParts := strings.Split(types, "=")
	typeName := typesParts[0]
	wrapperName := typesParts[0]
	if len(typesParts) > 1 {
		wrapperName = typesParts[1]
	}

	return argInterfaceType{
		Package:     packageName,
		TypeName:    typeName,
		WrapperName: wrapperName,
	}, nil
}

type argObjectType struct {
	Package  string
	TypeName string
}

func parseObjectType(val string) (argObjectType, error) {
	val = strings.TrimSpace(val)

	fromTypeDelimiterIdx := strings.LastIndex(val, ".")
	if fromTypeDelimiterIdx == -1 {
		return argObjectType{}, fmt.Errorf("invalid type: %s", val)
	}

	packageName := val[:fromTypeDelimiterIdx]
	typeName := val[fromTypeDelimiterIdx+1:]

	return argObjectType{
		Package:  packageName,
		TypeName: typeName,
	}, nil
}

package interface_creator

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
	fromType   string
	targetDir  string
	fileSuffix string
}

type Command struct {
	args commandArgs
}

func New() *Command {
	return new(Command)
}

func (c *Command) Name() string {
	return "interface"
}

func (c *Command) ShortName() string {
	return "i"
}

func (c *Command) InitFlags(flagSetter command.FlagSetter) error {
	const (
		flagFromType   = "type"
		flagTargetDir  = "target-dir"
		flagFileSuffix = "suffix"
	)

	flagSetter.Flags().StringVarP(
		&c.args.fromType,
		flagFromType,
		"t",
		"",
		"type for interface generation. Examples: <package>.<TypeName>; <package>.<TypeName>=<InterfaceName>; <package>.<TypeName1>=<InterfaceName1>,<package>.<TypeName2>=<InterfaceName2>",
	)
	flagSetter.Flags().StringVarP(
		&c.args.targetDir,
		flagTargetDir,
		"p",
		"",
		"target dir for the new interfaces",
	)
	flagSetter.Flags().StringVarP(
		&c.args.fileSuffix,
		flagFileSuffix,
		"",
		"",
		"result file suffix",
	)

	for _, flagName := range []string{flagFromType, flagTargetDir} {
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

	importList, objectSpecList, err := prepareObjectSpecList(c.args)
	if err != nil {
		return fmt.Errorf("prepare objects specifications: %w", err)
	}

	fileName := fmt.Sprintf("interfaces%s.go", c.args.fileSuffix)
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
		Interfaces: objectSpecList,
	}

	if err := g.Generate(f); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}

func prepareObjectSpecList(args commandArgs) (astpkg.ImportList, []objectSpec, error) {
	fromTypeList, err := parseFromType(args.fromType)
	if err != nil {
		return nil, nil, fmt.Errorf("parse types names: %w", err)
	}

	packageList, err := parsePackages(fromTypeList)
	if err != nil {
		return nil, nil, fmt.Errorf("parse packages: %w", err)
	}

	targetPackage, err := astpkg.GetPackagePath(args.targetDir)
	if err != nil {
		return nil, nil, fmt.Errorf("get target packages: %w", err)
	}

	if err := astpkg.InitSelfPackageImports(targetPackage, packageList...); err != nil {
		return nil, nil, fmt.Errorf("init self package imports: %w", err)
	}

	imports, err := astpkg.GetAllPackagesImports([]string{""}, packageList...)
	if err != nil {
		return nil, nil, fmt.Errorf("get all imports: %w", err)
	}

	interfaceList := make([]objectSpec, 0, len(fromTypeList))
	for _, pkg := range packageList {
		typeList := lo.Filter(fromTypeList, func(item argFromType, _ int) bool {
			return item.Package == pkg.Path
		})

		for _, item := range typeList {
			typeDecl, ok := pkg.TypeDeclList.GetByName(item.SourceName)
			if !ok {
				return nil, nil, fmt.Errorf("not found type: %s", item.SourceName)
			}

			methods := pkg.FuncDeclList.GetByReceiverName(item.SourceName)

			interfaceDesc, err := newObjectSpec(item.TargetName, typeDecl, methods, imports)
			if err != nil {
				return nil, nil, fmt.Errorf(
					"new object specification(%s): %w",
					item.TargetName, err,
				)
			}

			interfaceList = append(interfaceList, interfaceDesc)
		}
	}

	clearImports := func() astpkg.ImportList {
		usedImports := make(map[string]struct{})
		addUsedImport := func(t astpkg.Type) {
			for _, item := range t.Imports() {
				usedImports[item.Alias] = struct{}{}
				usedImports[item.AliasFromPath()] = struct{}{}
			}
		}
		for _, item := range interfaceList {
			for _, method := range item.Methods {
				for _, p := range method.Params {
					addUsedImport(p.Type)
				}
				for _, r := range method.Results {
					addUsedImport(r.Type)
				}
			}
		}

		return lo.Filter(imports, func(item astpkg.Import, _ int) bool {
			_, used := usedImports[item.Alias]
			return item.Path != "" && item.Path != targetPackage && used
		})
	}

	return clearImports(), interfaceList, nil
}

func parsePackages(fromTypeList []argFromType) ([]*astpkg.Package, error) {
	packagePathList := lo.Uniq(
		lo.Map(fromTypeList, func(item argFromType, _ int) string {
			return item.Package
		}),
	)

	packages := make([]*astpkg.Package, 0, len(packagePathList))
	for _, pkgPath := range packagePathList {
		pkg, err := astpkg.ParsePackage(pkgPath)
		if err != nil {
			return nil, fmt.Errorf("parse package(%s): %w", pkgPath, err)
		}

		packages = append(packages, pkg)
	}

	return packages, nil
}

type argFromType struct {
	Package    string
	SourceName string
	TargetName string
}

func parseFromType(val string) ([]argFromType, error) {
	list := make([]argFromType, 0)

	partList := strings.Split(val, ",")
	for _, part := range partList {
		part = strings.TrimSpace(part)

		fromTypeDelimiterIdx := strings.LastIndex(part, ".")
		if fromTypeDelimiterIdx == -1 {
			return nil, fmt.Errorf("invalid type: %s", part)
		}

		packageName := part[:fromTypeDelimiterIdx]
		types := part[fromTypeDelimiterIdx+1:]

		typesParts := strings.Split(types, "=")
		if len(typesParts) == 1 {
			list = append(list, argFromType{
				Package:    packageName,
				SourceName: typesParts[0],
				TargetName: typesParts[0],
			})
		} else {
			list = append(list, argFromType{
				Package:    packageName,
				SourceName: typesParts[0],
				TargetName: typesParts[1],
			})
		}
	}

	list = lo.Uniq(list)
	return list, nil
}

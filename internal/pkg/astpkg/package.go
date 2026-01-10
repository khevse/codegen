package astpkg

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

type Package struct {
	Path         string
	Dir          string
	TypeDeclList TypeDeclList
	FuncDeclList FuncDeclList
}

func ParsePackage(pkgName string) (*Package, error) {
	conf := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedImports |
			packages.LoadSyntax |
			packages.LoadAllSyntax |
			packages.NeedName,
	}

	pkgList, err := packages.Load(conf, pkgName)
	if err != nil {
		return nil, fmt.Errorf("load package:%w", err)
	}

	if err := getParsePackageError(pkgList); err != nil {
		return nil, err
	}

	var resPkg *Package

	for _, pkg := range pkgList {
		if pkg.ID != pkgName {
			continue
		}

		resPkg = &Package{
			Path:         pkg.ID,
			Dir:          pkg.Dir,
			TypeDeclList: nil,
			FuncDeclList: nil,
		}
		for _, file := range pkg.Syntax {
			importList := NewImportList(file.Imports)

			for _, decl := range file.Decls {
				switch castedDecl := decl.(type) {
				case *ast.GenDecl:
					if castedDecl.Tok == token.TYPE {
						for _, ts := range NewTypeDeclList(pkg.PkgPath, castedDecl) {
							err := InspectType(ts.Type, func(t Type) error {
								return SetPackageInformation(t, importList)
							})
							if err != nil {
								return nil, fmt.Errorf("set package information(%s): %w", ts, err)
							}

							resPkg.TypeDeclList = append(resPkg.TypeDeclList, ts)
						}
					}
				case *ast.FuncDecl:
					funcDecl := NewFuncDecl(castedDecl)
					err := InspectFuncDeclFields(
						funcDecl,
						func(f *Field) error {
							return InspectType(f.Type, func(t Type) error {
								return SetPackageInformation(t, importList)
							})
						},
					)
					if err != nil {
						return nil, fmt.Errorf("inspect func declaration(%s): %w", funcDecl, err)
					}

					resPkg.FuncDeclList = append(resPkg.FuncDeclList, funcDecl)
				}
			}
		}
	}
	if resPkg == nil {
		return nil, errors.New("not found")
	}

	return resPkg, nil
}

func GetPackagePath(pkgDir string) (string, error) {
	conf := &packages.Config{
		Mode: packages.NeedFiles,
		Dir:  pkgDir,
	}

	pkgList, err := packages.Load(conf)
	if err != nil {
		return "", fmt.Errorf("load package:%w", err)
	}

	if err := getParsePackageError(pkgList); err != nil {
		return "", err
	}

	for _, pkg := range pkgList {
		return pkg.ID, nil
	}

	return "", errors.New("package not found")
}

func SetPackagePathForAllDecl(pkg *Package) error {
	imp := NewImportWithAlias(pkg.Path)

	set := func(t Type) {
		casted, ok := t.(PackageCarrierType)
		if !ok || isBaseType(t) {
			return
		}

		required := casted.GetPackage() == "" && casted.GetPackagePath() == ""
		if required {
			casted.SetPackage(imp)
		}
	}

	for _, decl := range pkg.FuncDeclList {
		err := InspectFuncDeclFields(decl, func(f *Field) error {
			return InspectType(
				f.Type,
				func(t Type) error {
					set(t)
					return nil
				},
			)
		})
		if err != nil {
			return fmt.Errorf("inspect function declaration(%s): %w", decl, err)
		}
	}

	for _, decl := range pkg.TypeDeclList {
		err := InspectTypeDeclTypes(
			decl,
			func(t Type) error {
				set(t)
				return nil
			},
		)
		if err != nil {
			return fmt.Errorf("inspect type declaration(%s): %w", decl, err)
		}
	}

	return nil
}

func GetAllPackagesImports(baseImportPath []string, packageList ...*Package) (ImportList, error) {
	allImports := make([]string, len(baseImportPath))
	copy(allImports, baseImportPath)

	for _, pkg := range packageList {
		for _, decl := range pkg.FuncDeclList {
			imports, err := GetFuncDeclAllImportPath(decl)
			if err != nil {
				return nil, fmt.Errorf("get func declaration all imports(%s): %w", decl, err)
			}
			allImports = append(allImports, imports...)
		}
	}

	imports, err := NewImportListWithUniqAlias(allImports)
	if err != nil {
		return nil, fmt.Errorf("new imports list: %w", err)
	}

	return imports, nil
}

func InitSelfPackageImports(targetPackage string, packageList ...*Package) error {
	for _, pkg := range packageList {
		if targetPackage != "" && pkg.Path != targetPackage {
			if err := SetPackagePathForAllDecl(pkg); err != nil {
				return fmt.Errorf("set package path for all declaration: %w", err)
			}
		}
	}

	return nil
}

func getParsePackageError(pkgList []*packages.Package) error {
	var err error
	for _, pkg := range pkgList {
		if len(pkg.Errors) > 0 {
			err = errors.Join(
				err,
				fmt.Errorf(
					"package error(name=%s): %s",
					pkg.ID,
					lo.Map(pkg.Errors, func(item packages.Error, _ int) string {
						return fmt.Sprintf("%s(position:%s)", item.Msg, item.Pos)
					}),
				),
			)
		}
	}

	return err
}

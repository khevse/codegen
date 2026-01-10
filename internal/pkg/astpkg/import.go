package astpkg

import (
	"errors"
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

var ErrImportNotFound = errors.New("import not found")

type Import struct {
	Alias string
	Path  string
}

func (i Import) String() string { return fmt.Sprintf("%s(%s)", i.Alias, i.Path) }

func NewImport(alias, path string) Import {
	return Import{
		Alias: alias,
		Path:  path,
	}
}

func NewImportWithAlias(path string) Import {
	if path == "" {
		return NewImport("", "")
	}
	return NewImport(filepath.Base(path), path)
}

type ImportList []Import

func NewImportList(imports []*ast.ImportSpec) ImportList {
	importList := make(ImportList, 0, len(imports))
	for _, v := range imports {
		var alias string
		if name := v.Name; name != nil {
			alias = name.Name
		}

		pkgPath := strings.TrimFunc(
			v.Path.Value,
			func(r rune) bool { return r == '"' },
		)

		importList = append(
			importList,
			Import{Alias: alias, Path: pkgPath},
		)
	}

	return importList
}

var errNotUniqImport = errors.New("not uniq import")

func NewImportListWithUniqAlias(importPathList []string) (ImportList, error) {
	imports := make(ImportList, 0, len(importPathList))

	appendUniq := func(newItem Import) error {
		existsItem, err := imports.Get(newItem.Alias)
		if err != nil {
			if errors.Is(err, ErrImportNotFound) {
				imports = append(imports, newItem)
				return nil
			}

			return fmt.Errorf("get import by alias(%s): %w", newItem.Alias, err)
		}

		if exists := existsItem.Path == newItem.Path; exists {
			return nil
		}

		return errNotUniqImport
	}

	for _, path := range importPathList {
		newImport := NewImportWithAlias(path)

		if err := appendUniq(newImport); err != nil {
			if !errors.Is(err, errNotUniqImport) {
				return nil, fmt.Errorf("append import(%s): %w", newImport.String(), err)
			}

			const maxTries = 100
			sourceAlias := newImport.Alias
			for i := 1; i <= maxTries; i++ {
				newAlias := fmt.Sprintf("%s_%d", sourceAlias, i)
				newImport.Alias = newAlias

				if err := appendUniq(newImport); err != nil {
					if errors.Is(err, errNotUniqImport) && i < maxTries {
						continue
					}

					return nil, fmt.Errorf("append import(%s): %w", newImport.String(), err)
				}

				break
			}
		}
	}

	return imports, nil
}

func (l ImportList) Get(key string) (Import, error) {
	var val Import

	for _, item := range l {
		pkgName := filepath.Base(item.Path)
		if item.Alias == pkgName {
			item.Alias = ""
		}

		isEqualByPath := item.Alias == "" && pkgName == key
		isEqualByAlias := item.Alias == key

		if isEqualByPath || isEqualByAlias {
			isDuplicate := (val.Alias != "" && item.Path != val.Path) ||
				(val.Alias == "" && val.Path != "" && item.Path != val.Path)
			if isDuplicate {
				return Import{}, fmt.Errorf(
					"there are multiple imports with the same path: %+v; %+v",
					val, item,
				)
			}

			val = Import{
				Alias: lo.Ternary(item.Alias == pkgName, "", item.Alias),
				Path:  item.Path,
			}
		}
	}

	if val.Path == "" {
		return Import{}, ErrImportNotFound
	}

	return val, nil
}

func (l ImportList) GetByPath(path string) (Import, bool) {
	for _, item := range l {
		if item.Path == path {
			return item, true
		}
	}

	var empty Import
	return empty, false
}

package astpkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewImportList(t *testing.T) {
	newRawImports := func(t *testing.T, code string) []*ast.ImportSpec {
		f, err := parser.ParseFile(
			token.NewFileSet(),
			"",
			code,
			parser.DeclarationErrors,
		)
		require.NoError(t, err)
		return f.Imports
	}

	t.Run("empty", func(t *testing.T) {
		imports := newRawImports(t, `package p`)
		require.Empty(t, NewImportList(imports))
	})

	t.Run("success", func(t *testing.T) {
		imports := newRawImports(t, `package p; import p1 "./pkg1"; import "./pkg2"`)
		require.Equal(
			t,
			ImportList{
				{Alias: "p1", Path: "./pkg1"},
				{Alias: "", Path: "./pkg2"},
			},
			NewImportList(imports),
		)
	})
}

func TestImportListGet(t *testing.T) {
	t.Run("get by alias", func(t *testing.T) {
		list := ImportList{{Alias: "alias", Path: "example/pkg"}}
		item, err := list.Get("alias")
		require.NoError(t, err)
		require.Equal(t, Import{Alias: "alias", Path: "example/pkg"}, item)
	})

	t.Run("get by package name", func(t *testing.T) {
		list := ImportList{{Alias: "", Path: "example/pkg"}}
		item, err := list.Get("pkg")
		require.NoError(t, err)
		require.Equal(t, Import{Alias: "", Path: "example/pkg"}, item)
	})

	t.Run("two imports with some path", func(t *testing.T) {
		list := ImportList{{Alias: "pkg", Path: "example/pkg"}, {Alias: "", Path: "example/pkg"}}
		item, err := list.Get("pkg")
		require.NoError(t, err)
		require.Equal(t, Import{Alias: "", Path: "example/pkg"}, item)
	})

	t.Run("not found by package name", func(t *testing.T) {
		list := ImportList{{Alias: "", Path: "example/pkg"}}
		item, err := list.Get("unknown")
		require.ErrorIs(t, err, ErrImportNotFound)
		require.Empty(t, item)
	})

	t.Run("not found with alias", func(t *testing.T) {
		list := ImportList{{Alias: "alias", Path: "example/pkg"}}
		item, err := list.Get("unknown")
		require.ErrorIs(t, err, ErrImportNotFound)
		require.Empty(t, item)
	})

	t.Run("two imports with some path", func(t *testing.T) {
		list := ImportList{{Alias: "pkg", Path: "example/pkg1"}, {Alias: "pkg", Path: "example/pkg2"}}
		item, err := list.Get("pkg")
		require.EqualError(
			t,
			err,
			"there are multiple imports with the same path: pkg(example/pkg1); pkg(example/pkg2)",
		)
		require.Empty(t, item)
	})
}

func TestNewImportListWithUniqAlias(t *testing.T) {
	t.Parallel()

	t.Run("success for empty list", func(t *testing.T) {
		list, err := NewImportListWithUniqAlias([]string{})
		require.NoError(t, err)
		require.Empty(t, list)
	})

	t.Run("success with uniq imports", func(t *testing.T) {
		list, err := NewImportListWithUniqAlias(
			[]string{"./p1", "./p2"},
		)
		require.NoError(t, err)
		require.Equal(
			t,
			ImportList{
				{Alias: "p1", Path: "./p1"},
				{Alias: "p2", Path: "./p2"},
			},
			list,
		)
	})

	t.Run("success with not uniq imports", func(t *testing.T) {
		list, err := NewImportListWithUniqAlias(
			[]string{"./p1", "./p2/p1", "./p3/p2/p1"},
		)
		require.NoError(t, err)
		require.Equal(
			t,
			ImportList{
				{Alias: "p1", Path: "./p1"},
				{Alias: "p1_1", Path: "./p2/p1"},
				{Alias: "p1_2", Path: "./p3/p2/p1"},
			},
			list,
		)
	})

	t.Run("success with duplicates", func(t *testing.T) {
		list, err := NewImportListWithUniqAlias(
			[]string{"./p1", "./p1", "./p2/p1", "./p2/p1"},
		)
		require.NoError(t, err)
		require.Equal(
			t,
			ImportList{
				{Alias: "p1", Path: "./p1"},
				{Alias: "p1_1", Path: "./p2/p1"},
			},
			list,
		)
	})
}

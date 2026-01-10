package interface_creator

import (
	"embed"
	"io"
	"slices"
	"strings"

	"github.com/khevse/codegen/internal/pkg/application"
	"github.com/khevse/codegen/internal/pkg/astpkg"
	"github.com/khevse/codegen/internal/pkg/templatepkg"
)

//go:embed file.tmpl
var content embed.FS

type generator struct {
	Package    string
	Imports    astpkg.ImportList
	Interfaces []objectSpec
}

func (g generator) Generate(w io.Writer) error {
	slices.SortFunc(g.Imports, func(i, j astpkg.Import) int {
		return strings.Compare(i.Path, j.Path)
	})
	slices.SortFunc(g.Interfaces, func(i, j objectSpec) int {
		return strings.Compare(i.Name, j.Name)
	})
	for _, item := range g.Interfaces {
		slices.SortFunc(item.Methods, func(i, j methodSpec) int {
			return strings.Compare(i.Name, j.Name)
		})
	}

	params := templatepkg.ExecuteTemplateParams{
		Writer:       w,
		FS:           content,
		TemplateFile: "file.tmpl",
		Data: map[string]any{
			"package":    g.Package,
			"imports":    g.Imports,
			"interfaces": g.Interfaces,
			"appInfo":    application.GetInfo(),
		},
		Format: true,
	}

	return templatepkg.ExecuteTemplate(params)
}

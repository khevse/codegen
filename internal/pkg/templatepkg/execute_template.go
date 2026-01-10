package templatepkg

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"text/template"
)

type ExecuteTemplateParams struct {
	Writer       io.Writer
	FS           fs.FS
	TemplateFile string
	Data         any
	Format       bool
}

func ExecuteTemplate(params ExecuteTemplateParams) error {
	tmpl, err := template.New("").ParseFS(params.FS, params.TemplateFile)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.ExecuteTemplate(buf, params.TemplateFile, params.Data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	if !params.Format {
		if _, err = io.Copy(params.Writer, buf); err != nil {
			return fmt.Errorf("write code: %w", err)
		}
		return nil
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formate code: %w", err)
	}

	_, err = io.Copy(params.Writer, bytes.NewReader(formattedCode))
	if err != nil {
		return fmt.Errorf("write formatted code: %w", err)
	}

	return nil
}

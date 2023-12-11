package templates

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed *
var TemplateFS embed.FS

var parsed *template.Template

func Parse() (*template.Template, error) {
	if parsed != nil {
		return parsed, nil
	}

	tmpl := template.New("")

	err := fs.WalkDir(TemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".gohtml" {
			return nil
		}

		// Read the template file
		file, err := TemplateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Convert the path to template name
		// This will remove any preceding './' and replace '/' with '.'
		name := strings.TrimPrefix(path, "./")
		//name = strings.ReplaceAll(name, "/", ".")
		name = strings.TrimSuffix(name, ".gohtml")

		// Parse and add to the template set
		_, err = tmpl.New(name).Parse(string(file))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	parsed = tmpl

	return tmpl, nil
}

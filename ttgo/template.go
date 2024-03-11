package ttgo

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

type Template struct {
	file string
	tmpl *template.Template
}

func NewTemplate(file string) (*Template, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	abs = filepath.Clean(abs)

	p, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	te := &Template{
		file: abs,
	}

	tmpl := template.New(filepath.Base(abs))
	tmpl = tmpl.Funcs(te.FuncMap())
	tmpl = tmpl.Funcs(sprig.TxtFuncMap())
	tmpl, err = tmpl.Parse(string(p))
	if err != nil {
		return nil, err
	}
	te.tmpl = tmpl

	return te, nil
}

func (te *Template) FuncMap() map[string]any {
	return map[string]any{
		"partial":   te.partial,
		"writeFile": writeFile,
		"fromYaml":  fromYaml,
		"toYaml":    toYaml,
	}
}

func (te *Template) Execute(data any) (string, error) {
	b := new(bytes.Buffer)
	err := te.tmpl.Execute(b, data)
	return b.String(), err
}

func (te *Template) partial(file string, data any) (string, error) {
	abs := filepath.Clean(filepath.Join(filepath.Dir(te.file), file))
	pa, err := NewTemplate(abs)
	if err != nil {
		return "", err
	}
	return pa.Execute(data)
}

func writeFile(file, content string) (string, error) {
	return "", os.WriteFile(file, []byte(content), 0644)
}

func fromYaml(s string) (any, error) {
	var v any
	err := yaml.Unmarshal([]byte(s), &v)
	return v, err
}

func toYaml(v any) (string, error) {
	p, err := yaml.Marshal(v)
	return string(p), err
}

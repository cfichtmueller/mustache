package mustache

import (
	"fmt"
	"html/template"
)

var (
	ErrTemplateNotFound = fmt.Errorf("template not found")
)

type Engine struct {
	templates map[string]*Template
	partials  *StaticProvider
}

func NewEngine() *Engine {
	return &Engine{
		templates: make(map[string]*Template),
		partials: &StaticProvider{
			Partials: make(map[string]string),
		},
	}
}

func (e *Engine) Parse(name, t string) error {
	tmpl := &Template{t, "{{", "}}", 0, 1, []interface{}{}, false, e.partials, template.HTMLEscapeString}
	err := tmpl.parse()

	if err != nil {
		return err
	}

	e.templates[name] = tmpl

	return nil
}

func (e *Engine) AddPartial(name, t string) error {
	e.partials.Partials[name] = t
	return nil
}

func (e *Engine) Render(name string, context ...interface{}) (string, error) {
	t, ok := e.templates[name]
	if !ok {
		return "", ErrTemplateNotFound
	}
	return t.Render(context...)
}

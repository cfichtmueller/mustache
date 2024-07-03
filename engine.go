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

func (e *Engine) Parse(name, data string) error {
	tmpl, err := e.parseTemplate(data)
	if err != nil {
		return err
	}

	e.templates[name] = tmpl

	return nil
}

func (e *Engine) parseTemplate(data string) (*Template, error) {
	tmpl := &Template{
		data:              data,
		otag:              "{{",
		ctag:              "}}",
		p:                 0,
		curline:           1,
		elems:             []interface{}{},
		forceRaw:          false,
		partial:           e.partials,
		escape:            template.HTMLEscapeString,
		parserFunc:        e.parseTemplate,
		partialParserFunc: e.parsePartial,
	}
	err := tmpl.parse()

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (e *Engine) parsePartial(data string, partials PartialProvider) (*Template, error) {
	tmpl, err := e.parseTemplate(data)
	if err != nil {
		return nil, err
	}
	tmpl.partial = partials

	return tmpl, nil
}

func (e *Engine) MustParse(name, data string) *Engine {
	if err := e.Parse(name, data); err != nil {
		panic(err)
	}
	return e
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

func (e *Engine) RenderInLayout(name, layout string, context ...interface{}) (string, error) {
	tmpl, err := e.GetTemplate(name)
	if err != nil {
		return "", err
	}
	layoutTmpl, err := e.GetTemplate(layout)
	if err != nil {
		return "", err
	}

	return tmpl.RenderInLayout(layoutTmpl, context...)
}

func (e *Engine) GetTemplate(name string) (*Template, error) {
	t, ok := e.templates[name]
	if !ok {
		return nil, ErrTemplateNotFound
	}
	return t, nil
}

func (e *Engine) MustGetTemplate(name string) *Template {
	t, err := e.GetTemplate(name)
	if err != nil {
		panic(err)
	}
	return t
}

package mustache

import (
	"fmt"
	"html/template"
	"io/fs"
	"path"
)

var (
	ErrTemplateNotFound = fmt.Errorf("template not found")
)

type Engine struct {
	// AllowMissingVariables defines the behavior for a variable "miss." If it
	// is true (the default), an empty string is emitted. If it is false, an error
	// is generated instead.
	AllowMissingVariables bool
	templates             map[string]*Template
	partials              map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		AllowMissingVariables: true,
		templates:             make(map[string]*Template),
		partials:              make(map[string]string),
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

// ParseFS parses the template definitions from the named files. The template
// names will be the file names without path or extension. There must be at least
// one file. If an error occurs, parsing stops and an error is returned.
//
// It accepts a list of glob patterns.
// (Note that most file names serve as glob patterns matching only themselves.)
func (e *Engine) ParseFS(fs fs.FS, patterns ...string) error {
	return handleFS(e, fs, patterns, parseTemplate)
}

// MustParseFS is like ParseFS but will panic when an error occurs.
func (e *Engine) MustParseFS(fs fs.FS, patterns ...string) *Engine {
	if err := e.ParseFS(fs, patterns...); err != nil {
		panic(err)
	}
	return e
}

func handleFS(e *Engine, fsys fs.FS, patterns []string, handler func(e *Engine, name, content string) error) error {
	var filenames []string
	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return fmt.Errorf("cfichtmueller.com/mustache: pattern matches no files: %#q", pattern)
		}
		filenames = append(filenames, list...)
	}

	return handleFiles(e, readFileFS(fsys), filenames, handler)
}

func handleFiles(e *Engine, readFile func(string) (string, string, error), filenames []string, handler func(e *Engine, name, content string) error) error {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return fmt.Errorf("cfichtmueller.com/mustache: no files named in call to handleFS")
	}

	for _, filename := range filenames {
		name, content, err := readFile(filename)
		if err != nil {
			return err
		}
		if err := handler(e, name, content); err != nil {
			return err
		}
	}

	return nil
}

func readFileFS(fsys fs.FS) func(string) (string, string, error) {
	return func(file string) (name, content string, err error) {
		name = path.Base(file)
		b, err := fs.ReadFile(fsys, file)
		if err != nil {
			return "", "", err
		}
		return name[:len(name)-len(path.Ext(name))], string(b), nil
	}
}

func parseTemplate(e *Engine, name, content string) error {
	return e.Parse(name, content)
}

func (e *Engine) parseTemplate(data string) (*Template, error) {
	tmpl := &Template{
		allowMissingVariables: e.AllowMissingVariables,
		data:                  data,
		otag:                  "{{",
		ctag:                  "}}",
		p:                     0,
		curline:               1,
		elems:                 []interface{}{},
		forceRaw:              false,
		partial:               e,
		escape:                template.HTMLEscapeString,
		parserFunc:            e.parseTemplate,
		partialParserFunc:     e.parsePartial,
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
	e.partials[name] = t
	return nil
}

// AddPartialFS adds the partial definitions from the named files. The partial
// names will be the file names without path or extension. There must be at least
// one file. If an error occurs, parsing stops and an error is returned.
//
// It accepts a list of glob patterns.
// (Note that most file names serve as glob patterns matching only themselves.)
func (e *Engine) AddPartialFS(fs fs.FS, patterns ...string) error {
	return handleFS(e, fs, patterns, addPartial)
}

// MustAddPartialsFS is like AddPartialFS but will panic when an error occurs.
func (e *Engine) MustAddPartialFS(fs fs.FS, patterns ...string) *Engine {
	if err := e.AddPartialFS(fs, patterns...); err != nil {
		panic(err)
	}
	return e
}

func addPartial(e *Engine, name, content string) error {
	e.AddPartial(name, content)
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

func (e *Engine) GetPartial(name string) (string, error) {
	p, ok := e.partials[name]
	if !ok {
		return "", nil
	}
	return p, nil
}

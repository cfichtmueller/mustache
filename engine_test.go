package mustache

import (
	"embed"
	"testing"
)

func TestEngineParseTemplate(t *testing.T) {
	e := NewEngine()

	mustParse(t, e, "hello", "Hello {{Name}}")
	mustParse(t, e, "weather", "Sunny")
	mustParse(t, e, "report", "Hello {{Name}}, the weather is {{how}}")

	mustRender(t, e, "hello", Name{Name: "John"}, "Hello John")
	mustRender(t, e, "weather", nil, "Sunny")

	s, err := e.Render("report", Name{Name: "Alice"}, map[string]string{"how": "sunny"})
	if err != nil {
		t.Errorf("failed to render template: %v", err)
	}
	expected := "Hello Alice, the weather is sunny"

	if s != expected {
		t.Errorf("expected '%s', got '%s'", expected, s)
	}
}

func TestEngineRenderUnknownTemplate(t *testing.T) {
	e := NewEngine()

	s, err := e.Render("unknown", Name{Name: "Bob"})
	if err != ErrTemplateNotFound {
		t.Errorf("expected ErrTemplateNotFound, got %v", err)
	}
	if s != "" {
		t.Errorf("expected empty render result, got '%s'", s)
	}
}

func TestRenderPartials(t *testing.T) {
	e := NewEngine()

	e.AddPartial("p1", "<strong>{{Name}}</strong>")
	mustParse(t, e, "hello", "Hello {{>p1}}")

	mustRender(t, e, "hello", Name{Name: "Bob"}, "Hello <strong>Bob</strong>")
}

func TestAddPartialsAfterTemplates(t *testing.T) {
	e := NewEngine()

	e.MustParse("t", "Hello {{>p1}}")
	e.AddPartial("p1", "<strong>{{Name}}</strong>")

	mustRender(t, e, "t", Name{Name: "Bob"}, "Hello <strong>Bob</strong>")
}

func TestEngineIsolation(t *testing.T) {
	e1 := NewEngine()
	e2 := NewEngine()

	e1.Parse("t", "Hello from {{>g}}")
	e1.AddPartial("g", "<strong>{{Name}}</strong>")

	e2.Parse("t", "Guten Tag von {{>g}}")
	e2.AddPartial("g", "<i>{{Name}}</i>")

	n := Name{Name: "Alice"}

	mustRender(t, e1, "t", n, "Hello from <strong>Alice</strong>")
	mustRender(t, e2, "t", n, "Guten Tag von <i>Alice</i>")
}

var (
	//go:embed tests
	testFS embed.FS
)

func TestEngineParseFS(t *testing.T) {
	e := NewEngine()
	if err := e.ParseFS(testFS, "tests/templates/*.mustache"); err != nil {
		t.Error(err)
	}
	for _, name := range []string{"t1", "t2", "t3"} {
		tmpl, err := e.GetTemplate(name)
		if err != nil {
			t.Errorf("expected to find template %s: %v", name, err)
		}
		if tmpl == nil {
			t.Errorf("expected to find template %s but got nil", name)
		}
	}
}

func TestEngineAddPartialFS(t *testing.T) {
	e := NewEngine()
	if err := e.AddPartialFS(testFS, "tests/partials/*.mustache"); err != nil {
		t.Error(err)
	}
	for _, name := range []string{"p1", "p2"} {
		p, err := e.GetPartial(name)
		if err != nil {
			t.Errorf("expected to find partial %s: %v", name, err)
		}
		if p == "" {
			t.Errorf("expected to find partial %s but got empty", name)
		}
	}
}

func TestDecapitalizeStructNames(t *testing.T) {
	user := &User{Name: "Bob"}
	e := NewEngine().MustParse("field", "Hello {{name}}!").MustParse("method", "Hello {{initial}}!")

	mustRender(t, e, "field", user, "Hello Bob!")
	mustRender(t, e, "method", user, "Hello B!")

	e = NewEngine()
	e.DecapitalizeStructFieldNames = false
	e.MustParse("field", "Hello {{name}}!").MustParse("method", "Hello {{initial}}!")

	mustRender(t, e, "field", user, "Hello !")
	mustRender(t, e, "method", user, "Hello !")
}

func mustParse(t *testing.T, e *Engine, name, tmpl string) {
	if err := e.Parse(name, tmpl); err != nil {
		t.Errorf("failed to parse template: %v", err)
	}
}

func mustRender(t *testing.T, e *Engine, name string, context interface{}, expected string) {
	s, err := e.Render(name, context)
	if err != nil {
		t.Errorf("failed to render template: %v", err)
	}

	if s != expected {
		t.Errorf("Expected '%s', got '%s'", expected, s)
	}
}

type Name struct {
	Name string
}

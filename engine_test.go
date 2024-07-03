package mustache

import (
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

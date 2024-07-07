# mustache templates in Go

![Go Test](https://github.com/cfichtmueller/mustache/actions/workflows/testing.yml/badge.svg)

## Installation

To install mustache, run `go get github.com/cfichtmueller/mustache` and use `import "github.com/cfichtmueller/mustache"`.

## Usage

Usage of this library is pretty simple. It involves the following steps:

1. create a mustache engine for rendering
2. add partials (optional)
3. add templates
4. start rendering


Simplest Use Case

```go 
engine := mustache.NewEngine()

err := engine.Parse("greeting", "Hello {{name}}")

result, err := engine.Render("greeting", map[string]string{"name": "Bob"})
```

With Partials

```go 
engine := mustache.NewEngine()

engine.AddPartial("user", "<strong>{{name}}</strong>")

err := engine.Parse("greeting", "Hello {{> user}}")

result, err := engine.Render("greeting", map[string]string{"name": "Bob"})
```

Using Layouts

```go 
engine := mustache.NewEngine()

err := engine.Parse("layout", "<div>{{content}}</div>")

err = engine.Parse("greeting", "Hello {{name}}")

result, err := engine.RenderInLayout("greeting", "layout", map[string]string{"name": "Bob"})
```

You can also load templates and partials from a file system. E.g. you want to bundle up all your html templates within your app.


```ascii
├─templates/
│ ├─templates/
│ │ ├─index.mustache
│ │ └─profile.mustache
│ └─partials/
│   ├─user.mustache
│   └─project.mustache
└─main.go
```

```go
var (
    //go:embed templates
    templateFS embed.FS
    engine = mustache.
        NewEngine().
        MustParseFS(templateFS, "templates/templates/*.mustache").
        MustAddPartialFS(templateFS, "templates/partials/*.mustache")
)
```

## Decapitalize Struct Field Names

Go is unique in the sense as it uses casing for determining visibility.To shield template authors from the intricacies of Go, the engine automatically decapitalizes field and method names of structs when rendering a template.

This feature can be disabled.

```go

type User struct {
    Name string
}

func (u *User) Initial() string {
    return u.Name[:1]
}

e := NewEngine()
e.DecapitalizeStructFieldNames = true // default behavior
e.MustCompile("field", "Hello {{name}}!")
e.MustCompile("method", "Hello {{initial}}!")

val, err := e.Render("field", &User{Name: "Bob"}) // Hello Bob!
val, err := e.Render("method", &User{Name: "Alice"}) // Hello A!

e.DecapitalizeStructFieldNames = false // disable decapitalization
e.MustCompile("field", "Hello {{name}}!")
e.MustCompile("method", "Hello {{initial}}!")

val, err := e.Render("field", &User{Name: "Bob"}) // Hello !
val, err := e.Render("method", &User{Name: "Alice"}) // Hello !
```

## Escaping

mustache.go follows the official mustache HTML escaping rules. That is, if you enclose a variable with two curly brackets, `{{var}}`, the contents are HTML-escaped. For instance, strings like `5 > 2` are converted to `5 &gt; 2`. To use raw characters, use three curly brackets `{{{var}}}`.

## Layouts

It is a common pattern to include a template file as a "wrapper" for other templates. The wrapper may include a header and a footer, for instance. Mustache.go supports this pattern with the following two methods:

```go
engine.MustParse("template", templateString)
engine.MustParse("layout", layoutString)

res, err := engine.RenderInLayout("template", "layout", context)

```

The layout template must have a variable called `{{content}}`. For example, given the following layout:

```html
<html>
<head><title>Hi</title></head>
<body>
{{{content}}}
</body>
</html>
```

and the following template:

```html
<h1>Hello World!</h1>
```

A call to `engine.RenderInLayout("template", "layout", nil)` will produce:

```html
<html>
<head><title>Hi</title></head>
<body>
<h1>Hello World!</h1>
</body>
</html>
```

## A note about method receivers

Mustache.go supports calling methods on objects, but you have to be aware of Go's limitations. For example, lets's say you have the following type:

```go
type Person struct {
    FirstName string
    LastName string
}

func (p *Person) Name1() string {
    return p.FirstName + " " + p.LastName
}

func (p Person) Name2() string {
    return p.FirstName + " " + p.LastName
}
```

While they appear to be identical methods, `Name1` has a pointer receiver, and `Name2` has a value receiver. Objects of type `Person`(non-pointer) can only access `Name2`, while objects of type `*Person`(person) can access both. This is by design in the Go language.

So if you write the following:

```go
engine := mustache.NewEngine().MustParse("t", "{{Name1}}")

engine.Render("t", Person{"John", "Smith"})
```

It'll be blank. You either have to use `&Person{"John", "Smith"}`, or call `Name2`

## Supported features

- Variables
- Comments
- Change delimiter
- Sections (boolean, enumerable, and inverted)
- Partials

## Further Reading

- mustache project page
- mustache manual
- example mustache files

## Credits

Kudeos to [hoisie](https://github.com/hoisie), [cbroglie](https://github.com/cbroglie) and all other contributors for their original work on this project.
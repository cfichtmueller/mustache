# Mustache Template Engine for Go

![Go Test](https://github.com/cfichtmueller/mustache/actions/workflows/testing.yml/badge.svg)

----

## Why a Fork?

I forked [cbroglie/mustache](https://github.com/cbroglie/mustache) because I wanted to change some things:

- introduce the concept of engines which isolate template and execution contexts
- remove all implicit file system access

----

## Package Overview

This library is an implementation of the Mustache template language in Go.

### Mustache Spec Compliance

[mustache/spec](https://github.com/mustache/spec) contains the formal standard for Mustache, and it is included as a submodule (using v1.2.1) for testing compliance. All of the tests pass (big thanks to [kei10in](https://github.com/kei10in)), with the exception of the null interpolation tests added in v1.2.1. There is experimental support for a subset of the optional lambda functionality (thanks to [fromhut](https://github.com/fromhut)). The optional inheritance functionality has not been implemented.

----

## Documentation

For more information about mustache, check out the [mustache project page](https://github.com/mustache/mustache) or the [mustache manual](https://mustache.github.io/mustache.5.html).

Also check out some [example mustache files](http://github.com/mustache/mustache/tree/master/examples/).

----

## Installation

To use it in a program, run `go get github.com/cfichtmueller/mustache` and use `import "github.com/cfichtmueller/mustache"`.

----

## Usage

Usage of this library is pretty simple. It involves the following steps:

1. create a mustache engine for rendering
2. add partials (optional)
3. add templates
4. start rendering


```go 
engine := mustache.NewEngine()
engine.AddPartial("user", "<strong>{{name}}</strong>")
err := engine.Parse("base", "{{#names}}{{> user}}{{/names}}")

result, err := engine.Render("base", map[string]interface{}{"names": []map[string]string{{"name": "Alice"}, {"name": "Bob"}}})
```

----

## Escaping

mustache.go follows the official mustache HTML escaping rules. That is, if you enclose a variable with two curly brackets, `{{var}}`, the contents are HTML-escaped. For instance, strings like `5 > 2` are converted to `5 &gt; 2`. To use raw characters, use three curly brackets `{{{var}}}`.

----

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

----

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

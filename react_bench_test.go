package benchmarks_test

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/mark-rushakoff/go-reactjs-benchmarks/engine"
)

var (
	// Package-level private variable to store benchmark output,
	// to ensure compiler doesn't optimize anything away.
	resultString string
	resultEngine engine.Engine

	// Engines that we will initialize and then clone for test.
	otto013 engine.Engine
	otto014 engine.Engine
)

const (
	render013Div        = `React.renderToString(React.createElement("div", null, "Hello world"))`
	render013Components = `React.renderToString(React.createElement(ParentComponent))`

	render014Div        = `ReactDOMServer.renderToString(React.createElement("div", null, "Hello world"))`
	render014Components = `ReactDOMServer.renderToString(React.createElement(ParentComponent))`
)

var (
	renderDivRegexp = regexp.MustCompile(`^<div.*>Hello world</div>$`)
)

func init() {
	initOtto13()
	// initOtto14()
}

func initOtto13() {
	otto013 = engine.NewOttoEngine()
	if err := otto013.Load(react0133Source()); err != nil {
		panic(err)
	}
}

func initOtto14() {
	otto014 = engine.NewOttoEngine()
	if err := otto014.Load(react0143Source()); err != nil {
		panic(err)
	}
	if err := otto014.Load(reactDOMServer0143Source()); err != nil {
		panic(err)
	}
}

func BenchmarkOttoLoadReact013(b *testing.B) {
	src := react0133Source()
	b.ResetTimer()

	var e engine.Engine
	for n := 0; n < b.N; n++ {
		e = engine.NewOttoEngine()
		_ = e.Load(src)
	}
	resultEngine = e
}

func BenchmarkOttoLoadReact014(b *testing.B) {
	b.Skip("otto cannot yet handle React 0.14")
	reactSrc := react0143Source()
	domServerSrc := reactDOMServer0143Source()
	b.ResetTimer()

	var e engine.Engine
	for n := 0; n < b.N; n++ {
		e = engine.NewOttoEngine()
		_ = e.Load(reactSrc)
		_ = e.Load(domServerSrc)
	}
	resultEngine = e
}

func TestOttoReact013RenderDiv(t *testing.T) {
	e := otto013.Clone()
	testRenderDiv(t, e, render013Div)
}

func TestOttoReact014RenderDiv(t *testing.T) {
	t.Skip("otto cannot yet handle React 0.14")
	e := otto014.Clone()
	testRenderDiv(t, e, render014Div)
}

func testRenderDiv(t *testing.T, e engine.Engine, script string) {
	s, err := e.RunReact(script)
	if err != nil {
		t.Error("Error executing script:", err)
	}

	if !renderDivRegexp.MatchString(s) {
		t.Error("Expected result of render to roughly match <div>Hello world</div> but received: " + s)
	}
}

func BenchmarkOttoReact013RunRenderDiv(b *testing.B) {
	e := otto013.Clone()
	b.ResetTimer()

	benchmarkRunReact(b, e, render013Div)
}

func BenchmarkOttoReact014RunRenderDiv(b *testing.B) {
	b.Skip("otto cannot yet handle React 0.14")
	e := otto014.Clone()
	b.ResetTimer()

	benchmarkRunReact(b, e, render014Div)
}

func TestOttoReact013RenderComponents(t *testing.T) {
	e := otto013.Clone()
	if err := e.Load(componentSource()); err != nil {
		t.Error("Error executing `components.js` script file", err)
	}

	testRenderComponents(t, e, render013Components)
}

func TestOttoReact014RenderComponents(t *testing.T) {
	t.Skip("otto cannot yet handle React 0.14")
	e := otto014.Clone()
	if err := e.Load(componentSource()); err != nil {
		t.Error("Error executing `components.js` script file", err)
	}

	testRenderComponents(t, e, render014Components)
}

func testRenderComponents(t *testing.T, e engine.Engine, script string) {
	s, err := e.RunReact(script)
	if err != nil {
		t.Error("Error executing script:", err)
	}

	expectedText := []string{
		"the-parent",
		"child-1",
		"child-2",
		"child-3",
		"First child",
		"Second child",
		"Third child",
	}
	for _, exp := range expectedText {
		if !strings.Contains(s, exp) {
			t.Error("Expected result of rendering components to contain text " + exp + " but it did not.")
		}
	}
}

func BenchmarkOttoReact013RenderComponents(b *testing.B) {
	e := otto013.Clone()
	_ = e.Load(componentSource())
	b.ResetTimer()

	benchmarkRunReact(b, e, render013Components)
}

func BenchmarkOttoReact014RenderComponents(b *testing.B) {
	b.Skip("otto cannot yet handle React 0.14")
	e := otto014.Clone()
	_ = e.Load(componentSource())
	b.ResetTimer()

	benchmarkRunReact(b, e, render014Components)
}

func benchmarkRunReact(b *testing.B, e engine.Engine, script string) {
	var s string

	for n := 0; n < b.N; n++ {
		s, _ = e.RunReact(script)
	}
	resultString = s
}

func react0133Source() []byte {
	return mustLoadSource("assets/vendor/react-0.13.3.min.js")
}

func react0143Source() []byte {
	return mustLoadSource("assets/vendor/react-0.14.3.min.js")
}

func reactDOMServer0143Source() []byte {
	return mustLoadSource("assets/vendor/react-dom-server-0.14.3.min.js")
}

func componentSource() []byte {
	return mustLoadSource("assets/components.js")
}

func mustLoadSource(path string) []byte {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return src
}

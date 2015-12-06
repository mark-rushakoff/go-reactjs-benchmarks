package benchmarks_test

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

var (
	// Package-level private variable to store benchmark output,
	// to ensure compiler doesn't optimize anything away.
	resultString string
	resultOtto   *otto.Otto

	// Engines that we will initialize and then clone for test.
	otto013 *otto.Otto
)

const (
	render013Div        = `React.renderToString(React.createElement("div", null, "Hello world"))`
	render013Components = `React.renderToString(React.createElement(ParentComponent))`
)

var (
	renderDivRegexp = regexp.MustCompile(`^<div.*>Hello world</div>$`)
)

func init() {
	otto013 = otto.New()
	_, err := otto013.Run(react0133Source())
	if err != nil {
		panic(err)
	}
}

// Read the React source from disk once; then create a new interpreter and
// read/parse/execute the React entry point.
func BenchmarkOttoLoadReact0133(b *testing.B) {
	src := react0133Source()
	b.ResetTimer()

	var o *otto.Otto
	for n := 0; n < b.N; n++ {
		o = otto.New()
		_, _ = o.Run(src)
	}
	resultOtto = o
}

func TestOttoReact0133RenderDiv(t *testing.T) {
	o := otto013.Copy()
	v, err := o.Run(render013Div)
	if err != nil {
		t.Error("Error executing `render013Div` script", err)
	}
	if v.Class() == "string" {
		t.Error("Expected string result, actual type is: " + v.Class())
	}

	s := v.String()
	if !renderDivRegexp.MatchString(s) {
		t.Error("Expected result of render to roughly match <div>Hello world</div> but received: " + s)
	}
}

func BenchmarkOttoReact0133RunRenderDiv(b *testing.B) {
	o := otto013.Copy()
	b.ResetTimer()

	var s string
	for n := 0; n < b.N; n++ {
		v, _ := o.Run(render013Div)
		s = v.String()
	}
	resultString = s
}

func TestOttoReact0133RenderComponents(t *testing.T) {
	o := otto013.Copy()
	_, err := o.Run(componentSource())
	if err != nil {
		t.Error("Error executing `components.js` script file", err)
	}

	v, err := o.Run(render013Components)
	if err != nil {
		t.Error("Error executing `render013Components` script", err)
	}

	if v.Class() == "string" {
		t.Error("Expected string result, actual type is: " + v.Class())
	}

	s := v.String()

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
			t.Error("Expected result of componets to contain text " + exp + " but it did not.")
		}
	}
}

func BenchmarkOttoReact0133RenderComponents(b *testing.B) {
	o := otto013.Copy()
	_, _ = o.Run(componentSource())
	b.ResetTimer()

	var s string

	for n := 0; n < b.N; n++ {
		v, _ := o.Run(render013Components)
		s = v.String()
	}
	resultString = s
}

func react0133Source() []byte {
	return mustLoadSource("assets/vendor/react-0.13.3.min.js")
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

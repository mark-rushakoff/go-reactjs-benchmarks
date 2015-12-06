# go-reactjs-benchmarks

Benchmarks of different combinations of Go JS engines and different versions of React.js, primarily focusing on `React.renderToString`, to learn about server-side React.js rendering in Go.

Currently testing with:

* React versions:
  * 0.13.3
  * 0.14.3
* JS engines:
  * [Otto](https://github.com/robertkrimen/otto)

## Run the benchmarks

From the root folder (e.g. from `$GOPATH/github.com/mark-rushakoff/go-reactjs-benchmarks`), run `go test -bench=.`

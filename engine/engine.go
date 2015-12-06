package engine

type Engine interface {
	// Returns a copy of the engine, so you can do things without modifying the original engine's state.
	Clone() Engine

	// Parse the given source, updating the engine's internal state.
	// Assumes source contains no relevant output.
	Load(src []byte) error

	// Run the given React snippet.
	// Assumes the snippet results in string output.
	RunReact(src string) (string, error)
}

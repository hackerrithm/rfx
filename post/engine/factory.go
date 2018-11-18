package engine

type (
	// EngineFactory interface allows us to provide
	// other parts of the system with a way to make
	// instances of our use-case / interactors when
	// they need to
	EngineFactory interface {
		// NewPost creates a new Post interactor
		NewPost() Post
	}

	// engine factory stores the state of our engine
	// which only involves a storage factory instance
	engineFactory struct {
		StorageFactory
		// JWTSignParser
	}
)

// NewEngine creates a new engine factory that will
// make use of the passed in StorageFactory for any
// data persistence needs.
func NewEngine(s StorageFactory /*, jwt JWTSignParser*/) EngineFactory {
	return &engineFactory{
		s,
		// jwt,
	}
}

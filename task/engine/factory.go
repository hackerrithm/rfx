package engine

type (
	// Factory interface allows us to provide
	// other parts of the system with a way to make
	// instances of our use-case / interactors when
	// they need to
	Factory interface {
		// NewTask creates a new Task interactor
		NewTask() Task
	}

	// engine factory stores the state of our engine
	// which only involves a storage factory instance
	engineFactory struct {
		StorageFactory
	}
)

// NewEngine creates a new engine factory that will
// make use of the passed in StorageFactory for any
// data persistence needs.
func NewEngine(s StorageFactory) Factory {
	return &engineFactory{s}
}

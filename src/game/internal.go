package game

var (
	state *State
	scene Scene

	entities []Entity
	textEntities []TextEntity

	loading chan bool
)

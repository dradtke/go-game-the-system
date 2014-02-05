package game

import (
	"github.com/dradtke/go-allegro/allegro"
	"reflect"
	"runtime"
)

type State struct {
	current     status
	prev        status
	display     *allegro.Display
	sceneLoaded bool
}

func (s *State) Display() *allegro.Display {
	return s.display
}

func (s *State) SceneLoaded() bool {
	return s.sceneLoaded
}

func (s *State) MouseLeftDown() bool {
	return s.current.MouseLeftDown
}

func (s *State) MouseRightDown() bool {
	return s.current.MouseRightDown
}

func (s *State) MouseOnScreen() bool {
	return s.current.MouseOnScreen
}

func (s *State) MouseX() int {
	return s.current.MouseX
}

func (s *State) MouseY() int {
	return s.current.MouseY
}

func (s *State) KeyDown(code allegro.KeyCode) bool {
	return s.current.KeyDown[code]
}

// information regarding the game that can change between
// ticks, and needs to be double-buffered
type status struct {
	MouseLeftDown, MouseRightDown, MouseOnScreen bool
	MouseX, MouseY                               int
	KeyDown                                      map[allegro.KeyCode]bool
}

func sync(from *status, to *status) {
	from_val := reflect.ValueOf(from).Elem()
	to_val := reflect.ValueOf(to).Elem()
	n := from_val.NumField()
	for i := 0; i < n; i++ {
		to_val.Field(i).Set(from_val.Field(i))
	}
}

func Init(display *allegro.Display) {
	state = new(State)
	state.current = *new(status)
	state.prev = *new(status)
	state.current.KeyDown = make(map[allegro.KeyCode]bool)
	state.prev.KeyDown = make(map[allegro.KeyCode]bool)
	state.display = display
}

func GoTo(eventQueue *allegro.EventQueue, sc Scene) {
	if scene != nil {
		scene.Leave()
	}
	scene = sc
	unregisterEventSources(eventQueue)
	entities = make([]Entity, 0)
	runtime.GC()
	scene.Enter()
	registerEventSources(eventQueue)
	state.sceneLoaded = false
	loading = make(chan bool, 1)
	go func() {
		scene.Load(state)
		loading <- true
	}()
}

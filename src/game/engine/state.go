package engine

import (
	"github.com/dradtke/go-allegro/allegro"
	"reflect"
)

type State struct {
	current status
	prev    status
	sceneLoaded bool
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

// information regarding the game that can change between
// ticks, and needs to be double-buffered
type status struct {
	MouseLeftDown, MouseRightDown, MouseOnScreen bool
	MouseX, MouseY                               int
}

var (
	state        *State
	scene        Scene

	loading chan bool
)

func sync(from *status, to *status) {
	from_val := reflect.ValueOf(from).Elem()
	to_val := reflect.ValueOf(to).Elem()
	n := from_val.NumField()
	for i := 0; i < n; i++ {
		to_val.Field(i).Set(from_val.Field(i))
	}
}

func Init() {
	state = new(State)
	state.current = *new(status)
	state.prev = *new(status)
}

func GoTo(eventQueue *allegro.EventQueue, sc Scene) {
	if scene != nil {
		scene.Leave()
	}
	scene = sc
	unregisterEventSources(eventQueue)
	scene.Enter()
	registerEventSources(eventQueue)
	state.sceneLoaded = false
	loading = make(chan bool, 1)
	go func() {
		scene.Load()
		loading <- true
	}()
}

func Update() {
	scene.Update(state)
	select {
	case <-loading:
		state.sceneLoaded = true
	default:
		// not yet loaded
	}

	if state.current.MouseLeftDown && !state.prev.MouseLeftDown {
		scene.OnLeftPress(state)
	} else if !state.current.MouseLeftDown && state.prev.MouseLeftDown {
		scene.OnLeftRelease(state)
	}

	if state.current.MouseRightDown && !state.prev.MouseRightDown {
		scene.OnRightPress(state)
	} else if !state.current.MouseRightDown && state.prev.MouseRightDown {
		scene.OnRightRelease(state)
	}

	sync(&state.current, &state.prev)
}

// TODO: verify that this is the correct value for delta
func Render(delta float64) {
	allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
	allegro.HoldBitmapDrawing(true)
	scene.Render(state, delta)
	allegro.HoldBitmapDrawing(false)
	allegro.FlipDisplay()
}

func HandleEvent(event *allegro.Event) (unhandled bool) {
	if !scene.HandleEvent(state, event) {
		return false
	}
	switch event.Type {
	case allegro.EVENT_MOUSE_BUTTON_DOWN:
		switch event.Mouse.Button {
		case 1:
			state.current.MouseLeftDown = true
		case 2:
			state.current.MouseRightDown = true
		}

	case allegro.EVENT_MOUSE_BUTTON_UP:
		switch event.Mouse.Button {
		case 1:
			state.current.MouseLeftDown = false
		case 2:
			state.current.MouseRightDown = false
		}

	case allegro.EVENT_MOUSE_ENTER_DISPLAY:
		state.current.MouseOnScreen = true

	case allegro.EVENT_MOUSE_LEAVE_DISPLAY:
		state.current.MouseOnScreen = false

	case allegro.EVENT_MOUSE_AXES:
		state.current.MouseOnScreen = true
		state.current.MouseX = event.Mouse.X
		state.current.MouseY = event.Mouse.Y
	default:
		unhandled = true
	}
	return
}

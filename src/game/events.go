package game

import (
	"container/list"
	"github.com/dradtke/go-allegro/allegro"
)

var customEvents list.List

func registerEventSources(queue *allegro.EventQueue) {
	for e := customEvents.Front(); e != nil; e = e.Next() {
		if source, ok := e.Value.(*allegro.EventSource); ok {
			queue.RegisterEventSource(source)
		}
	}
}

// Used to clear out a scene's custom events when the game
// is changing scenes.
func unregisterEventSources(queue *allegro.EventQueue) {
	for e := customEvents.Front(); e != nil; e = e.Next() {
		if source, ok := e.Value.(*allegro.EventSource); ok {
			queue.UnregisterEventSource(source)
		}
	}
	customEvents.Init()
}

// This should be called in a scene's Enter() method to
// add some event source to the global event queue.
func RegisterEventSource(source *allegro.EventSource) {
	customEvents.PushBack(source)
}

func HandleEvent(event *allegro.Event) {
	if !scene.HandleEvent(state, event) {
		return
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

	case allegro.EVENT_KEY_DOWN:
		state.current.KeyDown[event.Keyboard.KeyCode] = true

	case allegro.EVENT_KEY_UP:
		state.current.KeyDown[event.Keyboard.KeyCode] = false
	}
}

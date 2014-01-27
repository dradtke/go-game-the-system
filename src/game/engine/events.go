package engine

import (
	"container/list"
	"github.com/dradtke/go-allegro/allegro"
)

var customEvents list.List

// This should be called in a scene's Enter() method to
// add some event source to the global event queue.
func RegisterEventSource(source *allegro.EventSource) {
	customEvents.PushBack(source)
}

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

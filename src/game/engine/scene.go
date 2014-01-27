package engine

import (
	"github.com/dradtke/go-allegro/allegro"
)

type EventMap map[*allegro.EventSource]func(*allegro.Event)

type Scene interface {
	Enter()
	Leave()
	Load()
	HandleEvent(*State, *allegro.Event) bool
	Update(*State)
	Render(*State, float64)
	OnLeftPress(*State)
	OnRightPress(*State)
	OnLeftRelease(*State)
	OnRightRelease(*State)
}

type BaseScene struct{}

func (s *BaseScene) Enter() {}
func (s *BaseScene) Leave() {}
func (s *BaseScene) Load()  {}
func (s *BaseScene) HandleEvent(state *State, event *allegro.Event) bool {
	return true
}
func (s *BaseScene) Update(state *State) {}
func (s *BaseScene) Render(state *State, delta float64) {}

// events
func (s *BaseScene) OnLeftPress(state *State)    {}
func (s *BaseScene) OnRightPress(state *State)   {}
func (s *BaseScene) OnLeftRelease(state *State)  {}
func (s *BaseScene) OnRightRelease(state *State) {}

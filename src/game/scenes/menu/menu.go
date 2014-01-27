package menu

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"game/engine"
	"game/engine/widget"
	"game/shared"
	"time"
)

type MenuScene struct {
	engine.BaseScene
	Name string // which menu are we in?

	dotTimer    *allegro.Timer
	loadingDots string
	images      map[string]*allegro.Bitmap

	button  *widget.Button
	widgets []widget.Widget
}

// Enter() is used for scene initialization. It should return a slice of
// EventSources that should be added to the global event queue. They
// will be automatically removed when the game changes scenes.
func (s *MenuScene) Enter() {
	fmt.Println("entering menu scene.")

	var err error
	if s.dotTimer, err = allegro.CreateTimer(0.5); err == nil {
		engine.RegisterEventSource(s.dotTimer.EventSource())
		s.dotTimer.Start()
	}

	s.button = &widget.Button{
		X: 200, Y: 200,
		OnClick: func() {
			fmt.Println("click!")
		},
		OnPress: func() {
			fmt.Println("pressing...")
		},
	}
	s.widgets = []widget.Widget{s.button}
}

// Load() should be used for images and other resources that may take a while
// to load into memory. This method is always run in its own goroutine, and when
// it finishes, the state's SceneLoaded() property will be set to true.
func (s *MenuScene) Load() {
	s.images = shared.LoadImages([]string{
		"src/game/scenes/menu/img/button.png",
		"src/game/scenes/menu/img/button-hover.png",
	})
	s.button.Base = s.images["button.png"]
	s.button.Hover = s.images["button-hover.png"]
	time.Sleep(3 * time.Second)
}

func (s *MenuScene) Render(state *engine.State, delta float64) {
	if !state.SceneLoaded() {
		font.DrawText(shared.BuiltinFont(), allegro.MapRGB(0xFF, 0xFF, 0xFF),
			10, 10, font.ALIGN_LEFT, "Loading"+s.loadingDots)
		return
	}
	s.button.Draw(state)
}

func (s *MenuScene) Leave() {
	fmt.Println("leaving menu scene.")
}

func (s *MenuScene) OnLeftPress(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
	for _, w := range s.widgets {
		w.Press(state)
	}
}

func (s *MenuScene) OnLeftRelease(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
	for _, w := range s.widgets {
		w.Release(state)
	}
}

func (s *MenuScene) HandleEvent(state *engine.State, event *allegro.Event) bool {
	switch event.Source {
	case s.dotTimer.EventSource():
		if s.loadingDots == "..." {
			s.loadingDots = ""
		} else {
			s.loadingDots += "."
		}
	default:
		return true
	}
	return false
}

package menu

import (
	"fmt"
	"game/engine"
	"game/engine/widget"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"time"
)

type MenuScene struct {
	engine.BaseScene
	Name string // which menu are we in?

	dotTimer    *allegro.Timer
	loadingDots string
}

// Enter() is used for scene initialization. After this method exits, any
// calls to engine.RegisterEventSource() that occurred within it will then
// have their respective event sources added to the global queue. They will
// be automatically removed when the game changes scenes.
func (s *MenuScene) Enter() {
	var err error
	if s.dotTimer, err = allegro.CreateTimer(0.5); err == nil {
		engine.RegisterEventSource(s.dotTimer.EventSource())
		s.dotTimer.Start()
	}
}

// Load() should be used for images and other resources that may take a while
// to load into memory. This method is always run in its own goroutine, and when
// it finishes, the state's SceneLoaded() property will be set to true.
func (s *MenuScene) Load(state *engine.State) {
	images := engine.LoadImages([]string{
		"src/game/scenes/menu/img/orb/norm.png",
		"src/game/scenes/menu/img/orb/lit.png",
		"src/game/scenes/menu/img/orb/down.png",
	})
	engine.AddEntity(&widget.Button{
		Base: images["norm.png"],
		Hover: images["lit.png"],
		Pressed: images["down.png"],
		X: 300,
		Y: 200,
		Radius: 45,
		Bound: widget.Circle,
		OnHover: func() {
			fmt.Println("hover")
		},
		OnPress: func() {
			fmt.Println("press")
		},
		OnClick: func() {
			fmt.Println("click")
		},
	})
	time.Sleep(3 * time.Second) // fake an additional 2-second load time
	s.dotTimer.Stop()
}

func (s *MenuScene) Update(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
}

// Render()
func (s *MenuScene) Render(state *engine.State, delta float32) {
	if !state.SceneLoaded() {
		font.DrawText(engine.BuiltinFont(), allegro.MapRGB(0xFF, 0xFF, 0xFF),
			10, 10, font.ALIGN_LEFT, "Loading"+s.loadingDots)
		return
	}
}

func (s *MenuScene) Leave() {
	fmt.Println("leaving menu scene.")
}

func (s *MenuScene) OnLeftPress(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
}

func (s *MenuScene) OnLeftRelease(state *engine.State) {
	if !state.SceneLoaded() {
		return
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

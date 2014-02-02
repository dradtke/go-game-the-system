package menu

import (
	"fmt"
	"game/engine"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	prim "github.com/dradtke/go-allegro/allegro/primitives"
	"time"
)

type MenuScene struct {
	engine.BaseScene
	Name string // which menu are we in?

	dotTimer    *allegro.Timer
	loadingDots string
	square      *square
}

type square struct {
	engine.Entity
	img *allegro.Bitmap
	dir int8
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
	s.square = &square{
		engine.Entity{X: 100, Y: 100},
		allegro.CreateBitmap(30, 30), // img
		1, // dir
	}
	s.square.img.AsTarget(func() {
		prim.DrawFilledRectangle(prim.Point{X: 0, Y: 0}, prim.Point{X: 30, Y: 30}, allegro.MapRGB(0, 0xFF, 0))
	})
	time.Sleep(2 * time.Second) // fake an additional 2-second load time
}

func (s *MenuScene) Update(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
	if state.KeyDown(allegro.KEY_UP) {
		s.square.Move(0, -5)
	} else if state.KeyDown(allegro.KEY_DOWN) {
		s.square.Move(0, 5)
	} else if state.KeyDown(allegro.KEY_RIGHT) {
		s.square.Move(5, 0)
	} else if state.KeyDown(allegro.KEY_LEFT) {
		s.square.Move(-5, 0)
	}
}

// Render()
func (s *MenuScene) Render(state *engine.State, delta float32) {
	if !state.SceneLoaded() {
		font.DrawText(engine.BuiltinFont(), allegro.MapRGB(0xFF, 0xFF, 0xFF),
			10, 10, font.ALIGN_LEFT, "Loading"+s.loadingDots)
		return
	}
	square_x := s.square.X + (float32(5*int(s.square.dir)) * delta)
	s.square.img.Draw(square_x, s.square.Y, allegro.FLIP_NONE)
}

func (s *MenuScene) Leave() {
	fmt.Println("leaving menu scene.")
}

func (s *MenuScene) OnLeftPress(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
	/*
		for _, w := range s.widgets {
			w.Press(state)
		}
	*/
}

func (s *MenuScene) OnLeftRelease(state *engine.State) {
	if !state.SceneLoaded() {
		return
	}
	/*
		for _, w := range s.widgets {
			w.Release(state)
		}
	*/
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

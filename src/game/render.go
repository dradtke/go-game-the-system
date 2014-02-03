package game

import (
	"github.com/dradtke/go-allegro/allegro"
)

func Render(delta float32) {
	allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
	allegro.HoldBitmapDrawing(true)
	if state.sceneLoaded {
		for _, e := range entities {
			e.Render(state, delta)
		}
	}
	scene.Render(state, delta)
	allegro.HoldBitmapDrawing(false)
	allegro.FlipDisplay()
}

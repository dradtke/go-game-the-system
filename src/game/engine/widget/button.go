package widget

import (
	"github.com/dradtke/go-allegro/allegro"
	"math"
	"game/engine"
)

type Button struct {
	Base, Hover, Pressed *allegro.Bitmap
	X, Y                 float32
	Radius               int // optional radius override for circular buttons
	Bound                Bound
	OnClick, OnPress     func()

	width, height            int
	left, right, top, bottom int
	is_pressed               bool
}

type Bound int

const (
	Rectangle Bound = iota // default
	Circle
	Pixel
)

func (b *Button) IsHovering(state *engine.State) bool {
	if !state.MouseOnScreen() {
		return false
	}
	var (
		mouse_x = state.MouseX()
		mouse_y = state.MouseY()
	)
	switch b.Bound {
	case Rectangle:
		if mouse_x >= b.left && mouse_x <= b.right && mouse_y >= b.top && mouse_y <= b.bottom {
			return true
		}
	case Circle:
		x := float32(mouse_x) - b.X
		y := float32(mouse_y) - b.Y
		dist := math.Sqrt(float64(x*x) + float64(y*y))
		// TODO: add a separate radius field?
		if dist <= float64(b.width) {
			return true
		}
	case Pixel:
		if mouse_x >= b.left && mouse_x <= b.right && mouse_y >= b.top && mouse_y <= b.bottom {
			x := float32(mouse_x) - b.X
			y := float32(mouse_y) - b.Y
			_, _, _, alpha := b.Base.Pixel(int(x), int(y)).UnmapRGBA()
			return alpha > 0
		}
	}
	return false
}

func (b *Button) Draw(state *engine.State) {
	if b.width == 0 {
		b.width = b.Base.Width()
	}
	if b.height == 0 {
		b.height = b.Base.Height()
	}
	if b.Bound == Circle && b.Radius == 0 {
		// If the button is circular, default its radius to the width.
		b.Radius = b.width
	}
	b.left = int(b.X)
	b.top = int(b.Y)
	b.right = b.left + b.width
	b.bottom = b.top + b.height

	if b.IsHovering(state) && b.Hover != nil {
		if b.is_pressed && b.Pressed != nil {
			b.Pressed.Draw(b.X, b.Y, allegro.FLIP_NONE)
		} else {
			b.Hover.Draw(b.X, b.Y, allegro.FLIP_NONE)
		}
	} else {
		b.Base.Draw(b.X, b.Y, allegro.FLIP_NONE)
	}
}

func (b *Button) Press(state *engine.State) {
	b.is_pressed = b.IsHovering(state)
	if b.is_pressed && b.OnPress != nil {
		b.OnPress()
	}
}

func (b *Button) Release(state *engine.State) {
	if b.is_pressed && b.IsHovering(state) && b.OnClick != nil {
		b.OnClick()
	}
	b.is_pressed = false
}

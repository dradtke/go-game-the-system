package widget

import (
	"game/engine"
	"github.com/dradtke/go-allegro/allegro"
	"math"
)

type Button struct {
	Base, Hover, Pressed      *allegro.Bitmap
	X, Y                      float32
	Radius                    int // optional radius override for circular buttons
	Bound                     Bound
	OnClick, OnPress, OnHover func()

	width, height            int
	left, right, top, bottom int
	is_hovering, is_pressed  bool
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
		x := float32(mouse_x) - float32(b.X + float32(b.width/2))
		y := float32(mouse_y) - float32(b.Y + float32(b.height/2))
		dist := math.Sqrt(float64(x*x) + float64(y*y))
		if (b.Radius != 0 && dist <= float64(b.Radius)) || (b.Radius == 0 && dist <= float64(b.width)) {
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

func (b *Button) Update(state *engine.State) {
	if !b.is_hovering && !b.is_pressed {
		if b.IsHovering(state) && !state.MouseLeftDown() {
			b.is_hovering = true
			if b.OnHover != nil {
				b.OnHover()
			}
		}
	} else if b.is_hovering && !b.is_pressed {
		if !b.IsHovering(state) {
			b.is_hovering = false
		} else if state.MouseLeftDown() {
			b.is_pressed = true
			if b.OnPress != nil {
				b.OnPress()
			}
		}
	} else if b.is_pressed {
		if !b.IsHovering(state) {
			b.is_hovering = false
			b.is_pressed = false
		} else if b.IsHovering(state) && !state.MouseLeftDown() {
			b.is_pressed = false
			if b.OnClick != nil {
				b.OnClick()
			}
		}
	}
}

func (b *Button) Render(state *engine.State, delta float32) {
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

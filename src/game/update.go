package game

import (
	"github.com/dradtke/go-allegro/allegro"
	"strconv"
)

func Update() {
	defer func() {
		sync(&state.current, &state.prev)
	}()

	if state.sceneLoaded {
		for _, e := range entities {
			e.Update(state)
		}
	}
	scene.Update(state)

	select {
	case <-loading:
		state.sceneLoaded = true
		return
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
}

func updateText(keyCode allegro.KeyCode, c rune) {
	if !state.sceneLoaded {
		return
	}
	// TODO: limit this to only the focused widget
	for _, t := range textEntities {
		switch keyCode {
		case allegro.KEY_BACKSPACE:
			t.Backspace()
		case allegro.KEY_LEFT:
			t.ShiftCursor(-1)
		case allegro.KEY_RIGHT:
			t.ShiftCursor(1)
		default:
			if strconv.IsPrint(c) {
				t.Write(c)
			}
		}
	}
}

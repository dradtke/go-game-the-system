package widget

import (
	"game"
)

type Widget interface {
	Press(state *game.State)
	Release(state *game.State)
	Draw(state *game.State)
}

package widget

import (
	"game/engine"
)

type Widget interface {
	Press(state *engine.State)
	Release(state *engine.State)
	Draw(state *engine.State)
}

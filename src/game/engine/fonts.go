package engine

import (
	"github.com/dradtke/go-allegro/allegro/font"
)

var builtinFont *font.Font

func BuiltinFont() *font.Font {
	if builtinFont == nil {
		var err error
		builtinFont, err = font.Builtin()
		if err != nil {
			panic(err)
		} else if builtinFont == nil {
			panic("Allegro returned nil built-in font!")
		}
	}
	return builtinFont
}

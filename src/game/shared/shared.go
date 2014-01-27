package shared

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"os"
	"path/filepath"
)

var builtinFont *font.Font

func init() {
	if f, err := font.Builtin(); err == nil {
		builtinFont = f
	} else {
		panic(err)
	}
}

func BuiltinFont() *font.Font {
	return builtinFont
}

func LoadImages(paths []string) map[string]*allegro.Bitmap {
	images := make(map[string]*allegro.Bitmap)
	for _, path := range paths {
		if bmp, err := allegro.LoadBitmap(path); err == nil {
			images[filepath.Base(path)] = bmp
		} else {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}
	return images
}

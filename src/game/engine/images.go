package engine

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"os"
	"path/filepath"
)

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

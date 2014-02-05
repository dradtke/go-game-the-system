package widget

import (
	"bytes"
	"container/list"
	"game"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
)

type Input struct {
	X, Y      float32
	MaxLength int
	Font      *font.Font
	Align     font.DrawFlags
	Color     allegro.Color
	// TODO: add caret blink option

	text       list.List
	val        string
	is_focused bool
	caret      *allegro.Bitmap
	caretRef   *list.Element
	caretPos   int
}

func (i *Input) save() {
	var buf bytes.Buffer
	for e := i.text.Front(); e != nil; e = e.Next() {
		buf.WriteRune(e.Value.(rune))
	}
	i.val = buf.String()
}

func (i *Input) Write(c rune) {
	if i.MaxLength > 0 && i.text.Len() == i.MaxLength {
		return
	}
	if i.caretRef != nil {
		i.caretRef = i.text.InsertAfter(c, i.caretRef)
	} else {
		i.caretRef = i.text.PushFront(c)
	}
	i.caretPos++
	i.save()
}

func (i *Input) Backspace() {
	if i.caretRef == nil {
		return
	}
	oldRef := i.caretRef
	i.caretRef = i.caretRef.Prev()
	i.text.Remove(oldRef)
	i.caretPos--
	i.save()
}

func (i *Input) ShiftCursor(shift int) {
	for shift < 0 && i.caretPos > 0 {
		shift++
		i.caretPos--
		i.caretRef = i.caretRef.Prev()
	}
	for shift > 0 && i.caretPos < i.text.Len() {
		shift--
		i.caretPos++
		if i.caretRef != nil {
			i.caretRef = i.caretRef.Next()
		} else {
			i.caretRef = i.text.Front()
		}
	}
}

func (i *Input) Update(state *game.State) {
	if i.Font == nil {
		i.Font = game.BuiltinFont()
	}
	if i.caret == nil {
		i.caret = allegro.CreateBitmap(2, i.Font.LineHeight())
		i.caret.AsTarget(func() {
			allegro.ClearToColor(i.Color)
		})
	}
}

func (i *Input) Render(state *game.State, delta float32) {
	font.DrawText(i.Font, i.Color, i.X, i.Y, i.Align, i.val)
	if i.caret != nil {
		i.caret.Draw(i.X+float32(i.Font.TextWidth(i.val[:i.caretPos])), i.Y, allegro.FLIP_NONE)
	}
}

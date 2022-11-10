package bento

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Editable struct {
	cursor        int
	cursorTime    int64
	displayCursor bool
	focus         bool
}

func (e *Editable) Cursor() int {
	if e.displayCursor {
		return e.cursor
	}
	return -1
}

func (e *Editable) Update(b *Box) error {
	if e == nil {
		return nil
	}
	if b.State == Disabled {
		return nil
	}
	if b.State == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		e.focus = true
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		e.focus = false
	}
	if e.focus {
		t := time.Now().UnixMilli()
		e.displayCursor = false
		if t-e.cursorTime <= 1000 {
			e.displayCursor = true
		} else if t-e.cursorTime >= 2000 {
			e.cursorTime = t
		}
		v := b.Content
		for _, k := range keys {
			if repeatingKeyPressed(k) {
				s := keyToString(k, ebiten.IsKeyPressed(ebiten.KeyShift))
				if s != "" && s != "\n" {
					v = v[:e.cursor] + s + v[e.cursor:]
					e.cursor++
					e.cursorTime = t
				} else if e.cursor > 0 && len(v) > 0 && k == ebiten.KeyBackspace {
					v = v[:e.cursor-1] + v[e.cursor:]
					e.cursor--
					e.cursorTime = t
				} else if e.cursor > 0 && k == ebiten.KeyLeft {
					e.cursor--
					e.cursorTime = t
				} else if e.cursor < len(v) && k == ebiten.KeyRight {
					e.cursor++
					e.cursorTime = t
				}
				// TODO: ebiten.KeyUp, ebiten.KeyDown
			}
		}
		if b.Content != v {
			b.Content = v
			b.fireEvent(Change)
		}
	} else {
		e.displayCursor = false
	}
	return nil
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func keyToString(k ebiten.Key, shift bool) string {
	if shift {
		switch k {
		case ebiten.Key0:
			return ")"
		case ebiten.Key1:
			return "!"
		case ebiten.Key2:
			return "@"
		case ebiten.Key3:
			return "#"
		case ebiten.Key4:
			return "$"
		case ebiten.Key5:
			return "%"
		case ebiten.Key6:
			return "^"
		case ebiten.Key7:
			return "&"
		case ebiten.Key8:
			return "*"
		case ebiten.Key9:
			return "("
		case ebiten.KeyA:
			return "A"
		case ebiten.KeyB:
			return "B"
		case ebiten.KeyC:
			return "C"
		case ebiten.KeyD:
			return "D"
		case ebiten.KeyE:
			return "E"
		case ebiten.KeyF:
			return "F"
		case ebiten.KeyG:
			return "G"
		case ebiten.KeyH:
			return "H"
		case ebiten.KeyI:
			return "I"
		case ebiten.KeyJ:
			return "J"
		case ebiten.KeyK:
			return "K"
		case ebiten.KeyL:
			return "L"
		case ebiten.KeyM:
			return "M"
		case ebiten.KeyN:
			return "N"
		case ebiten.KeyO:
			return "O"
		case ebiten.KeyP:
			return "P"
		case ebiten.KeyQ:
			return "Q"
		case ebiten.KeyR:
			return "R"
		case ebiten.KeyS:
			return "S"
		case ebiten.KeyT:
			return "T"
		case ebiten.KeyU:
			return "U"
		case ebiten.KeyV:
			return "V"
		case ebiten.KeyW:
			return "W"
		case ebiten.KeyX:
			return "X"
		case ebiten.KeyY:
			return "Y"
		case ebiten.KeyZ:
			return "Z"
		case ebiten.KeyComma:
			return "<"
		case ebiten.KeyPeriod:
			return ">"
		case ebiten.KeySemicolon:
			return ":"
		case ebiten.KeyQuote:
			return "\""
		case ebiten.KeyBracketLeft:
			return "{"
		case ebiten.KeyBracketRight:
			return "}"
		case ebiten.KeyMinus:
			return "_"
		case ebiten.KeyEqual:
			return "+"
		case ebiten.KeyBackslash:
			return "|"
		case ebiten.KeySlash:
			return "?"
		case ebiten.KeyBackquote:
			return "~"
		case ebiten.KeySpace:
			return " "
		case ebiten.KeyEnter:
			return "\n"
		}
	}
	switch k {
	case ebiten.Key0:
		return "0"
	case ebiten.Key1:
		return "1"
	case ebiten.Key2:
		return "2"
	case ebiten.Key3:
		return "3"
	case ebiten.Key4:
		return "4"
	case ebiten.Key5:
		return "5"
	case ebiten.Key6:
		return "6"
	case ebiten.Key7:
		return "7"
	case ebiten.Key8:
		return "8"
	case ebiten.Key9:
		return "9"
	case ebiten.KeyA:
		return "a"
	case ebiten.KeyB:
		return "b"
	case ebiten.KeyC:
		return "c"
	case ebiten.KeyD:
		return "d"
	case ebiten.KeyE:
		return "e"
	case ebiten.KeyF:
		return "f"
	case ebiten.KeyG:
		return "g"
	case ebiten.KeyH:
		return "h"
	case ebiten.KeyI:
		return "i"
	case ebiten.KeyJ:
		return "j"
	case ebiten.KeyK:
		return "k"
	case ebiten.KeyL:
		return "l"
	case ebiten.KeyM:
		return "m"
	case ebiten.KeyN:
		return "n"
	case ebiten.KeyO:
		return "o"
	case ebiten.KeyP:
		return "p"
	case ebiten.KeyQ:
		return "q"
	case ebiten.KeyR:
		return "r"
	case ebiten.KeyS:
		return "s"
	case ebiten.KeyT:
		return "t"
	case ebiten.KeyU:
		return "u"
	case ebiten.KeyV:
		return "v"
	case ebiten.KeyW:
		return "w"
	case ebiten.KeyX:
		return "x"
	case ebiten.KeyY:
		return "y"
	case ebiten.KeyZ:
		return "z"
	case ebiten.KeyComma:
		return ","
	case ebiten.KeyPeriod:
		return "."
	case ebiten.KeySemicolon:
		return ";"
	case ebiten.KeyQuote:
		return "'"
	case ebiten.KeyBracketLeft:
		return "["
	case ebiten.KeyBracketRight:
		return "]"
	case ebiten.KeyMinus:
		return "-"
	case ebiten.KeyEqual:
		return "="
	case ebiten.KeyBackslash:
		return "\\"
	case ebiten.KeySlash:
		return "/"
	case ebiten.KeyBackquote:
		return "`"
	case ebiten.KeySpace:
		return " "
	case ebiten.KeyEnter:
		return "\n"
	}
	return ""
}

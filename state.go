package bento

import (
	"image"
	"log"
	"reflect"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	Scrollspeed = 0.1
)

type State int

const (
	Idle     = State(0)
	Hover    = State(1)
	Active   = State(2)
	Disabled = State(3)
)

type Event struct {
	X, Y int
	Box  *Box
}

type state struct {
	state          State
	scrollState    [4]State
	cursor         int
	scrollLine     int
	scrollPosition float64
	cursorTime     int64
}

func (n *Box) updateState(keys []ebiten.Key) {
	if !n.style.display() || n.style.hidden() {
		return
	}
	if n.tag == "input" || n.tag == "textarea" {
		n.updateInput()
	} else {
		n.state.state = getState(n.innerRect())
	}
	if n.style.Scrollbar != nil {
		n.updateScroll()
	}
	if n.state.state == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && n.attrs["onClick"] != "" {
		x, y := ebiten.CursorPosition()
		attr := n.attrs["onClick"]
		m := reflect.ValueOf(n.component).MethodByName(attr)
		if !m.IsValid() {
			log.Fatalf("%s can't find onClick handler named %q in component %s", n.tag, attr, reflect.TypeOf(n.component))
		}
		var args []reflect.Value
		if m.Type().NumIn() == 1 {
			args = []reflect.Value{reflect.ValueOf(&Event{
				X:   x - n.X,
				Y:   y - n.Y,
				Box: n,
			})}
		}
		m.Call(args)
	}
	for _, c := range n.children {
		c.updateState(keys)
	}
}

func (n *Box) updateInput() {
	state := getState(n.innerRect())
	if state == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		n.state.state = Active
	} else if n.state.state != Active || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		n.state.state = state
	}
	if n.state.state == Active {
		v := n.attrs["value"]
		for _, k := range keys {
			if repeatingKeyPressed(k) {
				s := keyToString(k, ebiten.IsKeyPressed(ebiten.KeyShift))
				if s != "" && s != "\n" {
					v = v[:n.cursor] + s + v[n.cursor:]
					n.cursor++
					n.cursorTime = time.Now().UnixMilli()
				} else if n.cursor > 0 && len(v) > 0 && k == ebiten.KeyBackspace {
					v = v[:n.cursor-1] + v[n.cursor:]
					n.cursor--
					n.cursorTime = time.Now().UnixMilli()
				} else if n.cursor > 0 && k == ebiten.KeyLeft {
					n.cursor--
					n.cursorTime = time.Now().UnixMilli()
				} else if n.cursor < len(v) && k == ebiten.KeyRight {
					n.cursor++
					n.cursorTime = time.Now().UnixMilli()
				}
				// TODO: ebiten.KeyUp, ebiten.KeyDown
			}
		}
		n.attrs["value"] = v
	}
}

func (n *Box) updateScroll() {
	mt, _, _, ml := n.style.margin()
	pt, _, _, pl := n.style.padding()
	rects := n.scrollRects(n.scrollPosition)
	for i := 0; i < 4; i++ {
		if i == 2 {
			continue
		}
		// TODO: the math here works out but it's confusing
		r := rects[i].Add(image.Pt(n.X+ml+pl+pl, n.Y+mt+pt-pt))
		n.scrollState[i] = getState(r)
		if n.scrollState[i] == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if i == 0 {
				n.scrollLine--
				if n.scrollLine < 0 {
					n.scrollLine = 0
				}
			} else if i == 3 && n.scrollPosition < 1 {
				n.scrollLine++
			}
		}
	}
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

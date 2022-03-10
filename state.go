package bento

import (
	"reflect"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	Scrollspeed = 0.1
)

type state struct {
	buttonState    ButtonState
	inputState     ButtonState
	cursor         int
	scrollPosition int
	cursorTime     int64
}

func (n *Box) updateState(keys []ebiten.Key) {
	switch n.tag {
	case "button":
		n.updateButton()
	case "input", "textarea":
		n.updateInput()
	}
}

func (n *Box) updateButton() {
	n.buttonState = buttonState(n.innerRect())
	if n.buttonState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		attr := n.attrs["onClick"]
		m := reflect.ValueOf(n.component).MethodByName(attr)
		if m.IsValid() {
			m.Call(nil)
		}
	}
}

func (n *Box) updateInput() {
	state := buttonState(n.innerRect())
	if state == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		n.inputState = Active
	} else if n.inputState != Active || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		n.inputState = state
	}
	if n.inputState == Active {
		v := n.attrs["value"]
		for _, k := range keys {
			if repeatingKeyPressed(k) {
				s := keyToString(k, ebiten.IsKeyPressed(ebiten.KeyShift))
				if s != "" {
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
				} else if k == ebiten.KeyUp {
					// TODO
				} else if k == ebiten.KeyDown {
					// TODO
				}
			}
		}
		n.attrs["value"] = v
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

/*
type Scrollarea struct {
	states                                               [3][4]*NineSlice
	position                                             float64
	topBtnState, bottomBtnState, trackState, handleState ButtonState
	buffer                                               *ebiten.Image
}

func (s *Scrollarea) Update(keys []ebiten.Key, box *Box) ([]ebiten.Key, error) {
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	h := s.states[0][0].heights[0] + s.states[0][0].heights[1] + s.states[0][0].heights[2]

	r := box.innerRect()
	s.topBtnState = buttonState(r)
	if s.topBtnState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.position -= Scrollspeed
	}
	trackHeight := r.Dy() - 2*h
	s.trackState = buttonState(image.Rect(r.Min.X, r.Min.Y+h, r.Min.X+w, r.Min.Y+h+trackHeight))
	s.handleState = buttonState(image.Rect(r.Min.X, r.Min.Y+h+int(s.position*float64(trackHeight-h)), r.Min.X+w, r.Min.Y+h+int(s.position*float64(trackHeight-h))+h))
	s.bottomBtnState = buttonState(image.Rect(r.Min.X, r.Max.Y-h, r.Min.X+w, r.Max.X))
	if s.bottomBtnState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.position += Scrollspeed
	}

	_, dy := ebiten.Wheel()
	s.position += dy * Scrollspeed

	// TODO: bad position can crash during Draw
	if s.position < 0 {
		s.position = 0
	}
	if s.position > 1 {
		s.position = 1
	}
	return keys, nil
}

func (s *Scrollarea) Draw(screen *ebiten.Image, n *Box) {

	n.buffer = ebiten.NewImage(n.TextBounds.Dx(), n.TextBounds.Dy())
	text.DrawParagraph(n.buffer, n.templateContent(), n.style.Font, n.style.Color, 0, 0, n.style.MaxWidth, -n.TextBounds.Min.Y)
	offset := int(scroll.position * float64(n.TextBounds.Dy()))
	cropped := ebiten.NewImageFromImage(n.buffer.SubImage(image.Rect(0, offset, content.Dx(), content.Dy()+offset)))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
	img.DrawImage(cropped, op)
	n.style.Scrollbar.Draw(img, n)

	r := n.innerRect()
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	h := s.states[0][0].heights[0] + s.states[0][0].heights[1] + s.states[0][0].heights[2]
	trackHeight := r.Dy() - 2*h
	s.states[s.topBtnState][0].Draw(screen, r.Min.X, r.Min.Y, w, h)
	s.states[s.trackState][1].Draw(screen, r.Min.X, r.Min.Y+h, w, trackHeight)
	s.states[s.handleState][2].Draw(screen, r.Min.X, r.Min.Y+h+int(s.position*float64(trackHeight-h)), w, h)
	s.states[s.bottomBtnState][3].Draw(screen, r.Min.X, r.Max.Y-h, w, h)
}
*/

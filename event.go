package bento

import (
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type EventType string

const (
	Click  = EventType("Click")
	Scroll = EventType("Scroll")
	Hover  = EventType("Hover")
	Change = EventType("Change")
	Draw   = EventType("Draw")
	Update = EventType("Update")
)

type Event struct {
	X, Y             int // Mouse position, relative to the current box
	ScrollX, ScrollY float64
	Box              *Box
	Type             EventType
	Image            *ebiten.Image            // Canvas element
	Op               *ebiten.DrawImageOptions // Canvas element
	Value            string                   // Input and Textarea elements
}

func (n *Box) fireEvent(e EventType, value string, img *ebiten.Image, op *ebiten.DrawImageOptions) bool {
	x, y := ebiten.CursorPosition()
	sx, sy := ebiten.Wheel()
	return n.call("on"+string(e), &Event{
		X:       x - n.X,
		Y:       y - n.Y,
		ScrollX: sx,
		ScrollY: sy,
		Type:    e,
		Box:     n,
		Value:   value,
		Image:   img,
		Op:      op,
	})
}

func (n *Box) call(attr string, args ...interface{}) bool {
	fnName := n.Attrs[attr]
	if fnName == "" {
		return false
	}
	m := reflect.ValueOf(n.Component).MethodByName(fnName)
	if !m.IsValid() {
		log.Fatalf("%s can't find %s handler named %q", reflect.TypeOf(n.Component), attr, fnName)
	}
	t := m.Type()
	var reflectArgs []reflect.Value
	for i, arg := range args {
		if i < t.NumIn() {
			reflectArgs = append(reflectArgs, reflect.ValueOf(arg))
		}
	}
	out := m.Call(reflectArgs)
	root := n.root()
	if len(out) > 0 {
		v := out[0]
		if out[0].Type().Kind() == reflect.Bool {
			if v.Bool() {
				root.dirty = true
			}
			return true
		}
	}
	root.dirty = true
	return true
}

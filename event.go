package bento

import (
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type EventType string

const (
	Click  = EventType("Click")
	Change = EventType("Change")
	Update = EventType("Update")
	Draw   = EventType("Draw")
)

type Event struct {
	X, Y int
	Box  *Box
	Type EventType
}

func (n *Box) fireEvent(e EventType) {
	attr := n.attrs["on"+string(e)]
	if attr == "" {
		return
	}
	x, y := ebiten.CursorPosition()
	m := reflect.ValueOf(n.Component).MethodByName(attr)
	if !m.IsValid() {
		log.Fatalf("%s can't find on%s handler named %q in component %s", n.Tag, e, attr, reflect.TypeOf(n.Component))
	}
	var args []reflect.Value
	if m.Type().NumIn() == 1 {
		args = []reflect.Value{reflect.ValueOf(&Event{
			X:    x - n.X,
			Y:    y - n.Y,
			Box:  n,
			Type: e,
		})}
	}
	m.Call(args)
}

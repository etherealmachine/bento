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
)

type Event struct {
	X, Y int
	Box  *Box
	Type EventType
}

func (n *Box) fireEvent(e EventType) {
	x, y := ebiten.CursorPosition()
	n.call("on"+string(e), &Event{
		X:    x - n.X,
		Y:    y - n.Y,
		Box:  n,
		Type: e,
	})
}

func (n *Box) call(attr string, args ...interface{}) {
	fnName := n.attrs[attr]
	if fnName == "" {
		return
	}
	m := reflect.ValueOf(n.Component).MethodByName(fnName)
	if !m.IsValid() {
		log.Fatalf("%s can't find %s handler named %q", reflect.TypeOf(n.Component), attr, fnName)
	}
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}
	m.Call(reflectArgs)
}
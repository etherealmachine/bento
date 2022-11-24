package bento

import (
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type EventType string

const (
	Click  = EventType("Click")
	Hover  = EventType("Hover")
	Change = EventType("Change")
	Update = EventType("Update")
)

type Event struct {
	X, Y  int
	Box   *Box
	Type  EventType
	Value string
}

func (n *Box) fireEvent(e EventType, value string) bool {
	x, y := ebiten.CursorPosition()
	return n.call("on"+string(e), &Event{
		X:     x - n.X,
		Y:     y - n.Y,
		Box:   n,
		Value: value,
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

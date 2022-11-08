/*
Bento is an XML based UI builder for Ebiten
*/
package bento

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"reflect"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
)

func (n *Box) isSubcomponent() bool {
	r, _ := utf8.DecodeRuneInString(n.Tag)
	return unicode.IsUpper(r)
}

func (n *Box) buildSubcomponent() error {
	m := reflect.ValueOf(n.Component).MethodByName(n.Tag)
	if !m.IsValid() {
		return fmt.Errorf("%s has no method named %s", reflect.TypeOf(n.Component), n.Tag)
	}
	res := m.Call(nil)
	if style, ok := res[0].Interface().(*Style); ok {
		n.Style = style
		n.Tag = style.Extends
		return nil
	} else if sub, ok := res[0].Interface().(Component); ok {
		subNode, err := Build(sub)
		if err != nil {
			return err
		}
		subNode.Parent = n.Parent
		*n = *subNode
		return nil
	}
	return fmt.Errorf("%s.%s must return either Style or Component", reflect.TypeOf(n.Component), n.Tag)
}

func (n *Box) build(prev *Box) error {
	if n.Tag == "" {
		tmpl, err := template.New("").Parse(n.Component.UI())
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, n.Component); err != nil {
			return err
		}
		if err := xml.Unmarshal(buf.Bytes(), n); err != nil {
			return fmt.Errorf("error building %s: %s", reflect.ValueOf(n.Component).Elem().Type().Name(), err)
		}
	}
	if n.isSubcomponent() {
		if err := n.buildSubcomponent(); err != nil {
			return err
		}
	}
	if prev != nil && n.Tag == prev.Tag {
		n.Debug = prev.Debug
		n.state = prev.state
		if n.Tag == "input" || n.Tag == "textarea" {
			n.Attrs["value"] = prev.Attrs["value"]
			n.Attrs["placeholder"] = prev.Attrs["placeholder"]
		}
	}
	if n.Style == nil {
		n.Style = new(Style)
	}
	n.Style.adopt(n)
	if err := n.Style.parseAttributes(); err != nil {
		return err
	}
	if n.Tag == "input" || n.Tag == "textarea" {
		n.cursor = len(n.Attrs["value"])
	}
	if !n.Style.display() || n.Style.hidden() {
		return nil
	}
	for i, child := range n.Children {
		child.Parent = n
		if child.Component == nil {
			child.Component = n.Component
		}
		var prevChild *Box
		if prev != nil && i < len(prev.Children) {
			prevChild = prev.Children[i]
		}
		if err := child.build(prevChild); err != nil {
			return err
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func inside(r image.Rectangle, x, y int) bool {
	return x >= r.Min.X && x <= r.Max.X && y >= r.Min.Y && y <= r.Max.Y
}

func getState(rect image.Rectangle) State {
	x, y := ebiten.CursorPosition()
	if inside(rect, x, y) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return Active
		}
		return Hover
	}
	return Idle
}

func (n *Box) visit(depth int, f func(depth int, n *Box) error) error {
	if n == nil {
		return nil
	}
	if err := f(depth, n); err != nil {
		return err
	}
	for _, c := range n.Children {
		if err := c.visit(depth+1, f); err != nil {
			return err
		}
	}
	return nil
}

/*
Bento is an XML based UI builder for Ebiten
*/
package bento

import (
	"bytes"
	"fmt"
	"image"
	"reflect"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"encoding/xml"
)

type Component interface {
	UI() string
}

type Box struct {
	tag      string
	children []*Box
	style    *Style
	layout
	state
	parent    *Box
	component Component
	debug     bool
	attrs     map[string]string
	content   string
}

type layout struct {
	X, Y                        int `xml:",attr"`
	ContentWidth, ContentHeight int `xml:",attr"`
	InnerWidth, InnerHeight     int `xml:",attr"`
	OuterWidth, OuterHeight     int `xml:",attr"`
}

func (n *Box) visit(f func(n *Box) error) error {
	if err := f(n); err != nil {
		return err
	}
	for _, c := range n.children {
		if err := c.visit(f); err != nil {
			return err
		}
	}
	return nil
}

func Build(c Component) (*Box, error) {
	root := &Box{
		component: c,
	}
	if err := root.rebuild(); err != nil {
		return nil, err
	}
	return root, nil
}

func (old *Box) rebuild() error {
	n, err := inflateTemplate(old.component)
	if err != nil {
		return err
	}
	if n.tag != old.tag {
		*old = *n
	}
	return nil
}

func (n *Box) isSubcomponent() bool {
	r, _ := utf8.DecodeRuneInString(n.tag)
	return unicode.IsUpper(r)
}

func inflateTemplate(c Component) (*Box, error) {
	tmpl, err := template.New("").Parse(c.UI())
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, c); err != nil {
		return nil, err
	}
	n := &Box{}
	if err := xml.Unmarshal(buf.Bytes(), n); err != nil {
		return nil, err
	}
	n.inflate(c)
	return n, nil
}

func (n *Box) inflate(c Component) error {
	n.component = c
	if n.isSubcomponent() {
		if err := n.buildSubcomponent(); err != nil {
			return err
		}
	}
	if n.style == nil {
		n.style = new(Style)
	}
	n.style.adopt(n)
	if err := n.style.parseAttributes(); err != nil {
		return err
	}
	if n.tag == "input" || n.tag == "textarea" {
		n.cursor = len(n.attrs["value"])
	}
	if !n.style.display() || n.style.hidden() {
		return nil
	}
	for _, child := range n.children {
		child.inflate(c)
	}
	return nil
}

func (n *Box) buildSubcomponent() error {
	m := reflect.ValueOf(n.component).MethodByName(n.tag)
	if !m.IsValid() {
		return fmt.Errorf("%s: failed to find method for tag %s", reflect.TypeOf(n.component), n.tag)
	}
	res := m.Call(nil)
	if style, ok := res[0].Interface().(*Style); ok {
		n.style = style
		n.tag = style.Extends
	} else if child, ok := res[0].Interface().(Component); ok {
		b := &Box{component: child}
		if err := b.rebuild(); err != nil {
			return err
		}
		for name, value := range n.attrs {
			b.attrs[name] = value
		}
		b.parent = n.parent
		b.component = child
		*n = *b
	}
	return nil
}

var keys []ebiten.Key

func (n *Box) Update() error {
	keys = inpututil.AppendPressedKeys(keys)
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyD) {
		n.toggleDebug()
	}
	n.updateState(keys)
	keys = keys[:0]
	if err := n.rebuild(); err != nil {
		return err
	}
	n.size()
	n.grow()
	n.justify()
	return nil
}

func (n *Box) toggleDebug() {
	n.visit(func(n *Box) error {
		n.debug = !n.debug
		return nil
	})
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

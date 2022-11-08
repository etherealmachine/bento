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
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func Build(c Component) (*Box, error) {
	root := &Box{
		component: c,
	}
	if err := root.build(nil); err != nil {
		return nil, err
	}
	return root, nil
}

func (n *Box) isSubcomponent() bool {
	r, _ := utf8.DecodeRuneInString(n.tag)
	return unicode.IsUpper(r)
}

func buildSubcomponent(name string, component interface{}) (*Box, error) {
	m := reflect.ValueOf(component).MethodByName(name)
	if !m.IsValid() {
		return nil, fmt.Errorf("%s has no method named %s", reflect.TypeOf(component), name)
	}
	res := m.Call(nil)
	if style, ok := res[0].Interface().(*Style); ok {
		return &Box{
			tag:   style.Extends,
			style: style,
		}, nil
	} else if sub, ok := res[0].Interface().(Component); ok {
		return Build(sub)
	}
	return nil, fmt.Errorf("%s.%s must return either Style or Component", reflect.TypeOf(component), name)
}

func (n *Box) build(prev *Box) error {
	if n.tag == "" {
		tmpl, err := template.New("").Parse(n.component.UI())
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, n.component); err != nil {
			return err
		}
		if err := xml.Unmarshal(buf.Bytes(), n); err != nil {
			return err
		}
	}
	if n.isSubcomponent() {
		subNode, err := buildSubcomponent(n.tag, n.component)
		if err != nil {
			return err
		}
		*n = *subNode
	}
	if prev != nil {
		n.state = prev.state
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
	for i, child := range n.children {
		child.parent = n
		if child.component == nil {
			child.component = n.component
		}
		var prevChild *Box
		if prev != nil && i < len(prev.children) {
			prevChild = prev.children[i]
		}
		if err := child.build(prevChild); err != nil {
			return err
		}
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
	new := &Box{component: n.component}
	if err := new.build(n); err != nil {
		return err
	}
	*n = *new
	n.size()
	n.grow()
	n.justify()
	return nil
}

func (n *Box) toggleDebug() {
	n.visit(0, func(_ int, n *Box) error {
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

func (n *Box) visit(depth int, f func(depth int, n *Box) error) error {
	if n == nil {
		return nil
	}
	if err := f(depth, n); err != nil {
		return err
	}
	for _, c := range n.children {
		if err := c.visit(depth+1, f); err != nil {
			return err
		}
	}
	return nil
}

func (n *Box) String() string {
	if n.parent != nil {
		return n.parent.String()
	}
	buf := new(bytes.Buffer)
	n.visit(0, func(depth int, n *Box) error {
		for i := 0; i < depth; i++ {
			buf.WriteByte('\t')
		}
		buf.WriteString(n.tag)
		if n.parent == nil || n.component != n.parent.component {
			buf.WriteByte(' ')
			buf.WriteString(reflect.ValueOf(n.component).Elem().Type().Name())
		}
		if content := strings.TrimSpace(n.content); content != "" {
			buf.WriteByte(' ')
			buf.WriteString(content)
		}
		buf.WriteByte('\n')
		return nil
	})
	return buf.String()
}

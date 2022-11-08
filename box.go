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

func Build(c Component) (*Box, error) {
	root := &Box{
		component: c,
	}
	if err := root.build(nil); err != nil {
		return nil, err
	}
	return root, nil
}

func (n *Box) build(prev *Box) error {
	if n.parent == nil || (prev != nil && n.component != prev.component) {
		if prev != nil && n.isSubcomponent() && prev.componentType() == n.tag {
			n.component = prev.component
		}
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
		if err := n.buildSubcomponent(prev); err != nil {
			return err
		}
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

func (n *Box) isSubcomponent() bool {
	r, _ := utf8.DecodeRuneInString(n.tag)
	return unicode.IsUpper(r)
}

func (n *Box) componentType() string {
	if n == nil || n.component == nil {
		return "<nil>"
	}
	return reflect.ValueOf(n.component).Elem().Type().Name()
}

func (n *Box) buildSubcomponent(prev *Box) error {
	m := reflect.ValueOf(n.component).MethodByName(n.tag)
	if !m.IsValid() {
		return fmt.Errorf("%s: failed to find method for tag %s", reflect.TypeOf(n.component), n.tag)
	}
	res := m.Call(nil)
	if style, ok := res[0].Interface().(*Style); ok {
		n.style = style
		n.tag = style.Extends
	} else if child, ok := res[0].Interface().(Component); ok {
		b := &Box{
			component: child,
		}
		if err := b.build(prev); err != nil {
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

func (n *Box) String() string {
	if n.parent != nil {
		return n.parent.String()
	}
	buf := new(bytes.Buffer)
	n.visit(0, func(depth int, n *Box) error {
		for i := 0; i < depth; i++ {
			buf.WriteByte('\t')
		}
		content := strings.TrimSpace(n.content)
		if content != "" {
			fmt.Fprintf(buf, "%s %s %s\n", n.tag, n.componentType(), content)
		} else {
			fmt.Fprintf(buf, "%s %s\n", n.tag, n.componentType())
		}
		return nil
	})
	return buf.String()
}

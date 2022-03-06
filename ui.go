/*
	Bento is an XML based UI builder for Ebiten
*/
package bento

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"encoding/xml"
)

type Component interface {
	UI() string
	OnKeyDown(key ebiten.Key) bool
}

type Box struct {
	tag      string
	children []*Box
	style    *Style
	layout
	state
	parent      *Box
	component   Component
	context     interface{}
	repeat      *Box
	debug       bool
	attrs       map[string]string
	attrTmpls   map[string]*template.Template
	content     string
	contentTmpl *template.Template
}

type ButtonState int

const (
	Idle     = ButtonState(0)
	Hover    = ButtonState(1)
	Active   = ButtonState(2)
	Disabled = ButtonState(3)
)

func (n *Box) clone(parent *Box) *Box {
	attrs := make(map[string]string)
	for k, v := range n.attrs {
		attrs[k] = v
	}
	if parent == nil {
		parent = n.parent
	}
	clone := &Box{
		tag:         n.tag,
		parent:      parent,
		attrs:       attrs,
		style:       n.style,
		component:   n.component,
		context:     n.context,
		attrTmpls:   n.attrTmpls,
		content:     n.content,
		contentTmpl: n.contentTmpl,
	}
	clone.children = make([]*Box, len(n.children))
	for i, c := range n.children {
		clone.children[i] = c.clone(clone)
	}
	return clone
}

type layout struct {
	X, Y                        int `xml:",attr"`
	ContentWidth, ContentHeight int `xml:",attr"`
	InnerWidth, InnerHeight     int `xml:",attr"`
	OuterWidth, OuterHeight     int `xml:",attr"`
}

type state struct {
	buttonState    ButtonState
	inputState     ButtonState
	cursorPosition int
	scrollPosition int
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
	root := &Box{}
	root.component = c
	root.context = c
	if err := xml.Unmarshal([]byte(c.UI()), root); err != nil {
		return nil, err
	}
	return root, root.visit(func(n *Box) error {
		if n.component == nil {
			n.component = n.parent.component
		}
		if n.context == nil {
			n.context = n.parent.context
		}
		if r, _ := utf8.DecodeRuneInString(n.tag); unicode.IsUpper(r) {
			m := reflect.ValueOf(n.component).MethodByName(n.tag)
			if !m.IsValid() {
				return fmt.Errorf("%s: failed to find method for tag %s", reflect.TypeOf(n.component), n.tag)
			}
			res := m.Call(nil)
			if style, ok := res[0].Interface().(*Style); ok {
				n.style = style
				n.tag = style.Extends
			}
			if box, ok := res[0].Interface().(*Box); ok {
				n.tag = "col"
				n.children = append(n.children, box)
				box.parent = n
			}
		}
		if n.style == nil {
			n.style = new(Style)
		}
		n.style.adopt(n)
		if err := n.style.parseAttributes(); err != nil {
			return err
		}
		if strings.Contains(n.content, "{{") {
			tmpl := template.New("")
			var err error
			n.contentTmpl, err = tmpl.Parse(n.content)
			if err != nil {
				log.Fatalf("failed to parse template %s: %s", n.content, err)
			}
		}
		n.attrTmpls = make(map[string]*template.Template)
		for k, v := range n.attrs {
			if strings.Contains(v, "{{") {
				tmpl := template.New("")
				var err error
				tmpl, err = tmpl.Parse(v)
				if err != nil {
					log.Fatal(err)
				}
				n.attrTmpls[k] = tmpl
			}
		}
		if repeat := n.attrs["repeat"]; repeat != "" {
			n.repeat = n.children[0]
			v := reflect.ValueOf(n.context).Elem().FieldByName(repeat)
			n.children = make([]*Box, v.Len())
			for i := 0; i < v.Len(); i++ {
				val := v.Index(i)
				clone := n.repeat.clone(nil)
				ctx := make(map[string]interface{})
				ctx["item"] = val.Interface()
				ctx["index"] = i
				ctx["parent"] = n.context
				clone.context = ctx
				n.children[i] = clone
			}
		}
		return nil
	})
}

func (n *Box) path() string {
	if n.parent == nil {
		return n.tag
	}
	return n.parent.path() + "->" + n.tag
}

func (n *Box) templateContent() string {
	if n.contentTmpl != nil {
		buf := new(bytes.Buffer)
		if err := n.contentTmpl.Execute(buf, n.context); err != nil {
			return ""
		}
		content := buf.String()
		if content == "<nil>" {
			return ""
		}
		return content
	}
	return n.content
}

func (n *Box) Update(keys []ebiten.Key) ([]ebiten.Key, error) {
	if n.parent == nil {
		n.size()
		n.grow()
		n.justify()
	}
	if len(keys) == 0 {
		for _, k := range keys {
			if inpututil.IsKeyJustPressed(k) {
				if !n.component.OnKeyDown(k) {
					keys = append(keys, k)
				}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		n.toggleDebug()
		n.dump()
	}
	n.buttonState = buttonState(n.innerRect())
	if n.buttonState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		attr := n.attrs["onClick"]
		m := reflect.ValueOf(n.component).MethodByName(attr)
		if m.IsValid() {
			m.Call(nil)
		}
	}
	var err error
	for _, c := range n.children {
		keys, err = c.Update(keys)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}

func (n *Box) templateAttr(attr string, def bool) bool {
	if v := n.attrs[attr]; v != "" {
		if tmpl := n.attrTmpls[attr]; tmpl != nil {
			buf := new(bytes.Buffer)
			if err := tmpl.Execute(buf, n.context); err != nil {
				log.Fatal(err)
			} else {
				v = buf.String()
			}
		}
		if v == "true" {
			return true
		} else if v == "false" {
			return false
		}
		log.Fatalf("invalid value for attribute %s: %s", attr, v)
	}
	return def
}

func (n *Box) toggleDebug() {
	n.debug = !n.debug
}

func (n *Box) dump() {
	if !n.style.display() || n.style.hidden() {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: (%d,%d) %dx%d\n", n.path(), n.X, n.Y, n.OuterWidth, n.OuterHeight)
	for _, c := range n.children {
		c.dump()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func inside(r image.Rectangle, x, y int) bool {
	return x >= r.Min.X && x <= r.Max.X && y >= r.Min.Y && y <= r.Max.Y
}

func buttonState(rect image.Rectangle) ButtonState {
	x, y := ebiten.CursorPosition()
	if inside(rect, x, y) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return Active
		}
		return Hover
	}
	return Idle
}

/*
	Bento is an XML based UI builder for Ebiten
*/
package bento

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"reflect"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/etherealmachine/bento/v1/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"encoding/xml"
)

var (
	Scrollspeed = 10
)

type Component interface {
	UI() string
	OnClick(btn Box)
	OnKeyDown(key ebiten.Key) bool
}

type Box interface {
	Update(keys []ebiten.Key) ([]ebiten.Key, error)
	Draw(img *ebiten.Image)
	Bounds() image.Rectangle
}

type node struct {
	tag      string
	children []*node
	style    *Style
	layout
	parent      *node
	component   Component
	context     interface{}
	repeat      *node
	debug       bool
	attrs       map[string]string
	attrTmpls   map[string]*template.Template
	content     string
	contentTmpl *template.Template
	buffer      *ebiten.Image
}

func (n *node) clone(parent *node) *node {
	attrs := make(map[string]string)
	for k, v := range n.attrs {
		attrs[k] = v
	}
	if parent == nil {
		parent = n.parent
	}
	clone := &node{
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
	clone.children = make([]*node, len(n.children))
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
	TextBounds                  *image.Rectangle
}

func (n *node) visit(f func(n *node) error) error {
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

func Build(c Component) (Box, error) {
	root := &node{}
	root.component = c
	root.context = c
	if err := xml.Unmarshal([]byte(c.UI()), root); err != nil {
		return nil, err
	}
	return root, root.visit(func(n *node) error {
		n.component = c
		if n.context == nil {
			n.context = n.parent.context
		}
		if r, _ := utf8.DecodeRuneInString(n.tag); unicode.IsUpper(r) {
			m := reflect.ValueOf(n.component).MethodByName(n.tag)
			res := m.Call(nil)
			n.style = res[0].Interface().(*Style)
			n.tag = n.style.Extends
			if n.style.FontName != "" {
				n.style.Font = text.Font(n.style.FontName, n.style.FontSize)
			} else {
				n.style.Font = text.Font("sans", 16)
			}
		} else {
			n.style = &Style{
				Font:  text.Font("sans", 16),
				Color: &color.RGBA{A: 255},
			}
		}
		if err := n.style.adopt(n.attrs); err != nil {
			return err
		}
		if strings.Contains(n.content, "{{") {
			tmpl := template.New("")
			var err error
			n.contentTmpl, err = tmpl.Parse(n.content)
			if err != nil {
				log.Fatal(err)
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
			n.children = make([]*node, v.Len())
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

func (n *node) path() string {
	if n.parent == nil {
		return n.tag
	}
	return n.parent.path() + "->" + n.tag
}

func (n *node) templateContent() string {
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

func (n *node) Update(keys []ebiten.Key) ([]ebiten.Key, error) {
	//x, y := ebiten.CursorPosition()
	var unconsumedKeys []ebiten.Key
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			if !n.component.OnKeyDown(k) {
				unconsumedKeys = append(unconsumedKeys, k)
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		n.toggleDebug()
	}
	return unconsumedKeys, nil
}

func (n *node) templateAttr(attr string, def bool) bool {
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

func (n *node) margin() (int, int, int, int) {
	if n.style != nil {
		m := n.style.Margin
		if m != nil {
			return m.Top, m.Right, m.Bottom, m.Left
		}
	}
	return 0, 0, 0, 0
}

func (n *node) padding() (int, int, int, int) {
	if n.style != nil {
		p := n.style.Padding
		if p != nil {
			return p.Top, p.Right, p.Bottom, p.Left
		}
	}
	return 0, 0, 0, 0
}

func (n *node) toggleDebug() {
	n.visit(func(n *node) error {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

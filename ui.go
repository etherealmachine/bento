/*
	Bento is an XML based UI builder for Ebiten
*/
package bento

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"os"
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
	OnClick(btn *Node)
	OnKeyDown(key ebiten.Key) bool
}

type Node struct {
	XMLName xml.Name
	Layout
	State
	Parent    *Node             `xml:"-"`
	Children  []*Node           `xml:",any"`
	CharData  string            `xml:",chardata"`
	Attrs     map[string]string `xml:"-"`
	Style     *Style            `xml:",omitempty"`
	Component Component         `xml:"-"`
	Context   interface{}       `xml:"-"`
	Repeat    *Node
	Debug     bool                          `xml:",attr"`
	attrTmpls map[string]*template.Template `xml:"-"`
	tmpl      *template.Template            `xml:"-"`
	buffer    *ebiten.Image
}

func (n *Node) Clone(parent *Node) *Node {
	attrs := make(map[string]string)
	for k, v := range n.Attrs {
		attrs[k] = v
	}
	if parent == nil {
		parent = n.Parent
	}
	clone := &Node{
		XMLName: xml.Name{
			Local: n.XMLName.Local,
			Space: n.XMLName.Space,
		},
		Parent:    parent,
		CharData:  n.CharData,
		Attrs:     attrs,
		Style:     n.Style,
		Component: n.Component,
		Context:   n.Context,
		attrTmpls: n.attrTmpls,
		tmpl:      n.tmpl,
	}
	clone.Children = make([]*Node, len(n.Children))
	for i, c := range n.Children {
		clone.Children[i] = c.Clone(clone)
	}
	return clone
}

type Layout struct {
	X, Y                        int `xml:",attr"`
	ContentWidth, ContentHeight int `xml:",attr"`
	InnerWidth, InnerHeight     int `xml:",attr"`
	OuterWidth, OuterHeight     int `xml:",attr"`
	TextBounds                  *image.Rectangle
}

type State struct {
	Hover    bool `xml:",attr,omitempty"`
	Active   bool `xml:",attr,omitempty"`
	Disabled bool `xml:",attr,omitempty"`
	Hidden   bool `xml:",attr,omitempty"`
	Display  bool `xml:",attr,omitempty"`
	Scroll   int  `xml:",attr,omitempty"`
}

func (n *Node) Visit(f func(n *Node) error) error {
	if err := f(n); err != nil {
		return err
	}
	for _, c := range n.Children {
		if err := c.Visit(f); err != nil {
			return err
		}
	}
	return nil
}

func Build(c Component) (*Node, error) {
	root := &Node{}
	root.Children = nil
	root.Component = c
	root.Context = c
	if err := xml.Unmarshal([]byte(c.UI()), root); err != nil {
		return nil, err
	}
	return root, root.Visit(func(n *Node) error {
		n.Component = c
		if n.Context == nil {
			n.Context = n.Parent.Context
		}
		if r, _ := utf8.DecodeRuneInString(n.XMLName.Local); unicode.IsUpper(r) {
			m := reflect.ValueOf(n.Component).MethodByName(n.XMLName.Local)
			res := m.Call(nil)
			n.Style = res[0].Interface().(*Style)
			n.XMLName.Local = n.Style.Extends
			if n.Style.FontName != "" {
				n.Style.Font = text.Font(n.Style.FontName, n.Style.FontSize)
			} else {
				n.Style.Font = text.Font("sans", 16)
			}
		} else {
			n.Style = &Style{
				Font:  text.Font("sans", 16),
				Color: &color.RGBA{A: 255},
			}
		}
		if err := n.Style.Adopt(n.Attrs); err != nil {
			return err
		}
		if strings.Contains(n.CharData, "{{") {
			tmpl := template.New("")
			var err error
			n.tmpl, err = tmpl.Parse(n.CharData)
			if err != nil {
				log.Fatal(err)
			}
		}
		n.attrTmpls = make(map[string]*template.Template)
		for k, v := range n.Attrs {
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
		if repeat := n.Attrs["repeat"]; repeat != "" {
			n.Repeat = n.Children[0]
			v := reflect.ValueOf(n.Context).Elem().FieldByName(repeat)
			n.Children = make([]*Node, v.Len())
			for i := 0; i < v.Len(); i++ {
				val := v.Index(i)
				clone := n.Repeat.Clone(nil)
				ctx := make(map[string]interface{})
				ctx["item"] = val.Interface()
				ctx["index"] = i
				ctx["parent"] = n.Context
				clone.Context = ctx
				n.Children[i] = clone
			}
		}
		return nil
	})
}

func (n *Node) Path() string {
	if n.Parent == nil {
		return n.XMLName.Local
	}
	return n.Parent.Path() + "->" + n.XMLName.Local
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type node Node
	if err := d.DecodeElement((*node)(n), &start); err != nil {
		return err
	}
	n.CharData = strings.TrimSpace(n.CharData)
	n.Attrs = make(map[string]string)
	for _, attr := range start.Attr {
		n.Attrs[attr.Name.Local] = attr.Value
	}
	for _, c := range n.Children {
		c.Parent = n
	}
	return nil
}

func (n *Node) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = n.XMLName
	var context string
	if n.Context != nil {
		context = reflect.TypeOf(n.Context).String()
	}
	start.Attr = append(start.Attr,
		xml.Attr{
			Name:  xml.Name{Local: "Context"},
			Value: context,
		})
	type node Node
	return e.EncodeElement((*node)(n), start)
}

func (n *Node) Content() string {
	if n.tmpl != nil {
		buf := new(bytes.Buffer)
		if err := n.tmpl.Execute(buf, n.Context); err != nil {
			return ""
		}
		content := buf.String()
		if content == "<nil>" {
			return ""
		}
		return content
	}
	return n.CharData
}

func (n *Node) Update(keys []ebiten.Key) ([]ebiten.Key, error) {
	x, y := ebiten.CursorPosition()
	root := n
	n.Visit(func(n *Node) error {
		n.Hover = false
		n.Active = false

		n.Disabled = n.templateAttr("disabled", false)
		n.Hidden = n.templateAttr("hidden", false)
		n.Display = n.templateAttr("display", true)

		if n.Disabled {
			return nil
		}

		pt, _, _, pl := n.padding()
		if x >= n.X+pl && x <= n.X+pl+n.InnerWidth && y >= n.Y+pt && y <= n.Y+pt+n.InnerHeight {
			if len(n.Children) == 0 {
				n.Hover = true
				if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
					n.Active = true
				}
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					root.Component.OnClick(n)
				}
				_, dy := ebiten.Wheel()
				if n.TextBounds != nil {
					n.Scroll += int(dy) * Scrollspeed
					if n.Scroll < 0 {
						n.Scroll = 0
					}
					if n.Scroll >= n.TextBounds.Dy()-n.ContentHeight {
						n.Scroll = n.TextBounds.Dy() - n.ContentHeight
					}
				}
			}
		}
		return nil
	})
	var unconsumedKeys []ebiten.Key
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			if !root.Component.OnKeyDown(k) {
				unconsumedKeys = append(unconsumedKeys, k)
			}
		}
	}
	return unconsumedKeys, nil
}

func (n *Node) templateAttr(attr string, def bool) bool {
	if v := n.Attrs[attr]; v != "" {
		if tmpl := n.attrTmpls[attr]; tmpl != nil {
			buf := new(bytes.Buffer)
			if err := tmpl.Execute(buf, n.Context); err != nil {
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

func (n *Node) margin() (int, int, int, int) {
	if n.Style != nil {
		m := n.Style.Margin
		if m != nil {
			return m.Top, m.Right, m.Bottom, m.Left
		}
	}
	return 0, 0, 0, 0
}

func (n *Node) padding() (int, int, int, int) {
	if n.Style != nil {
		p := n.Style.Padding
		if p != nil {
			return p.Top, p.Right, p.Bottom, p.Left
		}
	}
	return 0, 0, 0, 0
}

func (n *Node) ToggleDebug() {
	n.Visit(func(n *Node) error {
		n.Debug = !n.Debug
		return nil
	})
}

func (n *Node) Dump(out string) {
	if f, err := os.Create(out); err == nil {
		enc := xml.NewEncoder(f)
		enc.Indent("", "  ")
		if err := enc.Encode(n); err != nil {
			log.Println(err)
		}
		f.Close()
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

package bento

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Component interface {
	UI() string
}

type Box struct {
	Tag       string
	Parent    *Box
	Children  []*Box
	Style     *Style
	Debug     bool
	Attrs     map[string]string
	Content   string
	Component Component
	Layout
	state
}

func Build(c Component) (*Box, error) {
	root := &Box{
		Component: c,
	}
	if err := root.build(nil); err != nil {
		return nil, err
	}
	return root, nil
}

var keys []ebiten.Key

func (n *Box) Update() error {
	keys = inpututil.AppendPressedKeys(keys)
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyD) {
		n.ToggleDebug()
	}
	n.updateState(keys)
	keys = keys[:0]
	return n.Rebuild()
}

func (n *Box) Rebuild() error {
	new := &Box{Component: n.Component}
	if err := new.build(n); err != nil {
		return err
	}
	*n = *new
	n.size()
	n.grow()
	n.justify()
	return nil
}

func (n *Box) ToggleDebug() {
	if n.Parent != nil {
		n.Parent.ToggleDebug()
		return
	}
	n.visit(0, func(_ int, n *Box) error {
		n.Debug = !n.Debug
		return nil
	})
	n.Rebuild()
}

func (n *Box) String() string {
	if n.Parent != nil {
		return n.Parent.String()
	}
	buf := new(bytes.Buffer)
	n.visit(0, func(depth int, n *Box) error {
		for i := 0; i < depth; i++ {
			buf.WriteByte('\t')
		}
		buf.WriteString(n.Tag)
		buf.WriteByte(' ')
		if n.Debug {
			buf.WriteString("Debug")
			buf.WriteByte(' ')
		}
		if n.Parent == nil || n.Component != n.Parent.Component {
			fmt.Fprintf(buf, "<%s>", reflect.ValueOf(n.Component).Elem().Type().Name())
			buf.WriteByte(' ')
		}
		if content := strings.TrimSpace(n.Content); content != "" {
			if len(n.Content) > 20 {
				buf.WriteString(strings.TrimSpace(content[:20]))
				buf.WriteString("...")
			} else {
				buf.WriteString(content)
			}
		}
		buf.WriteByte('\n')
		return nil
	})
	return buf.String()
}

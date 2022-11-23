package bento

import (
	"bytes"
	"fmt"
	"image"
	"reflect"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Component interface {
	UI() string
}

type State int

const (
	idle     = State(0)
	hover    = State(1)
	active   = State(2)
	disabled = State(3)
)

type Box struct {
	Tag        string
	Parent     *Box
	Children   []*Box
	Debug      bool
	Content    string
	Component  Component
	Attrs      map[string]string
	state      State
	style      *Style
	scrollable Scrollable
	editable   *Editable
	dirty      bool
	layout
}

func Build(c Component) (*Box, error) {
	root := &Box{
		Component: c,
		dirty:     true,
	}
	if err := root.build(nil); err != nil {
		return nil, err
	}
	return root, nil
}

var keys []ebiten.Key

func (n *Box) Update() error {
	if n.Parent == nil {
		keys = inpututil.AppendPressedKeys(keys)
		if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyD) {
			n.ToggleDebug()
		}
	}
	if !n.style.Display || n.style.Hidden {
		return nil
	}
	n.state = idle
	if n.Attrs["disabled"] == "true" {
		n.state = disabled
	} else {
		x, y := ebiten.CursorPosition()
		if inside(n.innerRect(), x, y) {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				n.state = active
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					n.fireEvent(Click, "")
				} else {
					n.state = hover
					n.fireEvent(Hover, "")
				}
			} else {
				n.state = hover
				n.fireEvent(Hover, "")
			}
		}
	}
	if err := n.editable.Update(n); err != nil {
		return err
	}
	if err := n.scrollable.Update(n); err != nil {
		return err
	}
	n.fireEvent(Update, "")
	for _, child := range n.Children {
		if err := child.Update(); err != nil {
			return err
		}
	}
	if n.Parent == nil {
		keys = keys[:0]
		if n.dirty {
			return n.Rebuild()
		}
	}
	return nil
}

func (n *Box) Rebuild() error {
	new := &Box{
		Component: n.Component,
	}
	if err := new.build(n); err != nil {
		return err
	}
	*n = *new
	for _, child := range n.Children {
		child.Parent = n
	}
	n.size()
	n.grow()
	n.justify()
	n.dirty = false
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
}

func (n *Box) String() string {
	buf := new(bytes.Buffer)
	n.visit(0, func(depth int, n *Box) error {
		for i := 0; i < depth; i++ {
			buf.WriteByte('\t')
		}
		row := []string{n.Tag}
		if n.Debug {
			row = append(row, "Debug")
		}
		if n.Component == nil {
			row = append(row, "<nil>")
		} else if n.Parent == nil || n.Component != n.Parent.Component {
			row = append(row, fmt.Sprintf("<%s>", reflect.ValueOf(n.Component).Elem().Type().Name()))
		}
		if n.Content != "" {
			if len(n.Content) > 20 {
				row = append(row, fmt.Sprintf("%q...", strings.TrimSpace(n.Content[:20])))
			} else {
				row = append(row, fmt.Sprintf("%q", n.Content))
			}
		}
		buf.WriteString(strings.Join(row, " "))
		buf.WriteByte('\n')
		return nil
	})
	return buf.String()
}

func (n *Box) root() *Box {
	if n.Parent == nil {
		return n
	}
	return n.Parent.root()
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

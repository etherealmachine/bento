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
	Idle     = State(0)
	Hover    = State(1)
	Active   = State(2)
	Disabled = State(3)
)

type Box struct {
	Tag        string
	Parent     *Box
	Children   []*Box
	Debug      bool
	Content    string
	Component  Component
	State      State
	style      *Style
	attrs      map[string]string
	scrollable Scrollable
	editable   *Editable
	layout
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
	if n.Parent == nil {
		keys = inpututil.AppendPressedKeys(keys)
		if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyD) {
			n.ToggleDebug()
		}
	}
	if !n.style.display() || n.style.hidden() {
		return nil
	}
	n.State = Idle
	if n.attrs["disabled"] == "true" {
		n.State = Disabled
	} else {
		x, y := ebiten.CursorPosition()
		if inside(n.innerRect(), x, y) {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				n.State = Active
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					n.fireEvent("Click")
				}
			} else {
				n.State = Hover
				n.fireEvent("Hover")
			}
		}
	}
	if err := n.editable.Update(n); err != nil {
		return err
	}
	if err := n.scrollable.Update(n); err != nil {
		return err
	}
	n.fireEvent(Update)
	for _, child := range n.Children {
		if err := child.Update(); err != nil {
			return err
		}
	}
	if n.Parent == nil {
		keys = keys[:0]
		return n.Rebuild()
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
				row = append(row, strings.TrimSpace(n.Content[:20])+"...")
			} else {
				row = append(row, n.Content)
			}
		}
		buf.WriteString(strings.Join(row, " "))
		buf.WriteByte('\n')
		return nil
	})
	return buf.String()
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

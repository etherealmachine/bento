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

var (
	debug bool
)

type Box struct {
	Tag        string
	Parent     *Box
	Children   []*Box
	Content    string
	Component  Component
	Attrs      map[string]string
	Style      Style
	state      State
	scrollable Scrollable
	editable   *Editable
	dirty      bool
	target     *ebiten.Image
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

type context struct {
	keys     []ebiten.Key
	consumed bool
}

var ctx context

func (n *Box) Update() error {
	if n.Parent != nil {
		return fmt.Errorf("Update called on non-root element %s", n.Tag)
	}
	ctx.consumed = false
	return n.update(&ctx)
}

func (n *Box) update(ctx *context) error {
	if !n.Style.Display || n.Style.Hidden {
		return nil
	}
	if n.Parent == nil {
		ctx.keys = inpututil.AppendPressedKeys(ctx.keys)
		if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyD) {
			n.ToggleDebug()
		}
	}
	// reverse draw order so that highest zIndex consumes events first
	for i := range n.Children {
		if err := n.Children[len(n.Children)-i-1].update(ctx); err != nil {
			return err
		}
	}
	n.state = idle
	if n.Attrs["disabled"] == "true" {
		n.state = disabled
	} else if x, y := ebiten.CursorPosition(); !ctx.consumed && inside(n.innerRect(), x, y) {
		if sx, sy := ebiten.Wheel(); sx != 0 || sy != 0 {
			ctx.consumed = n.fireEvent(Scroll, "", nil, nil)
		} else {
			switch {
			case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
				n.state = active
				ctx.consumed = n.fireEvent(Click, "", nil, nil)
			case ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft):
				n.state = active
				n.fireEvent(Hover, "", nil, nil)
				ctx.consumed = true
			default:
				n.state = hover
				n.fireEvent(Hover, "", nil, nil)
				ctx.consumed = true
			}
		}
	}
	if err := n.editable.update(n, ctx); err != nil {
		return err
	}
	if err := n.scrollable.update(n); err != nil {
		return err
	}
	n.fireEvent(Update, "", nil, nil)
	if n.Parent == nil {
		ctx.keys = ctx.keys[:0]
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
	n.relayout()
	n.dirty = false
	return nil
}

func (n *Box) ToggleDebug() {
	debug = !debug
}

func (n *Box) String() string {
	buf := new(bytes.Buffer)
	n.visit(0, func(depth int, n *Box) error {
		for i := 0; i < depth; i++ {
			buf.WriteByte('\t')
		}
		row := []string{n.Tag}
		if n.Component == nil {
			row = append(row, "<nil>")
		} else if n.Parent == nil || n.Component != n.Parent.Component {
			row = append(row, fmt.Sprintf("<%s>", reflect.ValueOf(n.Component).Elem().Type().Name()))
		}
		text := n.Content
		if text == "" {
			text = n.Attrs["value"]
		}
		if text != "" {
			if len(text) > 20 {
				row = append(row, fmt.Sprintf("%q...", strings.TrimSpace(text[:20])))
			} else {
				row = append(row, fmt.Sprintf("%q", text))
			}
		}
		if n.editable != nil {
			if n.editable.focus {
				row = append(row, "focus=true")
			} else {
				row = append(row, "focus=false")
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

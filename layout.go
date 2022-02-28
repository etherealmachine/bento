package bento

import (
	"fmt"
	"image"
	"strings"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
)

func (n *node) Bounds() image.Rectangle {
	return n.OuterRect()
}

func (n *node) OuterRect() image.Rectangle {
	return image.Rect(n.X, n.Y, n.X+n.OuterWidth, n.Y+n.OuterHeight)
}

func (n *node) InnerRect() image.Rectangle {
	mt, _, _, ml := n.margin()
	return image.Rect(n.X+ml, n.Y+mt, n.X+ml+n.InnerWidth, n.Y+mt+n.InnerHeight)
}

func (n *node) ContentRect() image.Rectangle {
	mt, _, _, ml := n.margin()
	pt, _, _, pl := n.padding()
	return image.Rect(n.X+ml+pl, n.Y+mt+pt, n.X+ml+pl+n.ContentWidth, n.Y+mt+pt+n.ContentHeight)
}

func (n *node) updateSize() {
	n.ContentWidth = 0
	n.ContentHeight = 0
	if !n.style.Display {
		n.InnerWidth = 0
		n.InnerHeight = 0
		n.OuterWidth = 0
		n.OuterHeight = 0
		return
	}
	if n.tag == "button" || n.tag == "text" {
		bounds := text.BoundString(n.style.Font, n.templateContent())
		textHeight := bounds.Dy()
		metrics := n.style.Font.Metrics()
		minHeight := metrics.Height.Round()
		if textHeight < minHeight {
			textHeight = minHeight
		}
		n.TextBounds = &bounds
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = textHeight
	} else if n.tag == "p" {
		bounds := text.BoundParagraph(n.style.Font, n.templateContent(), n.style.MaxWidth)
		n.TextBounds = &bounds
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = max(n.style.MinHeight, bounds.Dy())
		if n.style.MaxHeight > 0 && n.style.MaxHeight < n.ContentHeight {
			n.ContentHeight = n.style.MaxHeight
		}
	} else if n.tag == "img" {
		bounds := n.style.Background.Bounds()
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	}
	for _, c := range n.children {
		c.updateSize()
		switch n.tag {
		case "grid":
			n.ContentWidth = max(n.ContentWidth, c.OuterWidth)
			n.ContentHeight += c.OuterHeight
		case "row":
			n.ContentWidth += c.OuterWidth
			n.ContentHeight = max(n.ContentHeight, c.OuterHeight)
		case "col":
			n.ContentWidth = max(n.ContentWidth, c.OuterWidth)
			n.ContentHeight += c.OuterHeight
		}
	}
	n.styleSize()
	mt, mr, mb, ml := n.margin()
	pt, pr, pb, pl := n.padding()
	n.InnerWidth = n.ContentWidth + pl + pr
	n.InnerHeight = n.ContentHeight + pt + pb
	n.OuterWidth = n.InnerWidth + ml + mr
	if n.TextBounds != nil && n.TextBounds.Dy() > n.ContentHeight {
		n.OuterWidth += n.style.Scrollbar.Bounds().Dx()
	}
	n.OuterHeight = n.InnerHeight + mt + mb
}

func (n *node) styleSize() {
	if n.style == nil {
		return
	}
	if n.style.Attrs["width"] == "100%" {
		if n.parent == nil {
			w, _ := ebiten.WindowSize()
			n.ContentWidth = w
		} else {
			n.ContentWidth = max(n.ContentWidth, n.parent.InnerWidth)
		}
	}
	if n.style.Attrs["height"] == "100%" {
		if n.parent == nil {
			_, h := ebiten.WindowSize()
			n.ContentHeight = h
		} else {
			n.ContentHeight = max(n.ContentWidth, n.parent.InnerHeight)
		}
	}
	if n.style.MinWidth > 0 {
		n.ContentWidth = max(n.ContentWidth, n.style.MinWidth)
	}
	if n.style.MinHeight > 0 {
		n.ContentHeight = max(n.ContentHeight, n.style.MinHeight)
	}
	if n.style.Attrs["aspectRatio"] == "1:1" {
		n.ContentWidth = max(n.ContentWidth, n.ContentHeight)
		n.ContentHeight = max(n.ContentWidth, n.ContentHeight)
	}
	if n.style.Button != nil {
		n.style.Button.rect = n.InnerRect()
	}
	if n.style.Scrollbar != nil {
		inner := n.InnerRect()
		n.style.Scrollbar.x = inner.Max.X
		n.style.Scrollbar.y = inner.Min.Y
		n.style.Scrollbar.height = inner.Dy()
	}
}

func (n *node) fillWidth(w int) {
	_, mr, _, ml := n.margin()
	_, pr, _, pl := n.padding()
	n.OuterWidth = w
	n.InnerWidth = n.OuterWidth - ml - mr
	n.ContentWidth = n.InnerWidth - pr - pl
}

func (n *node) fillHeight(h int) {
	mt, _, mb, _ := n.margin()
	pt, _, pb, _ := n.padding()
	n.OuterHeight = h
	n.InnerHeight = n.OuterHeight - mt - mb
	n.ContentHeight = n.InnerHeight - pt - pb
}

func (n *node) doLayout() {
	hj, vj := Start, Start
	if n.style != nil && n.style.Attrs["justify"] != "" {
		justs := strings.Split(n.style.Attrs["justify"], " ")
		if len(justs) == 1 {
			hj = ParseJustification(justs[0])
			vj = hj
		} else {
			hj, vj = ParseJustification(justs[0]), ParseJustification(justs[1])
		}
	}
	dir := n.tag
	if dir == "grid" {
		dir = "col"
	}

	maxChildWidth := 0
	maxChildHeight := 0
	totalChildWidths := 0
	totalChildHeights := 0
	for _, c := range n.children {
		maxChildWidth = max(maxChildWidth, c.OuterWidth)
		maxChildHeight = max(maxChildHeight, c.OuterHeight)
		totalChildWidths += c.OuterWidth
		totalChildHeights += c.OuterHeight
	}

	var hspace, vspace int
	if dir == "row" {
		hspace = n.InnerWidth - totalChildWidths
		vspace = n.InnerHeight - maxChildHeight
	} else {
		hspace = n.InnerWidth - maxChildWidth
		vspace = n.InnerHeight - totalChildHeights
	}

	r := n.InnerRect()

	switch hj {
	case Start:
		for i, c := range n.children {
			offset := r.Min.X
			if i > 0 && dir == "row" {
				prev := n.children[i-1]
				offset = prev.X + prev.OuterWidth
			}
			c.X = offset
		}
	case End:
		for i, c := range n.children {
			offset := r.Min.X + hspace
			if i > 0 && dir == "row" {
				prev := n.children[i-1]
				offset = prev.X + prev.OuterWidth + hspace
			}
			c.X = offset
		}
	case Center:
		for i, c := range n.children {
			offset := r.Min.X + hspace/2
			if i > 0 && dir == "row" {
				prev := n.children[i-1]
				offset = prev.X + prev.OuterWidth
			}
			c.X = offset
		}
	case Evenly:
		gap := hspace / (len(n.children) + 1)
		for i, c := range n.children {
			offset := r.Min.X + gap
			if i > 0 && dir == "row" {
				prev := n.children[i-1]
				offset = prev.X + prev.OuterWidth + gap
			}
			c.X = offset
		}
	case Stretch:
		w := n.InnerWidth
		if dir == "row" {
			w = n.InnerWidth / len(n.children)
		}
		for i, c := range n.children {
			c.X = r.Min.X
			if dir == "row" {
				c.X += w * i
			}
			c.fillWidth(w)
		}
	default:
		panic(fmt.Errorf("can't handle horizontal justification %d", hj))
	}
	switch vj {
	case Start:
		for i, c := range n.children {
			offset := r.Min.Y
			if i > 0 && dir == "col" {
				prev := n.children[i-1]
				offset = prev.Y + prev.OuterHeight
			}
			c.Y = offset
		}
	case End:
		for i, c := range n.children {
			offset := r.Min.Y + vspace
			if i > 0 && dir == "col" {
				prev := n.children[i-1]
				offset = prev.Y + prev.OuterHeight + vspace
			}
			c.Y = offset
		}
	case Center:
		for i, c := range n.children {
			offset := r.Min.Y + vspace/2
			if i > 0 && dir == "col" {
				prev := n.children[i-1]
				offset = prev.Y + prev.OuterHeight
			}
			c.Y = offset
		}
	case Evenly:
		gap := vspace / (len(n.children) + 1)
		for i, c := range n.children {
			offset := r.Min.Y + gap
			if i > 0 && dir == "col" {
				prev := n.children[i-1]
				offset = prev.Y + prev.OuterHeight + gap
			}
			c.Y = offset
		}
	case Stretch:
		h := n.InnerHeight
		if dir == "col" {
			h = n.InnerHeight / len(n.children)
		}
		for i, c := range n.children {
			c.Y = r.Min.Y
			if dir == "col" {
				c.Y += h * i
			}
			c.fillHeight(h)
		}
	default:
		panic(fmt.Errorf("can't handle vertical justification %d", vj))
	}
	for _, c := range n.children {
		c.doLayout()
	}
}

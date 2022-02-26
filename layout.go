package bento

import (
	"fmt"
	"image"
	"strings"

	"github.com/etherealmachine/bento/v1/text"
	"github.com/hajimehoshi/ebiten/v2"
)

func (n *Node) OuterRect() *image.Rectangle {
	r := image.Rect(n.X, n.Y, n.X+n.OuterWidth, n.Y+n.OuterHeight)
	return &r
}

func (n *Node) InnerRect() *image.Rectangle {
	mt, _, _, ml := n.margin()
	r := image.Rect(n.X+ml, n.Y+mt, n.X+ml+n.InnerWidth, n.Y+mt+n.InnerHeight)
	return &r
}

func (n *Node) ContentRect() *image.Rectangle {
	mt, _, _, ml := n.margin()
	pt, _, _, pl := n.padding()
	r := image.Rect(n.X+ml+pl, n.Y+mt+pt, n.X+ml+pl+n.ContentWidth, n.Y+mt+pt+n.ContentHeight)
	return &r
}

func (n *Node) updateSize() {
	n.ContentWidth = 0
	n.ContentHeight = 0
	if !n.Display {
		n.InnerWidth = 0
		n.InnerHeight = 0
		n.OuterWidth = 0
		n.OuterHeight = 0
		return
	}
	if n.XMLName.Local == "button" || n.XMLName.Local == "text" {
		bounds := text.BoundString(n.Style.Font, n.Content())
		textHeight := bounds.Dy()
		metrics := n.Style.Font.Metrics()
		minHeight := metrics.Height.Round()
		if textHeight < minHeight {
			textHeight = minHeight
		}
		n.TextBounds = &bounds
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = textHeight
	} else if n.XMLName.Local == "p" {
		bounds := text.BoundParagraph(n.Style.Font, n.Content(), n.Style.MaxWidth)
		n.TextBounds = &bounds
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = max(n.Style.MinHeight, bounds.Dy())
		if n.Style.MaxHeight > 0 && n.Style.MaxHeight < n.ContentHeight {
			n.ContentHeight = n.Style.MaxHeight
		}
	} else if n.XMLName.Local == "img" {
		bounds := n.Style.Background.Bounds()
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	}
	for _, c := range n.Children {
		c.updateSize()
		switch n.XMLName.Local {
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
	n.OuterHeight = n.InnerHeight + mt + mb
}

func (n *Node) styleSize() {
	if n.Style == nil {
		return
	}
	if n.Style.Attrs["width"] == "100%" {
		if n.Parent == nil {
			w, _ := ebiten.WindowSize()
			n.ContentWidth = w
		} else {
			n.ContentWidth = max(n.ContentWidth, n.Parent.InnerWidth)
		}
	}
	if n.Style.Attrs["height"] == "100%" {
		if n.Parent == nil {
			_, h := ebiten.WindowSize()
			n.ContentHeight = h
		} else {
			n.ContentHeight = max(n.ContentWidth, n.Parent.InnerHeight)
		}
	}
	if n.Style.MinWidth > 0 {
		n.ContentWidth = max(n.ContentWidth, n.Style.MinWidth)
	}
	if n.Style.MinHeight > 0 {
		n.ContentHeight = max(n.ContentHeight, n.Style.MinHeight)
	}
	if n.Style.Attrs["aspectRatio"] == "1:1" {
		n.ContentWidth = max(n.ContentWidth, n.ContentHeight)
		n.ContentHeight = max(n.ContentWidth, n.ContentHeight)
	}
}

func (n *Node) fillWidth(w int) {
	_, mr, _, ml := n.margin()
	_, pr, _, pl := n.padding()
	n.OuterWidth = w
	n.InnerWidth = n.OuterWidth - ml - mr
	n.ContentWidth = n.InnerWidth - pr - pl
}

func (n *Node) fillHeight(h int) {
	mt, _, mb, _ := n.margin()
	pt, _, pb, _ := n.padding()
	n.OuterHeight = h
	n.InnerHeight = n.OuterHeight - mt - mb
	n.ContentHeight = n.InnerHeight - pt - pb
}

func (n *Node) layout() {
	hj, vj := Start, Start
	if n.Style != nil && n.Style.Attrs["justify"] != "" {
		justs := strings.Split(n.Style.Attrs["justify"], " ")
		if len(justs) == 1 {
			hj = ParseJustification(justs[0])
			vj = hj
		} else {
			hj, vj = ParseJustification(justs[0]), ParseJustification(justs[1])
		}
	}
	dir := n.XMLName.Local
	if dir == "grid" {
		dir = "col"
	}

	maxChildWidth := 0
	maxChildHeight := 0
	totalChildWidths := 0
	totalChildHeights := 0
	for _, c := range n.Children {
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
		for i, c := range n.Children {
			offset := r.Min.X
			if i > 0 && dir == "row" {
				prev := n.Children[i-1]
				offset = prev.X + prev.OuterWidth
			}
			c.X = offset
		}
	case End:
		for i, c := range n.Children {
			offset := r.Min.X + hspace
			if i > 0 && dir == "row" {
				prev := n.Children[i-1]
				offset = prev.X + prev.OuterWidth + hspace
			}
			c.X = offset
		}
	case Center:
		for i, c := range n.Children {
			offset := r.Min.X + hspace/2
			if i > 0 && dir == "row" {
				prev := n.Children[i-1]
				offset = prev.X + prev.OuterWidth
			}
			c.X = offset
		}
	case Evenly:
		gap := hspace / (len(n.Children) + 1)
		for i, c := range n.Children {
			offset := r.Min.X + gap
			if i > 0 && dir == "row" {
				prev := n.Children[i-1]
				offset = prev.X + prev.OuterWidth + gap
			}
			c.X = offset
		}
	case Stretch:
		w := n.InnerWidth
		if dir == "row" {
			w = n.InnerWidth / len(n.Children)
		}
		for i, c := range n.Children {
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
		for i, c := range n.Children {
			offset := r.Min.Y
			if i > 0 && dir == "col" {
				prev := n.Children[i-1]
				offset = prev.Y + prev.OuterHeight
			}
			c.Y = offset
		}
	case End:
		for i, c := range n.Children {
			offset := r.Min.Y + vspace
			if i > 0 && dir == "col" {
				prev := n.Children[i-1]
				offset = prev.Y + prev.OuterHeight + vspace
			}
			c.Y = offset
		}
	case Center:
		for i, c := range n.Children {
			offset := r.Min.Y + vspace/2
			if i > 0 && dir == "col" {
				prev := n.Children[i-1]
				offset = prev.Y + prev.OuterHeight
			}
			c.Y = offset
		}
	case Evenly:
		gap := vspace / (len(n.Children) + 1)
		for i, c := range n.Children {
			offset := r.Min.Y + gap
			if i > 0 && dir == "col" {
				prev := n.Children[i-1]
				offset = prev.Y + prev.OuterHeight + gap
			}
			c.Y = offset
		}
	case Stretch:
		h := n.InnerHeight
		if dir == "col" {
			h = n.InnerHeight / len(n.Children)
		}
		for i, c := range n.Children {
			c.Y = r.Min.Y
			if dir == "col" {
				c.Y += h * i
			}
			c.fillHeight(h)
		}
	default:
		panic(fmt.Errorf("can't handle vertical justification %d", vj))
	}
	for _, c := range n.Children {
		c.layout()
	}
}

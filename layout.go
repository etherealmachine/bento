package bento

import (
	"fmt"
	"image"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
)

func (n *node) Bounds() image.Rectangle {
	return n.outerRect()
}

func (n *node) outerRect() image.Rectangle {
	return image.Rect(n.X, n.Y, n.X+n.OuterWidth, n.Y+n.OuterHeight)
}

func (n *node) innerRect() image.Rectangle {
	mt, _, _, ml := n.style.margin()
	return image.Rect(n.X+ml, n.Y+mt, n.X+ml+n.InnerWidth, n.Y+mt+n.InnerHeight)
}

func (n *node) contentRect() image.Rectangle {
	mt, _, _, ml := n.style.margin()
	pt, _, _, pl := n.style.padding()
	return image.Rect(n.X+ml+pl, n.Y+mt+pt, n.X+ml+pl+n.ContentWidth, n.Y+mt+pt+n.ContentHeight)
}

func (n *node) size() {
	n.ContentWidth = 0
	n.ContentHeight = 0
	if !n.style.display() {
		n.InnerWidth = 0
		n.InnerHeight = 0
		n.OuterWidth = 0
		n.OuterHeight = 0
		return
	}
	if n.tag == "button" || n.tag == "text" {
		bounds := text.BoundString(n.style.Font, n.templateContent())
		n.TextBounds = &bounds
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.tag == "p" {
		bounds := text.BoundParagraph(n.style.Font, n.templateContent(), n.style.MaxWidth)
		n.TextBounds = &bounds
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.tag == "img" {
		bounds := n.style.Background.Bounds()
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	}
	for _, c := range n.children {
		c.size()
		switch n.tag {
		case "row":
			n.ContentWidth += c.OuterWidth
			n.ContentHeight = max(n.ContentHeight, c.OuterHeight)
		case "col":
			n.ContentWidth = max(n.ContentWidth, c.OuterWidth)
			n.ContentHeight += c.OuterHeight
		}
	}
	n.styleSize()
	mt, mr, mb, ml := n.style.margin()
	pt, pr, pb, pl := n.style.padding()
	n.InnerWidth = n.ContentWidth + pl + pr
	n.InnerHeight = n.ContentHeight + pt + pb
	n.OuterWidth = n.InnerWidth + ml + mr
	if n.TextBounds != nil && n.TextBounds.Dy() > n.ContentHeight {
		n.OuterWidth += n.style.Scrollbar.Bounds().Dx()
	}
	n.OuterHeight = n.InnerHeight + mt + mb
}

func (n *node) grow() {
	w, h := ebiten.WindowSize()
	if n.style != nil {
		if n.style.HGrow > 0 {
			if n.parent == nil {
				n.fillWidth(w)
			} else {
				// TODO
			}
		}
		if n.style.VGrow > 0 {
			if n.parent == nil {
				n.fillHeight(h)
			} else {
				// TODO
			}
		}
	}
	for _, c := range n.children {
		c.grow()
	}
}

func (n *node) fillWidth(w int) {
	_, mr, _, ml := n.style.margin()
	_, pr, _, pl := n.style.padding()
	n.OuterWidth = w
	n.InnerWidth = n.OuterWidth - ml - mr
	n.ContentWidth = n.InnerWidth - pr - pl
}

func (n *node) fillHeight(h int) {
	mt, _, mb, _ := n.style.margin()
	pt, _, pb, _ := n.style.padding()
	n.OuterHeight = h
	n.InnerHeight = n.OuterHeight - mt - mb
	n.ContentHeight = n.InnerHeight - pt - pb
}

func (n *node) space() (int, int) {
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
	if n.tag == "row" {
		hspace = n.InnerWidth - totalChildWidths
		vspace = n.InnerHeight - maxChildHeight
	} else {
		hspace = n.InnerWidth - maxChildWidth
		vspace = n.InnerHeight - totalChildHeights
	}
	return hspace, vspace
}

func (n *node) justify() {
	/*
		r := n.innerRect()
		extents := make([]int, len(n.children))
		for i, c := range n.children {
			extents[i] = c.InnerWidth
		}
		offset := distribute(hspace, hj, extents)
		for i, c := range n.children {
			c.X = r.Min.X + offset[i]
		}
		for i, c := range n.children {
			extents[i] = c.InnerHeight
		}
		//offset = distribute(vspace, vj, extents)
		for _, c := range n.children {
			c.Y = r.Min.Y // + offset[i]
		}
	*/
	for _, c := range n.children {
		c.justify()
	}
}

func distribute(space int, j Justification, extents []int) []int {
	switch j {
	case Start:
		return extents
	case End:
		return extents
	case Center:
		return extents
	case Evenly:
		return extents
	case Around:
		return extents
	case Between:
		return extents
	default:
		panic(fmt.Errorf("can't handle justification %s", j))
	}
}

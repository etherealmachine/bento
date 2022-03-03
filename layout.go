package bento

import (
	"fmt"
	"image"
	"math"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
)

func (n *Box) Bounds() image.Rectangle {
	return n.outerRect()
}

func (n *Box) outerRect() image.Rectangle {
	return image.Rect(n.X, n.Y, n.X+n.OuterWidth, n.Y+n.OuterHeight)
}

func (n *Box) innerRect() image.Rectangle {
	mt, _, _, ml := n.style.margin()
	return image.Rect(n.X+ml, n.Y+mt, n.X+ml+n.InnerWidth, n.Y+mt+n.InnerHeight)
}

func (n *Box) contentRect() image.Rectangle {
	mt, _, _, ml := n.style.margin()
	pt, _, _, pl := n.style.padding()
	return image.Rect(n.X+ml+pl, n.Y+mt+pt, n.X+ml+pl+n.ContentWidth, n.Y+mt+pt+n.ContentHeight)
}

func (n *Box) size() {
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
	} else if n.tag == "img" && n.style.Image != nil {
		bounds := n.style.Image.Bounds()
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

func (n *Box) grow() {
	w, h := ebiten.WindowSize()
	if n.style != nil && n.style.HGrow > 0 && n.parent == nil {
		n.fillWidth(w)
	}
	if n.style != nil && n.style.VGrow > 0 && n.parent == nil {
		n.fillHeight(h)
	}
	hgrow, vgrow := 0, 0
	for _, c := range n.children {
		hg, vg := c.style.growth()
		hgrow += hg
		vgrow += vg
	}
	hspace, vspace := n.space()
	for _, c := range n.children {
		hg, vg := c.style.growth()
		if hg > 0 {
			if n.tag == "row" {
				halloc := int(math.Floor(float64(hg) / float64(hgrow) * float64(hspace)))
				c.fillWidth(c.ContentWidth + halloc)
			} else {
				c.fillWidth(n.OuterWidth)
			}
		}
		if vg > 0 && vgrow > 0 {
			if n.tag == "col" {
				valloc := int(math.Floor(float64(vg) / float64(vgrow) * float64(vspace)))
				c.fillHeight(c.ContentHeight + valloc)
			} else {
				c.fillHeight(c.OuterHeight)
			}
		}
	}
	for _, c := range n.children {
		c.grow()
	}
}

func (n *Box) fillWidth(w int) {
	_, mr, _, ml := n.style.margin()
	_, pr, _, pl := n.style.padding()
	n.OuterWidth = w
	n.InnerWidth = n.OuterWidth - ml - mr
	n.ContentWidth = n.InnerWidth - pr - pl
}

func (n *Box) fillHeight(h int) {
	mt, _, mb, _ := n.style.margin()
	pt, _, pb, _ := n.style.padding()
	n.OuterHeight = h
	n.InnerHeight = n.OuterHeight - mt - mb
	n.ContentHeight = n.InnerHeight - pt - pb
}

func (n *Box) space() (int, int) {
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

func (n *Box) justify() {
	r := n.innerRect()
	hspace, vspace := n.space()
	hj, vj := n.style.justification()
	for _, c := range n.children {
		c.X = r.Min.X
		c.Y = r.Min.Y
	}
	extents := make([]int, len(n.children))
	for i, c := range n.children {
		extents[i] = c.OuterWidth
	}
	offsets := distribute(hspace, hj, extents, n.tag == "row")
	for i, c := range n.children {
		c.X += offsets[i]
	}
	for i, c := range n.children {
		extents[i] = c.OuterHeight
	}
	offsets = distribute(vspace, vj, extents, n.tag == "col")
	for i, c := range n.children {
		c.Y += offsets[i]
	}
	for _, c := range n.children {
		c.justify()
	}
}

func distribute(space int, j Justification, extents []int, mainAxis bool) []int {
	offsets := make([]int, len(extents))
	switch j {
	case Start:
		for i := range extents {
			if i == 0 || !mainAxis {
				offsets[i] = 0
			} else {
				offsets[i] = offsets[i-1] + extents[i-1]
			}
		}
	case End:
		for i := range extents {
			if i == 0 || !mainAxis {
				offsets[i] = space
			} else {
				offsets[i] = offsets[i-1] + extents[i-1]
			}
		}
	case Center:
		for i := range extents {
			if i == 0 || !mainAxis {
				offsets[i] = space / 2
			} else {
				offsets[i] = offsets[i-1] + extents[i-1]
			}
		}
	case Evenly:
		spacing := int(math.Floor(float64(space) / float64(len(extents)+1)))
		for i := range extents {
			if !mainAxis {
				offsets[i] = space / 2
			} else if i == 0 {
				offsets[i] = spacing
			} else {
				offsets[i] = offsets[i-1] + extents[i-1] + spacing
			}
		}
	case Around:
		spacing := int(math.Floor(float64(space) / float64(len(extents))))
		for i := range extents {
			if !mainAxis {
				offsets[i] = space / 2
			} else if i == 0 {
				offsets[i] = int(math.Floor(float64(spacing) / 2))
			} else {
				offsets[i] = offsets[i-1] + extents[i-1] + spacing
			}
		}
	case Between:
		spacing := int(math.Floor(float64(space) / float64(len(extents)-1)))
		for i := range extents {
			if !mainAxis {
				offsets[i] = space / 2
			} else if i == 0 {
				offsets[i] = 0
			} else {
				offsets[i] = offsets[i-1] + extents[i-1] + spacing
			}
		}
	default:
		panic(fmt.Errorf("can't handle justification %s", j))
	}
	return offsets
}

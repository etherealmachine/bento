package bento

import (
	"fmt"
	"image"
	"log"
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
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.tag == "p" {
		bounds := text.BoundParagraph(n.style.Font, n.templateContent(), n.style.MaxWidth)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.tag == "img" && n.style.Image != nil {
		bounds := n.style.Image.Bounds()
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.tag == "input" {
		bounds := text.BoundString(n.style.Font, n.attrs["placeholder"])
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = n.style.Font.Metrics().Height.Ceil()
	} else if n.tag == "textarea" {
		bounds := text.BoundString(n.style.Font, n.attrs["value"])
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.tag != "row" && n.tag != "col" {
		log.Fatalf("can't size %s", n.tag)
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
		if !c.style.display() {
			continue
		}
		hg, vg := c.style.growth()
		hgrow += hg
		vgrow += vg
	}
	hspace, vspace := n.space()
	for _, c := range n.children {
		if !c.style.display() {
			continue
		}
		hg, vg := c.style.growth()
		if hg > 0 {
			if n.tag == "row" {
				halloc := int(math.Floor(float64(hg) / float64(hgrow) * float64(hspace)))
				c.fillWidth(c.OuterWidth + halloc)
			} else {
				c.fillWidth(n.OuterWidth)
			}
		}
		if vg > 0 && vgrow > 0 {
			if n.tag == "col" {
				valloc := int(math.Floor(float64(vg) / float64(vgrow) * float64(vspace)))
				c.fillHeight(c.OuterHeight + valloc)
			} else {
				c.fillHeight(c.OuterHeight)
			}
		}
	}
	for _, c := range n.children {
		if !c.style.display() {
			continue
		}
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
	extents := make([][2]int, len(n.children))
	for i, c := range n.children {
		if n.tag == "row" {
			extents[i][0] = c.OuterWidth
			extents[i][1] = c.OuterHeight
		} else {
			extents[i][0] = c.OuterHeight
			extents[i][1] = c.OuterWidth
		}
	}
	var offsets [][2]int
	if n.tag == "row" {
		offsets = distribute(hspace, n.InnerHeight, hj, vj, extents)
	} else {
		offsets = distribute(vspace, n.InnerWidth, vj, hj, extents)
	}
	for i, c := range n.children {
		if n.tag == "row" {
			c.X += offsets[i][0]
			c.Y += offsets[i][1]
		} else {
			c.Y += offsets[i][0]
			c.X += offsets[i][1]
		}
	}
	for _, c := range n.children {
		c.justify()
	}
}

func distribute(mainspace, crossspace int, mainj, crossj Justification, extents [][2]int) [][2]int {
	offsets := make([][2]int, len(extents))
	for i := range extents {
		switch mainj {
		case Start:
			if i == 0 {
				offsets[i][0] = 0
			} else {
				offsets[i][0] = offsets[i-1][0] + extents[i-1][0]
			}
		case End:
			if i == 0 {
				offsets[i][0] = mainspace
			} else {
				offsets[i][0] = offsets[i-1][0] + extents[i-1][0]
			}
		case Center:
			if i == 0 {
				offsets[i][0] = mainspace / 2
			} else {
				offsets[i][0] = offsets[i-1][0] + extents[i-1][0]
			}
		case Evenly:
			spacing := int(math.Floor(float64(mainspace) / float64(len(extents)+1)))
			if i == 0 {
				offsets[i][0] = spacing
			} else {
				offsets[i][0] = offsets[i-1][0] + extents[i-1][0] + spacing
			}
		case Around:
			spacing := int(math.Floor(float64(mainspace) / float64(len(extents))))
			if i == 0 {
				offsets[i][0] = int(math.Floor(float64(spacing) / 2))
			} else {
				offsets[i][0] = offsets[i-1][0] + extents[i-1][0] + spacing
			}
		case Between:
			spacing := int(math.Floor(float64(mainspace) / float64(len(extents)-1)))
			if i == 0 {
				offsets[i][0] = 0
			} else {
				offsets[i][0] = offsets[i-1][0] + extents[i-1][0] + spacing
			}
		default:
			panic(fmt.Errorf("can't handle main axis justification %s", mainj))
		}
		switch crossj {
		case Start:
			offsets[i][1] = 0
		case End:
			offsets[i][1] = crossspace - extents[i][1]
		case Center, Evenly, Around, Between:
			offsets[i][1] = int(math.Floor(float64(crossspace)/2)) - int(math.Floor(float64(extents[i][1])/2))
		default:
			panic(fmt.Errorf("can't handle cross axis justification %s", crossj))
		}
	}
	return offsets
}

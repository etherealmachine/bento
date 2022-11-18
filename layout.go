package bento

import (
	"fmt"
	"image"
	"log"
	"math"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
)

type layout struct {
	X, Y                        int `xml:",attr"`
	ContentWidth, ContentHeight int `xml:",attr"`
	InnerWidth, InnerHeight     int `xml:",attr"`
	OuterWidth, OuterHeight     int `xml:",attr"`
}

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

func (n *Box) maxContentWidth() int {
	if n.style.MaxWidth > 0 {
		return n.style.MaxWidth
	}
	return n.ContentWidth
}

func (n *Box) maxContentHeight() int {
	if n.style.MaxHeight > 0 {
		return n.style.MaxHeight
	}
	return n.ContentHeight
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
	if n.Tag == "button" || n.Tag == "text" {
		bounds := text.BoundString(n.style.Font, n.Content)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = n.style.Font.Metrics().Height.Ceil()
	} else if n.Tag == "p" {
		bounds := text.BoundParagraph(n.style.Font, n.Content, n.style.MaxWidth)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = max(bounds.Dy(), n.style.Font.Metrics().Height.Ceil())
	} else if n.Tag == "img" && n.style.Image != nil {
		bounds := n.style.Image.Bounds()
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = bounds.Dy()
	} else if n.Tag == "input" {
		txt := n.Content
		if txt == "" {
			txt = n.Attrs["placeholder"]
		}
		bounds := text.BoundString(n.style.Font, n.Content)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = n.style.Font.Metrics().Height.Ceil()
	} else if n.Tag == "textarea" {
		txt := n.Content
		if txt == "" {
			txt = n.Attrs["placeholder"]
		}
		bounds := text.BoundParagraph(n.style.Font, n.Content, n.style.MaxWidth)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = max(bounds.Dy(), n.style.Font.Metrics().Height.Ceil())
	} else if n.Tag != "canvas" && n.Tag != "row" && n.Tag != "col" {
		log.Fatalf("can't size %s", n.Tag)
	}
	for _, c := range n.Children {
		c.size()
		switch n.Tag {
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
	if n.style != nil && n.style.HGrow > 0 && n.Parent == nil {
		n.fillWidth(w)
	}
	if n.style != nil && n.style.VGrow > 0 && n.Parent == nil {
		n.fillHeight(h)
	}
	hgrow, vgrow := 0, 0
	for _, c := range n.Children {
		if !c.style.display() {
			continue
		}
		hg, vg := c.style.growth()
		hgrow += hg
		vgrow += vg
	}
	hspace, vspace := n.space()
	for _, c := range n.Children {
		if !c.style.display() {
			continue
		}
		hg, vg := c.style.growth()
		if hg > 0 {
			if n.Tag == "row" {
				halloc := int(math.Floor(float64(hg) / float64(hgrow) * float64(hspace)))
				c.fillWidth(c.OuterWidth + halloc)
			} else {
				c.fillWidth(n.ContentWidth)
			}
		}
		if vg > 0 && vgrow > 0 {
			if n.Tag == "col" {
				valloc := int(math.Floor(float64(vg) / float64(vgrow) * float64(vspace)))
				c.fillHeight(c.OuterHeight + valloc)
			} else {
				c.fillHeight(n.ContentHeight)
			}
		}
	}
	for _, c := range n.Children {
		if !c.style.display() {
			continue
		}
		c.grow()
	}
}

func (n *Box) fillWidth(w int) {
	if w < n.OuterWidth {
		return
	}
	_, mr, _, ml := n.style.margin()
	_, pr, _, pl := n.style.padding()
	n.OuterWidth = w
	n.InnerWidth = n.OuterWidth - ml - mr
	n.ContentWidth = n.InnerWidth - pr - pl
}

func (n *Box) fillHeight(h int) {
	if h < n.OuterHeight {
		return
	}
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
	for _, c := range n.Children {
		maxChildWidth = max(maxChildWidth, c.OuterWidth)
		maxChildHeight = max(maxChildHeight, c.OuterHeight)
		totalChildWidths += c.OuterWidth
		totalChildHeights += c.OuterHeight
	}

	var hspace, vspace int
	if n.Tag == "row" {
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
	for _, c := range n.Children {
		c.X = r.Min.X
		c.Y = r.Min.Y
		var ox, oy int
		if c.style != nil {
			ox, oy = c.style.OffsetX, c.style.OffsetY
		}
		if ox < 0 {
			ox = n.InnerWidth - c.OuterWidth + ox + 1
		}
		if oy < 0 {
			oy = n.InnerHeight - c.OuterHeight + oy + 1
		}
		c.X += ox
		c.Y += oy
	}
	extents := make([][2]int, len(n.Children))
	for i, c := range n.Children {
		if n.Tag == "row" {
			extents[i][0] = c.OuterWidth
			extents[i][1] = c.OuterHeight
		} else {
			extents[i][0] = c.OuterHeight
			extents[i][1] = c.OuterWidth
		}
	}
	var offsets [][2]int
	if n.Tag == "row" {
		offsets = distribute(hspace, n.InnerHeight, hj, vj, extents)
	} else {
		offsets = distribute(vspace, n.InnerWidth, vj, hj, extents)
	}
	for i, c := range n.Children {
		if n.Tag == "row" {
			c.X += offsets[i][0]
			c.Y += offsets[i][1]
		} else {
			c.Y += offsets[i][0]
			c.X += offsets[i][1]
		}
	}
	for _, c := range n.Children {
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

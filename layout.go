package bento

import (
	"fmt"
	"image"
	"log"
	"math"
	"sort"

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

func (n *Box) relayout() {
	n.size()
	n.grow()
	n.justify()
	n.sort()
}

func (n *Box) outerRect() image.Rectangle {
	return image.Rect(n.X, n.Y, n.X+n.OuterWidth, n.Y+n.OuterHeight)
}

func (n *Box) innerRect() image.Rectangle {
	mt, ml := n.style.Margin.Top, n.style.Margin.Left
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

// determine the minimum size of the box
// the box must expand to fit all of its children
func (n *Box) size() {
	n.ContentWidth = 0
	n.ContentHeight = 0
	if !n.style.Display {
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
		n.ContentWidth = int(float64(bounds.Dx()) * n.style.ScaleX)
		n.ContentHeight = int(float64(bounds.Dy()) * n.style.ScaleY)
	} else if n.Tag == "input" {
		txt := n.Content
		if txt == "" {
			txt = n.Attrs["placeholder"]
		}
		bounds := text.BoundString(n.style.Font, txt)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = n.style.Font.Metrics().Height.Ceil()
	} else if n.Tag == "textarea" {
		txt := n.Content
		if txt == "" {
			txt = n.Attrs["placeholder"]
		}
		bounds := text.BoundParagraph(n.style.Font, txt, n.style.MaxWidth)
		n.ContentWidth = bounds.Dx()
		n.ContentHeight = max(bounds.Dy(), n.style.Font.Metrics().Height.Ceil())
	} else if n.Tag != "canvas" && n.Tag != "row" && n.Tag != "col" {
		log.Fatalf("can't size %s", n.Tag)
	}
	for _, c := range n.Children {
		c.size()
		if c.style.Float {
			continue
		}
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
	n.InnerWidth = n.ContentWidth + n.style.Padding.Left + n.style.Padding.Right
	n.InnerHeight = n.ContentHeight + n.style.Padding.Top + n.style.Padding.Bottom
	n.OuterWidth = n.InnerWidth + n.style.Margin.Left + n.style.Margin.Right
	n.OuterHeight = n.InnerHeight + n.style.Margin.Top + n.style.Margin.Bottom
}

// grow children of the box to fit the space available, using their "grow" attribute
func (n *Box) grow() {
	// special case - the root of the tree can grow in either direction to fit the full window
	if n.Parent == nil {
		w, h := ebiten.WindowSize()
		if n.style.HGrow > 0 {
			n.fillWidth(w)
		}
		if n.style.VGrow > 0 {
			n.fillHeight(h)
		}
	}
	// sum the growth attributes to determine the relative growth of each child later
	hgrow, vgrow := 0, 0
	for _, c := range n.Children {
		if !c.style.Display || c.style.Float {
			continue
		}
		hgrow += c.style.HGrow
		vgrow += c.style.VGrow
	}
	// allocate leftover space to each child according to their relative growth terms
	hspace, vspace := n.space()
	for _, c := range n.Children {
		if !c.style.Display || c.style.Float {
			continue
		}
		if c.style.HGrow > 0 {
			if n.Tag == "row" {
				halloc := int(math.Floor(float64(c.style.HGrow) / float64(hgrow) * float64(hspace)))
				c.fillWidth(c.OuterWidth + halloc)
			} else {
				c.fillWidth(n.ContentWidth)
			}
		}
		if c.style.VGrow > 0 && vgrow > 0 {
			if n.Tag == "col" {
				valloc := int(math.Floor(float64(c.style.VGrow) / float64(vgrow) * float64(vspace)))
				c.fillHeight(c.OuterHeight + valloc)
			} else {
				c.fillHeight(n.ContentHeight)
			}
		}
	}
	// finally, allow the children to grow their own children using this newly allocated space
	for _, c := range n.Children {
		if !c.style.Display {
			continue
		}
		c.grow()
	}
}

func (n *Box) fillWidth(w int) {
	if w < n.OuterWidth {
		return
	}
	mr, ml := n.style.Margin.Right, n.style.Margin.Left
	pr, pl := n.style.Padding.Right, n.style.Padding.Left
	n.OuterWidth = w
	n.InnerWidth = n.OuterWidth - ml - mr
	n.ContentWidth = n.InnerWidth - pr - pl
}

func (n *Box) fillHeight(h int) {
	if h < n.OuterHeight {
		return
	}
	mt, mb := n.style.Margin.Top, n.style.Margin.Bottom
	pt, pb := n.style.Padding.Top, n.style.Padding.Bottom
	n.OuterHeight = h
	n.InnerHeight = n.OuterHeight - mt - mb
	n.ContentHeight = n.InnerHeight - pt - pb
}

// returns the horizontal and vertical space leftover in the box
func (n *Box) space() (int, int) {
	maxChildWidth := 0
	maxChildHeight := 0
	totalChildWidths := 0
	totalChildHeights := 0
	for _, c := range n.Children {
		if !c.style.Display || c.style.Float {
			continue
		}
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

// place the children in the box according to their justification
func (n *Box) justify() {
	r := n.innerRect()
	hspace, vspace := n.space()
	for _, c := range n.Children {
		if !c.style.Display || c.style.Float {
			continue
		}
		c.X = r.Min.X
		c.Y = r.Min.Y
		var ox, oy int
		ox, oy = c.style.OffsetX, c.style.OffsetY
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
		if !c.style.Display || c.style.Float {
			continue
		}
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
		offsets = distribute(hspace, n.InnerHeight, n.style.HJust, n.style.VJust, extents)
	} else {
		offsets = distribute(vspace, n.InnerWidth, n.style.VJust, n.style.HJust, extents)
	}
	for i, c := range n.Children {
		if c.style.Float {
			switch c.style.HJustSelf {
			case Start:
				c.X = r.Min.X
			case End:
				c.X = r.Min.X + n.InnerWidth - c.OuterWidth
			case Center:
				c.X = r.Min.X + (n.InnerWidth / 2) - (c.OuterWidth / 2)
			}
			switch c.style.VJustSelf {
			case Start:
				c.Y = r.Min.Y
			case End:
				c.Y = r.Min.Y + n.InnerHeight - c.OuterHeight
			case Center:
				c.Y = r.Min.Y + (n.InnerHeight / 2) - (c.OuterHeight / 2)
			}
			continue
		}
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
			panic(fmt.Errorf("can't handle main axis justification %q", mainj))
		}
		switch crossj {
		case Start:
			offsets[i][1] = 0
		case End:
			offsets[i][1] = crossspace - extents[i][1]
		case Center, Evenly, Around, Between:
			offsets[i][1] = int(math.Floor(float64(crossspace)/2)) - int(math.Floor(float64(extents[i][1])/2))
		default:
			panic(fmt.Errorf("can't handle cross axis justification %q", crossj))
		}
	}
	return offsets
}

// sort children by zIndex
func (n *Box) sort() {
	for i, c := range n.Children {
		if c.style.ZIndex == 0 {
			c.style.ZIndex = i + 1
		}
		c.sort()
	}
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].style.ZIndex < n.Children[j].style.ZIndex
	})
}

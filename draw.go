package bento

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (n *Box) Draw(img *ebiten.Image) {
	if !n.style.Display || n.style.Hidden {
		return
	}
	mt, ml := n.style.Margin.Top, n.style.Margin.Left
	pt, pl := n.style.Padding.Top, n.style.Padding.Left

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(float64(n.X), float64(n.Y))

	if debug {
		// Outer
		drawBox(img, n.OuterWidth, n.OuterHeight, color.White, true, op)
	}

	op.GeoM.Translate(float64(ml), float64(mt))
	if debug {
		// Inner
		drawBox(img, n.InnerWidth, n.InnerHeight, &color.RGBA{R: 200, G: 200, B: 200, A: 255}, true, op)
	}

	if n.style.Border != nil {
		n.style.Border.Draw(img, 0, 0, n.InnerWidth, n.InnerHeight, op)
	}

	switch n.Tag {
	case "button":
		n.style.Button[n.state].Draw(img, 0, 0, n.InnerWidth, n.InnerHeight, op)
	case "input", "textarea":
		state := n.state
		if n.editable.focus {
			state = active
		}
		n.style.Input[state].Draw(img, 0, 0, n.InnerWidth, n.InnerHeight, op)
	}

	op.GeoM.Translate(float64(pl), float64(pt))
	if debug {
		// Content
		drawBox(img, n.ContentWidth, n.ContentHeight, &color.RGBA{R: 100, G: 100, B: 100, A: 255}, true, op)
	}

	switch n.Tag {
	case "button", "text":
		text.DrawString(img, n.Content, n.style.Font, n.style.Color, n.style.Underline, n.ContentWidth, n.ContentHeight, text.Center, text.Center, -1, *op)
	case "p":
		if n.style.Scrollbar != nil {
			n.scrollable.position = text.DrawParagraph(
				img, n.Content, n.style.Font, n.style.Color, n.style.Underline,
				n.maxContentWidth(), n.maxContentHeight(),
				-1, n.scrollable.line,
				*op)
			if n.scrollable.position >= 0 {
				op.GeoM.Translate(float64(pl), -float64(pt))
				n.drawScrollbar(img, op)
			}
		} else {
			text.DrawParagraph(
				img, n.Content, n.style.Font, n.style.Color, n.style.Underline,
				n.maxContentWidth(), n.maxContentHeight(),
				-1, -1,
				*op)
		}
	case "img":
		imgOp := new(ebiten.DrawImageOptions)
		imgOp.GeoM.Scale(n.style.ScaleX, n.style.ScaleY)
		imgOp.GeoM.Translate(float64(n.X), float64(n.Y))
		imgOp.GeoM.Translate(float64(ml), float64(mt))
		imgOp.GeoM.Translate(float64(pl), float64(pt))
		img.DrawImage(n.style.Image, imgOp)
	case "input", "textarea":
		txt := n.Attrs["value"]
		if txt == "" {
			txt = n.Attrs["placeholder"]
		}
		if n.Tag == "input" {
			text.DrawString(img, txt, n.style.Font, n.style.Color, n.style.Underline, n.ContentWidth, n.ContentHeight, text.Start, text.Center, n.editable.Cursor(), *op)
		} else {
			n.scrollable.position = text.DrawParagraph(
				img, txt, n.style.Font, n.style.Color, n.style.Underline,
				n.maxContentWidth(), n.maxContentHeight(),
				n.editable.Cursor(), n.scrollable.line,
				*op)
			if n.scrollable.position >= 0 {
				op.GeoM.Translate(float64(pl), -float64(pt))
				n.drawScrollbar(img, op)
			}
		}
	case "canvas":
		n.call("onDraw", img)
	case "row", "col":
	default:
		log.Fatalf("can't draw %s", n.Tag)
	}

	if debug {
		op := new(ebiten.DrawImageOptions)
		op.GeoM.Translate(float64(n.X), float64(n.Y))
		// Annotate
		font := text.Font("RobotoMono", 18)
		txt := fmt.Sprintf("%s %dx%d", n.Tag, n.OuterWidth, n.OuterHeight)
		bounds := text.BoundString(font, txt)
		drawBox(img, bounds.Dx()+8, bounds.Dy()+4, color.White, false, op)
		op.GeoM.Translate(4, 4)
		text.DrawString(img, txt, font, color.Black, false,
			n.OuterWidth, n.OuterHeight, text.Start, text.Start, -1, *op)
	}

	for _, c := range n.Children {
		c.Draw(img)
	}

	if debug && n.Parent == nil {
		op := new(ebiten.DrawImageOptions)
		op.GeoM.Translate(float64(img.Bounds().Dx()-48), 24)
		txt := fmt.Sprintf("%.0f", ebiten.CurrentFPS())
		font := text.Font("RobotoMono", 24)
		bounds := text.BoundString(font, txt)
		drawBox(img, bounds.Dx(), bounds.Dy(), color.White, false, op)
		text.DrawString(
			img,
			fmt.Sprintf("%.0f", ebiten.CurrentFPS()),
			font, color.Black, false,
			0, 0, text.Start, text.Start, -1, *op)
	}
}

func drawBox(img *ebiten.Image, width, height int, c color.Color, border bool, op *ebiten.DrawImageOptions) {
	x1, y1 := op.GeoM.Apply(float64(0), float64(0))
	x2, y2 := op.GeoM.Apply(float64(width), float64(height))
	ebitenutil.DrawRect(img, x1, y1, x2-x1, y2-y1, c)
	if border {
		ebitenutil.DrawLine(img, x1, y1, x2, y1, color.Black)
		ebitenutil.DrawLine(img, x1, y1, x1, y2, color.Black)
		ebitenutil.DrawLine(img, x2, y2, x2, y1, color.Black)
		ebitenutil.DrawLine(img, x2, y2, x1, y2, color.Black)
	}
}

func (n *Box) drawScrollbar(img *ebiten.Image, op *ebiten.DrawImageOptions) {
	for i, r := range n.scrollRects(n.scrollable.position) {
		btn := n.style.Scrollbar[n.scrollable.state[i]][i]
		btn.Draw(img, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), op)
	}
}

func (n *Box) scrollRects(scrollPos float64) [4]image.Rectangle {
	var rects [4]image.Rectangle
	s := n.style.Scrollbar[0][0].Width()
	sf := float64(s)
	trackHeight := float64(n.InnerHeight) - 2.5*sf
	rects[0] = image.Rect(n.ContentWidth-s, 0, n.ContentWidth, s)               // top button
	rects[1] = image.Rect(n.ContentWidth-s, s, n.ContentWidth, n.InnerHeight-s) // track
	rects[2] = image.Rect(
		n.ContentWidth-s,
		int(trackHeight*scrollPos+0.75*sf),
		n.ContentWidth,
		int(trackHeight*scrollPos+1.75*sf)) // handle
	rects[3] = image.Rect(n.ContentWidth-s, n.InnerHeight-s, n.ContentWidth, n.InnerHeight) // bottom button
	return rects
}

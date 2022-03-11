package bento

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (n *Box) Draw(img *ebiten.Image) {
	if n.style.hidden() || !n.style.display() {
		return
	}
	mt, _, _, ml := n.style.margin()
	pt, _, _, pl := n.style.padding()

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(float64(n.X), float64(n.Y))

	if n.debug {
		// Outer
		drawBox(img, n.OuterWidth, n.OuterHeight, color.White, true, op)
	}

	op.GeoM.Translate(float64(ml), float64(mt))
	if n.debug {
		// Inner
		drawBox(img, n.InnerWidth, n.InnerHeight, &color.RGBA{R: 200, G: 200, B: 200, A: 255}, true, op)
	}

	if n.style != nil && n.style.Border != nil {
		n.style.Border.Draw(img, 0, 0, n.InnerWidth, n.InnerHeight, op)
	}

	switch n.tag {
	case "button":
		n.style.Button[int(n.buttonState)].Draw(img, 0, 0, n.InnerWidth, n.InnerHeight, op)
	case "input", "textarea":
		n.style.Input[int(n.inputState)].Draw(img, 0, 0, n.InnerWidth, n.InnerHeight, op)
	}

	op.GeoM.Translate(float64(pl), float64(pt))
	if n.debug {
		// Content
		drawBox(img, n.ContentWidth, n.ContentHeight, &color.RGBA{R: 100, G: 100, B: 100, A: 255}, true, op)
	}

	switch n.tag {
	case "button":
		text.DrawString(img, n.templateContent(), n.style.Font, n.style.Color, n.ContentWidth, n.ContentHeight, text.Center, text.Center, -1, *op)
	case "text":
		text.DrawString(img, n.templateContent(), n.style.Font, n.style.Color, n.ContentWidth, n.ContentHeight, text.Center, text.Center, -1, *op)
	case "p":
		txt := n.templateContent()
		maxHeight := n.style.MaxHeight
		if n.style.MaxHeight != 0 {
			maxHeight = max(maxHeight, n.ContentHeight)
		}
		if text.DrawParagraph(img, txt, n.style.Font, n.style.Color, n.style.MaxWidth, maxHeight, -1, n.scrollPosition, *op) {
			op.GeoM.Translate(float64(pl), -float64(pt))
			n.drawScrollbar(img, op)
		}
	case "img":
		img.DrawImage(n.style.Image, op)
	case "input", "textarea":
		txt := n.attrs["value"]
		if txt == "" && n.inputState != Active {
			txt = n.attrs["placeholder"]
		}
		cursor := -1
		if n.inputState == Active {
			t := time.Now().UnixMilli()
			if t-n.cursorTime <= 1000 {
				cursor = n.cursor
			} else if t-n.cursorTime >= 2000 {
				n.cursorTime = t
			}
		}
		if n.tag == "input" {
			text.DrawString(img, txt, n.style.Font, n.style.Color, n.ContentWidth, n.ContentHeight, text.Start, text.Center, cursor, *op)
		} else if text.DrawParagraph(img, txt, n.style.Font, n.style.Color, n.style.MaxWidth, n.style.MaxHeight, cursor, n.scrollPosition, *op) {
			op.GeoM.Translate(float64(pl), -float64(pt))
			n.drawScrollbar(img, op)
		}
	case "row", "col":
	default:
		log.Fatalf("can't draw %s", n.tag)
	}

	if n.debug {
		op := new(ebiten.DrawImageOptions)
		op.GeoM.Translate(float64(n.X), float64(n.Y))
		// Annotate
		text.DrawString(
			img,
			fmt.Sprintf("%s %dx%d", n.tag, n.OuterWidth, n.OuterHeight),
			text.Font("mono", 10), color.Black, n.OuterWidth, n.OuterHeight, text.Start, text.Start, -1, *op)
	}

	for _, c := range n.children {
		c.Draw(img)
	}

	if n.debug && n.parent == nil {
		op := new(ebiten.DrawImageOptions)
		op.GeoM.Translate(float64(img.Bounds().Dx()-48), 24)
		text.DrawString(
			img,
			fmt.Sprintf("%.0f", ebiten.CurrentFPS()),
			text.Font("mono", 24), color.White, 0, 0, text.Start, text.Start, -1, *op)
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
	for i, r := range n.scrollRects() {
		btn := n.style.Scrollbar[int(n.scrollState[i])][i]
		btn.Draw(img, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), op)
	}
}

func (n *Box) scrollRects() [4]image.Rectangle {
	var rects [4]image.Rectangle
	s := n.style.Scrollbar[0][0].Width()
	rects[0] = image.Rect(n.ContentWidth-s, 0, n.ContentWidth, s)                           // top button
	rects[1] = image.Rect(n.ContentWidth-s, s, n.ContentWidth, n.InnerHeight-s)             // track
	rects[2] = image.Rect(0, 0, 0, 0)                                                       // handle
	rects[3] = image.Rect(n.ContentWidth-s, n.InnerHeight-s, n.ContentWidth, n.InnerHeight) // bottom button
	return rects
}

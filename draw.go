package bento

import (
	"image"
	"image/color"

	"github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (n *node) Draw(img *ebiten.Image) {
	if n.parent == nil {
		n.size()
		n.grow()
		n.justify()
	}
	if n.style.hidden() || !n.style.display() {
		return
	}
	content := n.contentRect()
	inner := n.innerRect()
	outer := n.outerRect()
	if n.debug {
		// Outer
		drawBox(img, outer, color.White, true)
		// Inner
		drawBox(img, inner, &color.RGBA{R: 200, G: 200, B: 200, A: 255}, true)
	}
	if n.style != nil && n.style.Border != nil {
		n.style.Border.Draw(img, inner.Min.X, inner.Min.Y, inner.Dx(), inner.Dy())
	}
	if n.style.Button != nil {
		n.style.Button.Draw(img)
	}
	if n.debug {
		// Content
		drawBox(img, content, &color.RGBA{R: 100, G: 100, B: 100, A: 255}, true)
	}
	switch n.tag {
	case "button":
		text.DrawString(img, n.templateContent(), n.style.Font, n.style.Color, content, text.Center, text.Center)
	case "text":
		text.DrawString(img, n.templateContent(), n.style.Font, n.style.Color, content, text.Center, text.Center)
	case "p":
		if scroll := n.style.Scrollbar; scroll != nil {
			n.buffer = ebiten.NewImage(n.TextBounds.Dx(), n.TextBounds.Dy())
			text.DrawParagraph(n.buffer, n.templateContent(), n.style.Font, n.style.Color, 0, 0, n.style.MaxWidth, -n.TextBounds.Min.Y)
			offset := int(scroll.Position() * float64(n.TextBounds.Dy()))
			cropped := ebiten.NewImageFromImage(n.buffer.SubImage(image.Rect(0, offset, content.Dx(), content.Dy()+offset)))
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
			img.DrawImage(cropped, op)
			n.style.Scrollbar.Draw(img)
		} else {
			text.DrawParagraph(img, n.templateContent(), n.style.Font, n.style.Color, content.Min.X, content.Min.Y, n.style.MaxWidth, -n.TextBounds.Min.Y)
		}
	case "img":
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
		img.DrawImage(n.style.Background, op)
	}
	if n.debug {
		// Annotate
		text.DrawString(img, n.tag, text.Font("mono", 16), color.Black, outer.Add(image.Pt(4, 4)), text.Start, text.Start)
	}
	for _, c := range n.children {
		c.Draw(img)
	}
}

func drawBox(img *ebiten.Image, rect image.Rectangle, c color.Color, border bool) {
	x, y := float64(rect.Min.X), float64(rect.Min.Y)
	w, h := float64(rect.Dx()), float64(rect.Dy())
	ebitenutil.DrawRect(img, x, y, w, h, c)
	if border {
		ebitenutil.DrawLine(img, x, y, x+w, y, color.Black)
		ebitenutil.DrawLine(img, x, y, x, y+h, color.Black)
		ebitenutil.DrawLine(img, x+w, y+h, x+w, y, color.Black)
		ebitenutil.DrawLine(img, x+w, y+h, x, y+h, color.Black)
	}
}

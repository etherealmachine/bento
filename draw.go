package bento

import (
	"image"
	"image/color"

	"github.com/etherealmachine/bento/v1/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (n *Node) Draw(img *ebiten.Image) {
	if n.Parent == nil {
		n.updateSize()
		n.layout()
	}
	if n.Hidden || !n.Display {
		return
	}
	content := n.ContentRect()
	inner := n.InnerRect()
	outer := n.OuterRect()
	if n.Debug {
		// Outer
		drawBox(img, outer, color.White, true)
		// Inner
		drawBox(img, inner, &color.RGBA{R: 200, G: 200, B: 200, A: 255}, true)
	}
	if n.Style != nil && n.Style.Border != nil {
		n.Style.Border.Draw(img, inner.Min.X, inner.Min.Y, inner.Dx(), inner.Dy())
	}
	state := Idle
	if n.Hover {
		state = Hover
	}
	if n.Active {
		state = Active
	}
	if n.Disabled {
		state = Disabled
	}
	if n.XMLName.Local == "button" && n.Style != nil && n.Style.Button != nil {
		n.Style.Button.Draw(img, inner.Min.X, inner.Min.Y, inner.Dx(), inner.Dy(), state)
	}
	if n.Debug {
		// Content
		drawBox(img, content, &color.RGBA{R: 100, G: 100, B: 100, A: 255}, true)
	}
	switch n.XMLName.Local {
	case "button":
		text.DrawString(img, n.Content(), n.Style.Font, n.Style.Color, content, text.Center, text.Center)
	case "text":
		text.DrawString(img, n.Content(), n.Style.Font, n.Style.Color, content, text.Center, text.Center)
	case "p":
		if n.TextBounds.Dy() > content.Dy() {
			n.buffer = ebiten.NewImage(n.TextBounds.Dx(), n.TextBounds.Dy())
			text.DrawParagraph(n.buffer, n.Content(), n.Style.Font, n.Style.Color, 0, 0, n.Style.MaxWidth, -n.TextBounds.Min.Y)
			cropped := ebiten.NewImageFromImage(n.buffer.SubImage(image.Rect(0, n.Scroll, content.Dx(), content.Dy()+n.Scroll)))
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
			img.DrawImage(cropped, op)
			n.Style.Scrollbar.Draw(img, inner.Max.X, inner.Min.Y, inner.Dy(), float64(n.Scroll)/float64(n.TextBounds.Dy()-n.ContentHeight), state)
		} else {
			text.DrawParagraph(img, n.Content(), n.Style.Font, n.Style.Color, content.Min.X, content.Min.Y, n.Style.MaxWidth, -n.TextBounds.Min.Y)
		}
	case "img":
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
		img.DrawImage(n.Style.Background, op)
	}
	if n.Debug {
		// Annotate
		r := outer.Add(image.Pt(4, 4))
		text.DrawString(img, n.XMLName.Local, text.Font("mono", 16), color.Black, &r, text.Start, text.Start)
	}
	for _, c := range n.Children {
		c.Draw(img)
	}
}

func drawBox(img *ebiten.Image, rect *image.Rectangle, c color.Color, border bool) {
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

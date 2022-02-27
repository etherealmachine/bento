package bento

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/*
	Branched from https://github.com/blizzy78/ebitenui/blob/ca1a302d930b/image/nineslice.go
	Copyright 2020 Maik Schreiber

	Permission is hereby granted, free of charge, to any person obtaining a copy of
	this software and associated documentation files (the "Software"), to deal in
	the Software without restriction, including without limitation the rights to
	use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
	of the Software, and to permit persons to whom the Software is furnished to do
	so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
	FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
	COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
	IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
	CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// A NineSlice is an image that can be drawn with any width and height. It is basically a 3x3 grid of image tiles:
// The corner tiles are drawn as-is, while the center columns and rows of tiles will be stretched to fit the desired
// width and height.
type NineSlice struct {
	image   *ebiten.Image
	width   int
	widths  [3]int
	heights [3]int

	frames [][9]*ebiten.Image
}

func NewNineSlice(filename string, width int, widths, heights [3]int, frames int, horizontal bool) *NineSlice {
	n := &NineSlice{
		width:   width,
		widths:  widths,
		heights: heights,
	}
	n.image, _, _ = ebitenutil.NewImageFromFile(filename)
	n.frames = make([][9]*ebiten.Image, frames)
	for i := 0; i < frames; i++ {
		n.frames[i] = n.createTiles(i, horizontal)
	}
	n.image = nil
	return n
}

func (n *NineSlice) draw(screen *ebiten.Image, frame int, x, y, width, height int) {
	n.drawTiles(screen, n.frames[frame], x, y, width, height)
}

func (n *NineSlice) drawScrollbar(screen *ebiten.Image, x, y, height int) {
	h := n.heights[0] + n.heights[1] + n.heights[2]
	n.drawTiles(screen, n.frames[0], x, y, n.width, h)
	n.drawTiles(screen, n.frames[1], x, y+h, n.width, height-2*h)
	n.drawTiles(screen, n.frames[2], x, y+h, n.width, h)
	n.drawTiles(screen, n.frames[5], x, y+height-h, n.width, h)
}

func (n *NineSlice) drawTiles(screen *ebiten.Image, tiles [9]*ebiten.Image, x, y, width, height int) {
	sy := 0
	ty := y
	for r, sh := range n.heights {
		sx := 0
		tx := x

		var th int
		if r == 1 {
			th = height - n.heights[0] - n.heights[2]
		} else {
			th = sh
		}

		for c, sw := range n.widths {
			var tw int
			if c == 1 {
				tw = width - n.widths[0] - n.widths[2]
			} else {
				tw = sw
			}

			n.drawTile(screen, tiles[r*3+c], tx, ty, sw, sh, tw, th)

			sx += sw
			tx += tw
		}

		sy += sh
		ty += th
	}
}

func (n *NineSlice) drawTile(screen *ebiten.Image, tile *ebiten.Image, tx int, ty int, sw int, sh int, tw int, th int) {
	if sw <= 0 || sh <= 0 || tw <= 0 || th <= 0 {
		return
	}

	opts := ebiten.DrawImageOptions{
		Filter: ebiten.FilterNearest,
	}

	if tw != sw || th != sh {
		opts.GeoM.Scale(float64(tw)/float64(sw), float64(th)/float64(sh))
	}

	opts.GeoM.Translate(float64(tx), float64(ty))

	screen.DrawImage(tile, &opts)
}

func (n *NineSlice) createTiles(frame int, horizontal bool) [9]*ebiten.Image {
	var tiles [9]*ebiten.Image
	min := n.image.Bounds().Min

	height := n.heights[0] + n.heights[1] + n.heights[2]
	sy := min.Y
	for r, sh := range n.heights {
		sx := min.X
		for c, sw := range n.widths {
			if sh > 0 && sw > 0 {
				rect := image.Rect(0, 0, sw, sh)
				if horizontal {
					rect = rect.Add(image.Point{Y: frame * height})
				} else {
					rect = rect.Add(image.Point{X: frame * n.width})
				}
				rect = rect.Add(image.Point{sx, sy})
				tiles[r*3+c] = n.image.SubImage(rect).(*ebiten.Image)
			}
			sx += sw
		}
		sy += sh
	}
	return tiles
}

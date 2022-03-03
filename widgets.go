package bento

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ButtonState int

const (
	Idle     = ButtonState(0)
	Hover    = ButtonState(1)
	Active   = ButtonState(2)
	Disabled = ButtonState(3)
)

type Button struct {
	states  [4]*NineSlice
	rect    image.Rectangle
	state   ButtonState
	onClick func(id string)
	box     *Box
}

func NewButton(img *ebiten.Image, widths, heights [3]int, onClick func(id string)) *Button {
	b := &Button{
		onClick: onClick,
	}
	w := widths[0] + widths[1] + widths[2]
	for i := 0; i < 4; i++ {
		b.states[i] = NewNineSlice(img, widths, heights, w*i, 0)
	}
	return b
}

func (b *Button) Bounds() image.Rectangle {
	return b.rect
}

func (b *Button) Update(keys []ebiten.Key) ([]ebiten.Key, error) {
	b.state = buttonState(b.rect)
	if b.state == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.onClick(b.box.attrs["id"])
	}
	return keys, nil
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.states[int(b.state)].Draw(screen, b.rect.Min.X, b.rect.Min.Y, b.rect.Dx(), b.rect.Dy())
}

var (
	Scrollspeed = 0.1
)

type Scrollbar struct {
	states                                               [3][4]*NineSlice
	x, y, height                                         int
	position                                             float64
	topBtnState, bottomBtnState, trackState, handleState ButtonState
}

func NewScrollbar(img *ebiten.Image, widths, heights [3]int) *Scrollbar {
	s := &Scrollbar{}
	w := widths[0] + widths[1] + widths[2]
	h := heights[0] + heights[1] + heights[2]
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			s.states[i][j] = NewNineSlice(img, widths, heights, w*i, j*h)
		}
	}
	return s
}

func (s *Scrollbar) Position() float64 {
	return s.position
}

func (s *Scrollbar) Update(keys []ebiten.Key) ([]ebiten.Key, error) {
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	h := s.states[0][0].heights[0] + s.states[0][0].heights[1] + s.states[0][0].heights[2]

	s.topBtnState = buttonState(image.Rect(s.x, s.y, s.x+w, s.y+h))
	if s.topBtnState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.position -= Scrollspeed
	}
	trackHeight := s.height - 2*h
	s.trackState = buttonState(image.Rect(s.x, s.y+h, s.x+w, s.y+h+trackHeight))
	s.handleState = buttonState(image.Rect(s.x, s.y+h+int(s.position*float64(trackHeight-h)), s.x+w, s.y+h+int(s.position*float64(trackHeight-h))+h))
	s.bottomBtnState = buttonState(image.Rect(s.x, s.y+s.height-h, s.x+w, s.y+s.height))
	if s.bottomBtnState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.position += Scrollspeed
	}

	_, dy := ebiten.Wheel()
	s.position += dy * Scrollspeed

	// TODO: bad position can crash during Draw
	if s.position < 0 {
		s.position = 0
	}
	if s.position > 1 {
		s.position = 1
	}
	return keys, nil
}

func (s *Scrollbar) Draw(screen *ebiten.Image) {
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	h := s.states[0][0].heights[0] + s.states[0][0].heights[1] + s.states[0][0].heights[2]
	trackHeight := s.height - 2*h
	s.states[s.topBtnState][0].Draw(screen, s.x, s.y, w, h)
	s.states[s.trackState][1].Draw(screen, s.x, s.y+h, w, trackHeight)
	s.states[s.handleState][2].Draw(screen, s.x, s.y+h+int(s.position*float64(trackHeight-h)), w, h)
	s.states[s.bottomBtnState][3].Draw(screen, s.x, s.y+s.height-h, w, h)
}

func (s *Scrollbar) Bounds() image.Rectangle {
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	return image.Rect(s.x, s.y, s.x+w, s.y+s.height)
}

func buttonState(rect image.Rectangle) ButtonState {
	x, y := ebiten.CursorPosition()
	if inside(rect, x, y) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			return Active
		}
		return Hover
	}
	return Idle
}

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
	widths  [3]int
	heights [3]int
	tiles   [9]*ebiten.Image
}

func NewNineSlice(img *ebiten.Image, widths, heights [3]int, offsetX, offsetY int) *NineSlice {
	n := &NineSlice{
		widths:  widths,
		heights: heights,
	}
	n.createTiles(img, offsetX, offsetY)
	return n
}

func (n *NineSlice) createTiles(img *ebiten.Image, ox, oy int) {
	min := img.Bounds().Min
	sy := min.Y
	for r, sh := range n.heights {
		sx := min.X
		for c, sw := range n.widths {
			if sh > 0 && sw > 0 {
				rect := image.Rect(ox, oy, ox+sw, oy+sh).Add(image.Point{sx, sy})
				n.tiles[r*3+c] = img.SubImage(rect).(*ebiten.Image)
			}
			sx += sw
		}
		sy += sh
	}
}

func (n *NineSlice) Draw(screen *ebiten.Image, x, y, width, height int) {
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

			n.drawTile(screen, n.tiles[r*3+c], tx, ty, sw, sh, tw, th)

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

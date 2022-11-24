package bento

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Scrollable struct {
	state    [4]State
	line     int
	position float64
}

func (s *Scrollable) update(b *Box) error {
	if b.style.Scrollbar == nil || b.Attrs["disabled"] == "true" {
		return nil
	}
	mt, ml := b.style.Margin.Top, b.style.Margin.Left
	pt, pl := b.style.Padding.Top, b.style.Padding.Left
	rects := b.scrollRects(s.position)
	for i := 0; i < 4; i++ {
		if i == 2 {
			continue
		}
		// TODO: the math here works out but it's confusing
		r := rects[i].Add(image.Pt(b.X+ml+pl+pl, b.Y+mt+pt-pt))
		s.state[i] = idle
		x, y := ebiten.CursorPosition()
		if inside(r, x, y) {
			s.state[i] = hover
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				s.state[i] = active
			}
		}
		if s.state[i] == active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if i == 0 {
				s.line--
				if s.line < 0 {
					s.line = 0
				}
			} else if i == 3 && s.position < 1 {
				s.line++
			}
		}
	}
	return nil
}

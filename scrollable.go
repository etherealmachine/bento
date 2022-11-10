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

func (s *Scrollable) Update(b *Box) error {
	if b.style.Scrollbar == nil {
		return nil
	}
	mt, _, _, ml := b.style.margin()
	pt, _, _, pl := b.style.padding()
	rects := b.scrollRects(s.position)
	for i := 0; i < 4; i++ {
		if i == 2 {
			continue
		}
		// TODO: the math here works out but it's confusing
		r := rects[i].Add(image.Pt(b.X+ml+pl+pl, b.Y+mt+pt-pt))
		s.state[i] = getState(r)
		if s.state[i] == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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

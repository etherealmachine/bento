package main

import (
	"image/color"
	"time"

	"github.com/etherealmachine/bento"
)

type Page4 struct {
	t *time.Time
}

func (p *Page4) OnUpdate(e *bento.Event) bool {
	if p.t == nil {
		t := time.Now()
		p.t = &t
	}
	x := float64(time.Since(*p.t)) / float64(5*time.Second)
	if x >= 1 {
		x = 1
	}
	e.Box.Style.Color = color.RGBA{R: 255, G: 255, B: 255, A: uint8(255 * x)}
	return false
}

func (p *Page4) UI() string {
	return `<col grow="1" justify="center" margin="24px">
		<text font="NotoSans 24" margin="4px" padding="12px" onUpdate="OnUpdate">
			github.com/etherealmachine/bento
		</text>
	</col>`
}

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Page1 struct {
	Clicks int
}

func (p *Page1) OnKeyDown(key ebiten.Key) bool {
	return false
}

func (p *Page1) Click(_ string) {
	p.Clicks++
}

func (p *Page1) Reset(_ string) {
	p.Clicks = 0
}

func (p *Page1) UI() string {
	return `<col grow="1" justify="center" border="frame.png 10">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum</text>
		<text font="RobotoMono 24" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum</text>
		<row>
			<button onClick="Click" color="#ffffff" margin="4px" padding="12px" btn="button.png 6">Clicks: {{.Clicks}}</button>
			<button onClick="Reset" color="#ffffff" margin="4px" padding="12px" btn="button.png 6">Reset</button>
		</row>
	</col>`
}

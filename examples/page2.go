package main

import (
	"image/color"

	"github.com/etherealmachine/bento"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Page2 struct {
	Clicks     int
	Paragraphs []string
}

func (p *Page2) OnKeyDown(key ebiten.Key) bool {
	return false
}

func (p *Page2) handleClick(_ string) {
	p.Clicks++
}

func (p *Page2) Reset(_ string) {
	p.Clicks = 0
}

func (p *Page2) Text() *bento.Style {
	scrollbar, _, _ := ebitenutil.NewImageFromFile("scrollbar.png")
	return &bento.Style{
		Extends:   "p",
		FontName:  "NotoSans",
		FontSize:  16,
		Margin:    &bento.Spacing{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Padding:   &bento.Spacing{Top: 12, Bottom: 12, Left: 12, Right: 12},
		Color:     &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Scrollbar: bento.NewScrollbar(scrollbar, [3]int{6, 20, 6}, [3]int{6, 20, 6}),
		Attrs: map[string]string{
			"maxWidth":  "50em",
			"maxHeight": "8em",
		},
	}
}

func (p *Page2) Button() *bento.Style {
	button, _, _ := ebitenutil.NewImageFromFile("button.png")
	return &bento.Style{
		Extends:  "button",
		FontName: "NotoSans",
		FontSize: 16,
		Margin:   &bento.Spacing{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Padding:  &bento.Spacing{Top: 12, Bottom: 12, Left: 12, Right: 12},
		Color:    &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Button:   bento.NewButton(button, [3]int{6, 20, 6}, [3]int{6, 20, 6}, p.handleClick),
	}
}

func (p *Page2) UI() string {
	return `<col grow="1" justify="center" border="frame.png 10">
		<row grow="1">
			<img bg="profile.png"/>
			<p grow="1" font="NotoSans 16" color="#ffffff" maxWidth="40em">{{index .Paragraphs 0}}</p>
		</row>
		<text font="RobotoMono 24" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum</text>
		<Text border="frame.png 10">
	{{index .Paragraphs 1}}
	{{index .Paragraphs 2}}
		</Text>
		<row>
			<Button>Clicks: {{.Clicks}}</Button>
			<button onClick="Reset" color="#ffffff" margin="4px" padding="12px" btn="button.png 6">Reset</button>
		</row>
	</col>`
}

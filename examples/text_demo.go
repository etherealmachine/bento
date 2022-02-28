package main

import (
	"image/color"

	"github.com/etherealmachine/bento"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TextDemo struct {
	Clicks int
}

func (t *TextDemo) OnKeyDown(key ebiten.Key) bool {
	return false
}

func (t *TextDemo) handleClick(_ string) {
	t.Clicks++
}

func (t *TextDemo) Text() *bento.Style {
	scrollbar, _, _ := ebitenutil.NewImageFromFile("scrollbar.png")
	return &bento.Style{
		Extends:   "p",
		FontName:  "Gentium",
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

func (t *TextDemo) Button() *bento.Style {
	button, _, _ := ebitenutil.NewImageFromFile("button.png")
	return &bento.Style{
		Extends:  "button",
		FontName: "Gentium",
		FontSize: 16,
		Margin:   &bento.Spacing{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Padding:  &bento.Spacing{Top: 12, Bottom: 12, Left: 12, Right: 12},
		Color:    &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Button:   bento.NewButton(button, [3]int{6, 20, 6}, [3]int{6, 20, 6}, t.handleClick),
	}
}

func (t *TextDemo) UI() string {
	return `<col width="100%" height="100%" justify="center" border="frame.png 10 12 10 10 12 10">
		<row flex="1">
			<img bg="profile.png"/>
			<p flex="1" font="NotoSans 16" color="#ffffff" maxWidth="40em">
Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.
The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.
			</p>
		</row>
		<text font="RobotoMono 24" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum</text>
		<Text border="frame.png 10">
Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.
The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.
		</Text>
		<Button>Clicks: {{.Clicks}}</Button>
	</col>`
}

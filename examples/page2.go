package main

import (
	"image/color"
	"log"

	"github.com/etherealmachine/bento"
)

type Page2 struct {
	Clicks         int
	Paragraphs     []string
	Title, Content string
}

func (p *Page2) Click() {
	p.Clicks++
}

func (p *Page2) Change(old, new string) {
	p.Title = new
}

func (p *Page2) Reset() {
	p.Clicks = 0
}

func (p *Page2) Text() *bento.Style {
	scrollbar, err := bento.ParseScrollbar("scrollbar.png 6")
	if err != nil {
		log.Fatal(err)
	}
	return &bento.Style{
		Extends:   "p",
		FontName:  "NotoSans",
		FontSize:  16,
		Margin:    &bento.Spacing{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Padding:   &bento.Spacing{Top: 12, Bottom: 12, Left: 12, Right: 12},
		Color:     &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Scrollbar: scrollbar,
		Attrs: map[string]string{
			"maxWidth":  "50em",
			"maxHeight": "8lh",
		},
	}
}

func (p *Page2) Button() *bento.Style {
	btn, err := bento.ParseButton("button.png 6")
	if err != nil {
		log.Fatal(err)
	}
	return &bento.Style{
		Extends:  "button",
		FontName: "NotoSans",
		FontSize: 16,
		Margin:   &bento.Spacing{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Padding:  &bento.Spacing{Top: 12, Bottom: 12, Left: 12, Right: 12},
		Color:    &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Button:   btn,
	}
}

func (p *Page2) Input() *bento.Style {
	input, err := bento.ParseButton("input.png 6")
	if err != nil {
		log.Fatal(err)
	}
	return &bento.Style{
		Extends:  "input",
		FontName: "NotoSans",
		FontSize: 16,
		Margin:   &bento.Spacing{Top: 4, Bottom: 4, Left: 4, Right: 4},
		Padding:  &bento.Spacing{Top: 12, Bottom: 12, Left: 12, Right: 12},
		Color:    &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Input:    input,
	}
}

func (p *Page2) UI() string {
	return `<col grow="1" justify="center" border="frame.png 10">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">Page 2</text>
		<col grow="1">
			<text font="RobotoMono 24" color="#ffffff" margin="4px" padding="12px">{{.Title}}</text>
			<row>
				<img src="profile.png"/>
				<p grow="1" padding="24px" font="NotoSans 16" color="#ffffff" maxWidth="80em">{{.Content}}</p>
			</row>
			<Input onChange="Change" grow="1 0" placeholder="Title" color="#ffffff" />
			<textarea grow="1" value="Content" color="#ffffff" margin="4px" padding="1em" input="textarea.png 6" maxHeight="6em" scrollbar="scrollbar.png 6" />
		</col>
		<Text grow="0 1" border="frame.png 10">{{index .Paragraphs 1}}
{{index .Paragraphs 2}}</Text>
		<row grow="1 0">
			<Button onClick="Click">Clicks: {{.Clicks}}</Button>
			<button onClick="Reset" color="#ffffff" margin="4px" padding="12px" btn="button.png 6">Reset</button>
		</row>
	</col>`
}

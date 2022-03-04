package main

import (
	"log"

	"github.com/etherealmachine/bento"
	"github.com/hajimehoshi/ebiten/v2"
)

type Demo struct {
	CurrentPage int
}

func (d *Demo) OnKeyDown(key ebiten.Key) bool {
	return false
}

func (d *Demo) Prev(_ string) {
	d.CurrentPage = 0
}

func (d *Demo) Next(_ string) {
	d.CurrentPage = 1
}

func (d *Demo) Page1() *bento.Box {
	p := &Page1{}
	b, err := bento.Build(p)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (d *Demo) Page2() *bento.Box {
	p := &Page2{Title: "Loomings", Content: paragraphs[0], Paragraphs: paragraphs}
	b, err := bento.Build(p)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (d *Demo) UI() string {
	return `<col grow="1" border="frame.png 10">
		<Page1 grow="1" display="{{eq .CurrentPage 0}}"/>
		<Page2 grow="1" display="{{eq .CurrentPage 1}}"/>
		<row grow="1 0" justify="between">
			<button onClick="Prev" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{eq .CurrentPage 0}}">Prev</button>
			<button onClick="Next" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{eq .CurrentPage 1}}">Next</button>
		</row>
	</col>`
}

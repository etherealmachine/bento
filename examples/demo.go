package main

import (
	"github.com/etherealmachine/bento"
)

type Demo struct {
	CurrentPage int
}

func (d *Demo) Prev() {
	d.CurrentPage = 0
}

func (d *Demo) Next() {
	d.CurrentPage = 1
}

func (d *Demo) Page1() bento.Component {
	return &Page1{}
}

func (d *Demo) Page2() bento.Component {
	return &Page2{Title: "Loomings", Content: paragraphs[0], Paragraphs: paragraphs}
}

func (d *Demo) UI() string {
	return `<col grow="1" border="frame.png 10">
		{{ if eq .CurrentPage 0 }}
			<Page1 grow="1" />
		{{ end }}
		{{ if eq .CurrentPage 1 }}
			<Page2 grow="1" />
		{{ end }}
		<row grow="1 0" justify="between">
			<button onClick="Prev" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{eq .CurrentPage 0}}">Prev</button>
			<button onClick="Next" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{eq .CurrentPage 1}}">Next</button>
		</row>
	</col>`
}

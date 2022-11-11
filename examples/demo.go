package main

import (
	"github.com/etherealmachine/bento"
)

type Demo struct {
	CurrentPage int
	page1       *Page1
	page2       *Page2
}

func (d *Demo) Prev(_ *bento.Event) {
	d.CurrentPage = 0
}

func (d *Demo) Next(_ *bento.Event) {
	d.CurrentPage = 1
}

func (d *Demo) Page1() bento.Component {
	if d.page1 == nil {
		d.page1 = &Page1{Paragraphs: paragraphs}
	}
	return d.page1
}

func (d *Demo) Page2() bento.Component {
	if d.page2 == nil {
		d.page2 = &Page2{}
	}
	return d.page2
}

func (d *Demo) UI() string {
	return `<col grow="1">
		{{ if eq .CurrentPage 0 }}
			<Page1 />
		{{ end }}
		{{ if eq .CurrentPage 1 }}
			<Page2 />
		{{ end }}
		<row grow="1 0" justify="between">
			<button onClick="Prev" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{eq .CurrentPage 0}}">Prev</button>
			<button onClick="Next" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{eq .CurrentPage 1}}">Next</button>
		</row>
	</col>`
}

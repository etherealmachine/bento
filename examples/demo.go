package main

import (
	"github.com/etherealmachine/bento"
)

type Demo struct {
	CurrentPage int
	Pages       []bento.Component
}

func NewDemo() *Demo {
	return &Demo{
		Pages: []bento.Component{
			&Page1{Paragraphs: paragraphs[:3]},
			&Page2{},
			&Page3{},
			&Page4{},
		},
	}
}

func (d *Demo) Prev(_ *bento.Event) {
	d.CurrentPage--
	if d.CurrentPage < 0 {
		d.CurrentPage = 0
	}
}

func (d *Demo) Next(_ *bento.Event) {
	d.CurrentPage++
	if d.CurrentPage >= len(d.Pages) {
		d.CurrentPage = len(d.Pages) - 1
	}
}

func (d *Demo) PageN() bento.Component {
	return d.Pages[d.CurrentPage]
}

func (d *Demo) UI() string {
	return `<col grow="1">
		<PageN />
		<row grow="1 0" justify="between">
			<button onClick="Prev" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{ eq .CurrentPage 0 }}">Prev</button>
			<button onClick="Next" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" disabled="{{ ge .CurrentPage (len .Pages) }}">Next</button>
		</row>
	</col>`
}

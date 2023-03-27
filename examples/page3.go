package main

import (
	"github.com/etherealmachine/bento"
)

type Page3 struct {
	TextInput, TextArea string
}

func (p *Page3) UpdateTextInput(event *bento.Event) {
	p.TextInput = event.Value
}

func (p *Page3) UpdateTextArea(event *bento.Event) {
	p.TextArea = event.Value
}

func (p *Page3) UI() string {
	return `<col grow="1" justify="start" margin="24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px" zIndex="300">
			But wait, there's more!
		</text>
		<input
				minWidth="20em"
				onChange="UpdateTextInput"
				placeholder="Editable Inputs"
				input="input.png 6"
				color="#ffffff"
				margin="4px"
				padding="16px"
				value="{{ .TextInput }}" />
		<textarea
				minWidth="40em"
				minHeight="8lh"
				onChange="UpdateTextArea"
				placeholder="And editable multi-line text areas!"
				input="input.png 6"
				color="#ffffff"
				margin="4px"
				padding="16px"
				value="{{ .TextArea }}" />
	</col>`
}

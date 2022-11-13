package main

import (
	"github.com/etherealmachine/bento"
)

type Page3 struct {
	TextInput string
}

func (p *Page3) TextChange(event *bento.Event) {
	p.TextInput = event.Box.Content
}

func (p *Page3) UI() string {
	return `<col grow="1" justify="start" margin="24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">
			But wait, there's more!
		</text>
		<p
				font="RobotoMono 16"
				color="#ffffff"
				margin="4px"
				padding="16px 24px">
<![CDATA[
// Images with the <img> tag
<img src="profile.png" />
]]>	
		</p>
		<img src="profile.png" />
		<input
				minWidth="20em"
				onChange="TextChange"
				placeholder="Editable Inputs"
				input="input.png 6"
				color="#ffffff"
				margin="4px"
				padding="16px">
			{{ .TextInput }}
		</input>
		<textarea
				minWidth="40em"
				minHeight="8lh"
				onChange="TextChange"
				placeholder="And editable multi-line text areas!"
				input="input.png 6"
				color="#ffffff"
				margin="4px"
				padding="16px">
			{{ .TextInput }}
		</textarea>
	</col>`
}

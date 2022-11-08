package main

import "github.com/etherealmachine/bento"

type Page1 struct {
	Clicks     int
	Paragraphs []string
}

func (p *Page1) Click() {
	p.Clicks++
}

func (p *Page1) Debug(event *bento.Event) {
	event.Box.ToggleDebug()
}

func (p *Page1) UI() string {
	return `<col grow="1" justify="start" margin="24px 0 0 24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">
			What is Bento?
		</text>
		<col margin="0 0 0 16px">
			<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
				Bento is a GUI library for Ebiten that puts your UI in a box!
			</text>
			<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
				Bento has support for 9 components, including:
			</text>
			<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
				Text boxes, like this one
			</text>
			<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
				Clickable buttons, with hover, active, and disabled states, and events
			</text>
			<row justify="start center">
				<button onClick="Click" color="#ffffff" margin="12px" padding="12px" btn="button.png 6">Click Me</button>
				<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
					I've been clicked {{ .Clicks }} times
				</text>
			</row>
			<p
					border="frame.png 10"
					font="NotoSans 18"
					color="#ffffff"
					margin="4px"
					padding="16px 24px"
					maxWidth="50em"
					maxHeight="8lh"
					scrollbar="scrollbar.png 6">
It even supports scrollable paragraphs of text!
{{ range .Paragraphs }}
{{ . }}
{{ end }}
			</p>
			<row justify="start center">
				<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
					See how it works by enabling debug mode with CTRL-D or by toggling
				</text>
				<button onClick="Debug" color="#ffffff" margin="12px" padding="12px" btn="button.png 6">Debug Mode</button>
			</row>
		</col>
	</col>`
}

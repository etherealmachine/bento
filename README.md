[![Go Reference](https://pkg.go.dev/badge/github.com/etherealmachine/bento.svg)](https://pkg.go.dev/github.com/etherealmachine/bento)
[![Build Status](https://github.com/etherealmachine/bento/workflows/Go/badge.svg)](https://github.com/etherealmachine/bento/actions?query=workflow%3AGo)
# Bento

[Demo](https://etherealmachine.github.io/bento/)

# An XML based UI builder for Ebiten
It's like HTML/React, but for Ebiten!

* Support for buttons, text, scrollable paragraphs, text input, and text areas
* Small but useful set of layout options to build a responsive UI

![Screenshot of Page 1 of the Demo](https://user-images.githubusercontent.com/460276/202970256-4555e26d-62b2-4e09-9edb-1cb490187237.png)
![Screenshot of Page 2 of the Demo](https://user-images.githubusercontent.com/460276/202970335-794b1d26-6b0c-4f5a-9fe6-982d47b84421.png)
![Screenshot of Page 1 in debug mode](https://user-images.githubusercontent.com/460276/202970503-b016aee2-d29a-478f-a0fa-8bff5fdf0024.png)

```
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
	return `<col grow="1" justify="start" margin="24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px" underline="true">
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
					Check under the hood by enabling debug mode with CTRL-D or by toggling
				</text>
				<button onClick="Debug" color="#ffffff" margin="12px" padding="12px" btn="button.png 6">Debug Mode</button>
			</row>
		</col>
	</col>`
}
```

## Supported Styling
* Font Name and Size `font="NotoSans 24"`
* Font Color `color="#ffffff"`
* Margin `margin="20px"`
* Padding `padding="10px"`
* Underline Text `underline="true"`
* minWidth, minHeight, maxWidth, maxHeight
* Horizontal/Vertical Growth `grow="1 2"`
* Horizontal/Vertical Justification `justify="start center"`
  * `start`, `end`, `center`, `evenly`, `around`, `between`
* X/Y Offset `offset="4 12"`
* Scale (`img` only) `scale="2"`


## Events and Callbacks
```
type Event struct {
	X, Y  int       # For Click and Hover events
	Box   *Box
	Type  EventType # Click, Hover, Change, Update
	Value string    # For Input and Textarea Change events
}
```

**onClick** is fired on the first click event on a box (mousedown). Any element can register an `onClick` handler, not just buttons.
The X and Y fields of the `Event` will have the mouse coordinates relative to the top left corner of the box.

**onHover** is fired if the mouse is over the bounds of the box (including the margin and padding). `onHover` is NOT fired if the mouse was just clicked
(`onClick` takes precedence).

**onChange** is fired if when the value of `input` or `textarea` element changes, the `Value` field of the event contains the new value of the element

**onUpdate** is fired every frame

**onDraw** is fired on the `canvas` element. An onDraw handler **must** take a single `*ebiten.Image` argument

Click, Hover, Change, and Update handlers **may** either take a single `*bento.Event` argument, or zero arguments.

All handlers **may** return a boolean as an optimization hint. If the handler returns true, bento will automatically recompute the template and regenerate the UI tree.
This is an expensive operation, so handlers should return false if no variables were changed during the callback.

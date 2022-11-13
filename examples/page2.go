package main

type Page2 struct {
}

func (p *Page2) UI() string {
	return `<col grow="1" justify="start" margin="24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">
			How does it work?
		</text>
		<p
				grow="1"
				border="frame.png 10"
				font="RobotoMono 16"
				color="#ffffff"
				margin="4px"
				padding="16px 24px">
			<![CDATA[
// Define a struct to hold the UI state
type MyUI struct {
	Clicks int
}

// Add some event handlers
func (u *MyUI) Click(event *bento.Event) {
	u.Clicks++
}

// Layout the UI with XML
// Event handlers get attached with attributes
// text/template can be used to modify the content based on the state
func (u *MyUI) UI() string {
	return ` + "`" + `
		<col grow="1" justify="start" margin="24px">
			<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">
				Your text goes here
			</text>
			<row justify="start center">
				<button onClick="Click" color="#ffffff" margin="12px" padding="12px" btn="button.png 6">
					Clickable Button
				</button>
				<text font="NotoSans 18" color="#ffffff" margin="4px" padding="12px">
					{{ "{{ .Clicks }}" }} times
				</text>
			</row>
		</col>` + "`" + `
}
			]]>	
		</p>
	</col>`
}

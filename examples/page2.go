package main

type Page2 struct {
}

func (p *Page2) UI() string {
	return `<col grow="1" justify="start" margin="24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">
			How does it work?
		</text>
		<p
				border="frame.png 10"
				font="RobotoMono 18"
				color="#ffffff"
				margin="4px"
				padding="16px 24px">
			Very well, thanks!
		</p>
	</col>`
}

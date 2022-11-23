package main

type Page4 struct {
}

func (p *Page4) UI() string {
	return `<col grow="1" justify="center" margin="24px">
		<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">
			github.com/etherealmachine/bento
		</text>
	</col>`
}

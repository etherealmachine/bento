package main

type Page1 struct {
	Clicks int
}

func (p *Page1) Click() {
	p.Clicks++
}

func (p *Page1) Reset() {
	p.Clicks = 0
}

func (p *Page1) UI() string {
	return `<col grow="1" justify="center" border="frame.png 10">
		<row grow="1" justify="center">
			<text font="NotoSans 24" color="#ffffff" margin="4px" padding="12px">Page 1</text>
		</row>
		<row grow="1" justify="evenly center">
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum A</text>
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum B</text>
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum C</text>
		</row>
		<row grow="1" justify="around center">
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum A</text>
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum B</text>
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum C</text>
		</row>
		<row grow="1" justify="between center">
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum A</text>
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem ipsum B</text>
			<text font="RobotoMono 20" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum C</text>
		</row>
		<row grow="1" justify="start start" margin="0 0 0 48px">
			<text font="RobotoMono 20" color="#ffffff" padding="8px">a</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">b</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">c</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">A</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">B</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">C</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">1</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">2</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">3</text>
		</row>
		<row grow="1" justify="start center" margin="0 0 0 48px">
			<text font="RobotoMono 20" color="#ffffff" padding="8px">a</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">b</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">c</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">A</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">B</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">C</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">1</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">2</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">3</text>
		</row>
		<row grow="1" justify="start end" margin="0 0 0 48px">
			<text font="RobotoMono 20" color="#ffffff" padding="8px">a</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">b</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">c</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">A</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">B</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">C</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">1</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">2</text>
			<text font="RobotoMono 20" color="#ffffff" padding="8px">3</text>
		</row>
		<row grow="1" justify="end end">
			<text font="RobotoMono 24" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum</text>
		</row>
		<row grow="2" justify="end center">
			<text font="RobotoMono 24" color="#ffffff" margin="4px" padding="12px">Lorem Ipsum</text>
			<img src="profile.png" margin="4px" padding="8px">
				<button onClick="EditProfile" color="#ffffff" margin="4px" padding="12px" btn="button.png 6" offset="-1 -1">Edit</button>
			</img>
		</row>
		<row grow="1 0" justify="between">
			<button onClick="Click" color="#ffffff" margin="4px" padding="12px" btn="button.png 6">Clicks: {{.Clicks}}</button>
			<button onClick="Reset" color="#ffffff" margin="4px" padding="12px" btn="button.png 6">Reset</button>
		</row>
	</col>`
}

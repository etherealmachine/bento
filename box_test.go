package bento

import "testing"

type TestComponent struct {
	Items []string
}

func (c *TestComponent) UI() string {
	return `<col>
		<input />
		<row>
			{{ range .Items }}
				<text>{{ . }}</text>
			{{ end }}
		</row>
	</col>`
}

func TestUpdate(t *testing.T) {
	c := &TestComponent{}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}
	box.Update()
}

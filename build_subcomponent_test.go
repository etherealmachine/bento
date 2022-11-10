package bento

import (
	"testing"
)

type ComponentWithSub struct {
	Count int
	sub   *SubComponent
}

func (c *ComponentWithSub) SubComponent() *SubComponent {
	if c.sub != nil {
		return c.sub
	}
	c.sub = &SubComponent{Count: 2}
	return c.sub
}

func (c *ComponentWithSub) UI() string {
	return `<col>
		{{ if eq .Count 1 }}
			<SubComponent />
		{{ end }}
		<text>Hello World</text>
	</col>`
}

type SubComponent struct {
	Count int
}

func (c *SubComponent) UI() string {
	return `<col>
		<text>{{ .Count }}</text>
	</col>`
}

func TestBuildSubcomponent(t *testing.T) {
	c := &ComponentWithSub{Count: 1}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}
	want := `col <ComponentWithSub>
	col <SubComponent>
		text 2
	text Hello World
`
	got := box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

func TestRebuildSubcomponent(t *testing.T) {
	c := &ComponentWithSub{Count: 1}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}
	want := `col <ComponentWithSub>
	col <SubComponent>
		text 2
	text Hello World
`
	got := box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	c.Count = 0
	if err := box.Rebuild(); err != nil {
		t.Fatal(err)
	}
	want = `col <ComponentWithSub>
	text Hello World
`
	got = box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	c.Count = 1
	c.sub.Count = 3
	if err := box.Rebuild(); err != nil {
		t.Fatal(err)
	}
	want = `col <ComponentWithSub>
	col <SubComponent>
		text 3
	text Hello World
`
	got = box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

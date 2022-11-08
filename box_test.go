package bento

import (
	"fmt"
	"strings"
	"testing"
)

type TestComponent struct {
	Count int
	Array []string
	Map   map[string]string
	sub   *TestSubcomponent
}

func (c *TestComponent) Incr() {
	c.Count++
}

func (c *TestComponent) Decr() {
	c.Count++
}

func (c *TestComponent) TestSubcomponent() *TestSubcomponent {
	c.sub = &TestSubcomponent{Label: "foo"}
	return c.sub
}

func (c *TestComponent) UI() string {
	return `<col>
		{{ if lt .Count 3 }}
			<TestSubcomponent />
		{{ end}}
		<row>
			<text>{{ .Count }}</text>
			<button onClick="Incr">Increment</button>
			<button onClick="Decr">Decrement</button>
		</row>
		<col>
		{{ range .Array }}
			<text>{{ . }}</text>
		{{ end }}
		</col>
		<col>
		{{ range $key, $value := .Map }}
			<text>{{ $key }}: {{ $value }}</text>
		{{ end }}
		</col>
	</col>`
}

type TestSubcomponent struct {
	Label string
}

func (c *TestSubcomponent) UI() string {
	return `<col>
		<text>{{ .Label }}</text>
	</col>`
}

func TestBuild(t *testing.T) {
	c := &TestComponent{}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.TrimSpace(`
col TestComponent
	col TestSubcomponent
		text TestSubcomponent foo
	row TestComponent
		text TestComponent 0
		button TestComponent Increment
		button TestComponent Decrement
	col TestComponent
	col TestComponent
	`)
	got := strings.TrimSpace(box.String())
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

func TestRebuild(t *testing.T) {
	c := &TestComponent{}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		c.Count = i
		c.sub.Label = fmt.Sprintf("foo %d", i)
		new := &Box{component: c}
		if err := new.build(box); err != nil {
			t.Fatal(err)
		}
		*box = *new
		want := strings.TrimSpace(fmt.Sprintf(`
col TestComponent
	col TestSubcomponent
		text TestSubcomponent foo %d
	row TestComponent
		text TestComponent %d
		button TestComponent Increment
		button TestComponent Decrement
	col TestComponent
	col TestComponent
	`, i, i))
		got := strings.TrimSpace(box.String())
		if got != want {
			t.Fatalf("got\n%s\nwant\n%s\n", got, want)
		}
	}
	c.Count = 3
	new := &Box{component: c}
	if err := new.build(box); err != nil {
		t.Fatal(err)
	}
	*box = *new
	want := strings.TrimSpace(`
col TestComponent
	row TestComponent
		text TestComponent 3
		button TestComponent Increment
		button TestComponent Decrement
	col TestComponent
	col TestComponent
	`)
	got := strings.TrimSpace(box.String())
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

package bento

import (
	"testing"
)

type BasicComponent struct {
	Count int
	Array []string
	Map   map[string]string
}

func (c *BasicComponent) UI() string {
	return `<col>
		<row>
			<text>{{ .Count }}</text>
		</row>
		{{ if eq .Count 1 }}
		<col>
			<text>One</text>
		</col>
		{{ end }}
		{{ if .Array }}
			<col>
			{{ range .Array }}
				<text>{{ . }}</text>
			{{ end }}
			</col>
		{{ end }}
		{{ if .Map }}
			<col>
			{{ range $key, $value := .Map }}
				<text>{{ $key }}: {{ $value }}</text>
			{{ end }}
			</col>
		{{ end }}
	</col>`
}

func TestBuild(t *testing.T) {
	c := &BasicComponent{
		Count: 5,
		Array: []string{"a", "b", "c"},
		Map: map[string]string{
			"foo": "bar",
			"bar": "baz",
		},
	}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}
	want := `col <BasicComponent>
	row
		text 5
	col
		text a
		text b
		text c
	col
		text bar: baz
		text foo: bar
`
	got := box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
	nodes := 0
	box.visit(0, func(_ int, n *Box) error {
		nodes++
		for _, child := range n.Children {
			if child.Parent != n {
				t.Fatal("child has incorrect parent")
			}
		}
		return nil
	})
	if want := 10; nodes != want {
		t.Fatalf("got %d nodes, want %d", nodes, want)
	}
}

func TestRebuild(t *testing.T) {
	c := &BasicComponent{}
	box, err := Build(c)
	if err != nil {
		t.Fatal(err)
	}

	c.Count = 0
	if err := box.Rebuild(); err != nil {
		t.Fatal(err)
	}
	want := `col <BasicComponent>
	row
		text 0
`
	got := box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	nodes := 0
	box.visit(0, func(_ int, n *Box) error {
		nodes++
		for _, child := range n.Children {
			if child.Parent != n {
				t.Fatal("parent has kidnapped a child")
			}
		}
		return nil
	})
	if want := 3; nodes != want {
		t.Fatalf("got %d nodes, want %d", nodes, want)
	}

	c.Count = 1
	if err := box.Rebuild(); err != nil {
		t.Fatal(err)
	}
	want = `col <BasicComponent>
	row
		text 1
	col
		text One
`
	got = box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	c.Count = 2
	if err := box.Rebuild(); err != nil {
		t.Fatal(err)
	}
	want = `col <BasicComponent>
	row
		text 2
`
	got = box.String()
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

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

	nodes := 0
	box.visit(0, func(_ int, n *Box) error {
		nodes++
		for _, child := range n.Children {
			if child.Parent != n {
				t.Fatal("child has incorrect parent")
			}
		}
		return nil
	})
	if want := 4; nodes != want {
		t.Fatalf("got %d nodes, want %d", nodes, want)
	}
}

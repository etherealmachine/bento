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

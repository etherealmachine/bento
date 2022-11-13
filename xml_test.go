package bento

import (
	"encoding/xml"
	"log"
	"testing"
)

func TestUnmarshalXML(t *testing.T) {
	got := &Box{}
	if err := xml.Unmarshal([]byte(`<col>
		<row>
			<p>Hello World</p>
		</row>
		<row>
			<p>
				<![CDATA[
<Foo />
<Bar />
<Baz />
				]]>
			</p>
		</row>
	</col>`), got); err != nil {
		t.Fatal(err)
	}
	want := &Box{
		Tag: "col",
		Children: []*Box{
			{
				Tag: "row",
				Children: []*Box{
					{Tag: "p", Content: "Hello World"},
				},
			},
			{
				Tag: "row",
				Children: []*Box{
					{Tag: "p", Content: "<Foo />\n<Bar />\n<Baz />"},
				},
			},
		},
	}
	if err := got.diff(want); err != nil {
		log.Fatal(err)
	}
}

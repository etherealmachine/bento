package bento

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

var allowedTags = []string{
	"row",
	"col",
	"p",
	"text",
	"button",
	"img",
	"input",
	"textarea",
	"canvas",
}

func (n *Box) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Tag = start.Name.Local
	allowed := false
	for _, tag := range allowedTags {
		if n.Tag == tag {
			allowed = true
		}
	}
	if !allowed {
		if r, _ := utf8.DecodeRuneInString(n.Tag); !unicode.IsUpper(r) {
			return fmt.Errorf("unsupported tag %s, allowed tags: %v", n.Tag, allowedTags)
		}
	}
	n.Attrs = make(map[string]string)
	for _, attr := range start.Attr {
		n.Attrs[attr.Name.Local] = attr.Value
	}
	for {
		next, err := d.Token()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		switch next := next.(type) {
		case xml.StartElement:
			child := &Box{}
			if err := child.UnmarshalXML(d, next); err != nil {
				return err
			}
			n.Children = append(n.Children, child)
			child.Parent = n
		case xml.CharData:
			n.Content += strings.TrimSpace(string(next))
		case xml.EndElement:
			return nil
		case xml.ProcInst:
			return fmt.Errorf("unsupported xml processing instruction (<?target inst?>)")
		case xml.Directive:
			return fmt.Errorf("unsupported xml processing instruction (<!text>)")
		default:
			return fmt.Errorf("unsupported xml type %s", reflect.TypeOf(next))
		}
	}
}

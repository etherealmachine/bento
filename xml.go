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
	"grid",
	"row",
	"col",
	"p",
	"text",
	"button",
	"img",
}

func (n *node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.tag = start.Name.Local
	allowed := false
	for _, tag := range allowedTags {
		if n.tag == tag {
			allowed = true
		}
	}
	if !allowed {
		if r, _ := utf8.DecodeRuneInString(n.tag); !unicode.IsUpper(r) {
			return fmt.Errorf("unsupported tag %s, allowed tags: %v", n.tag, allowedTags)
		}
	}
	n.attrs = make(map[string]string)
	for _, attr := range start.Attr {
		n.attrs[attr.Name.Local] = attr.Value
	}
	for {
		next, err := d.Token()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		switch next.(type) {
		case xml.StartElement:
			child := &node{}
			if err := child.UnmarshalXML(d, next.(xml.StartElement)); err != nil {
				return err
			}
			n.children = append(n.children, child)
			child.parent = n
		case xml.EndElement:
			return nil
		case xml.CharData:
			n.content = strings.TrimSpace(string(next.(xml.CharData)))
		case xml.ProcInst:
			return fmt.Errorf("unsupported xml processing instruction (<?target inst?>)")
		case xml.Directive:
			return fmt.Errorf("unsupported xml processing instruction (<!text>)")
		default:
			return fmt.Errorf("unsupported xml type %s", reflect.TypeOf(next))
		}
	}
}

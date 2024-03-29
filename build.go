/*
Bento is an XML based UI builder for Ebiten
*/
package bento

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

func (n *Box) isSubcomponent() bool {
	r, _ := utf8.DecodeRuneInString(n.Tag)
	return unicode.IsUpper(r)
}

func (n *Box) buildSubcomponent(prev *Box) error {
	var subComponent reflect.Value
	if m := reflect.ValueOf(n.Component).MethodByName(n.Tag); m.IsValid() {
		res := m.Call(nil)
		subComponent = res[0]
	} else if field := reflect.ValueOf(n.Component).Elem().FieldByName(n.Tag); field.IsValid() {
		subComponent = field
	} else {
		return fmt.Errorf("%s must have a field or method named %s that returns a bento.Component", reflect.TypeOf(n.Component), n.Tag)
	}
	if style, ok := subComponent.Interface().(*Style); ok {
		n.Style = *style
		n.Tag = style.Extends
		return nil
	} else if sub, ok := subComponent.Interface().(Component); ok {
		subNode := &Box{
			Component: sub,
			Parent:    n.Parent,
		}
		if err := subNode.build(prev); err != nil {
			return err
		}
		*n = *subNode
		for _, child := range n.Children {
			child.Parent = n
		}
		return nil
	}
	return fmt.Errorf("%s.%s must return either Style or Component", reflect.TypeOf(n.Component), n.Tag)
}

func (n *Box) build(prev *Box) error {
	if n.Tag == "" {
		if err := n.expandComponent(); err != nil {
			return err
		}
	}
	if n.isSubcomponent() {
		return n.buildSubcomponent(prev)
	}
	n.Style.adopt(n)
	if err := n.Style.parseAttributes(); err != nil {
		return err
	}
	if prev != nil && n.Tag == prev.Tag {
		n.state = prev.state
		n.editable = prev.editable
		n.scrollable = prev.scrollable
	} else if n.Tag == "input" || n.Tag == "textarea" {
		n.editable = &Editable{}
	}
	if !n.Style.Display || n.Style.Hidden {
		return nil
	}
	i := 0
	for _, child := range n.Children {
		if child.Component == nil {
			child.Component = n.Component
		}
		var prevChild *Box
		if prev != nil && i < len(prev.Children) {
			prevChild = prev.Children[i]
		}
		if err := child.build(prevChild); err != nil {
			return err
		}
		if prevChild != nil && child.Tag == prevChild.Tag {
			i++
		}
	}
	return nil
}

func (n *Box) expandComponent() error {
	tmpl, err := template.New("").Parse(n.Component.UI())
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, n.Component); err != nil {
		return err
	}
	if err := xml.Unmarshal(buf.Bytes(), n); err != nil {
		// TODO: Good error reporting if parsing the XML fails
		re := regexp.MustCompile(`error on line (\d+)`)
		matches := re.FindStringSubmatch(err.Error())
		if len(matches) == 2 {
			lineNo, _ := strconv.Atoi(matches[1])
			lines := strings.Split(buf.String(), "\n")
			ctx := new(bytes.Buffer)
			for i, line := range lines {
				if i >= lineNo-3 && i <= lineNo+3 {
					ctx.WriteString(fmt.Sprintf("%d: %s", i+1, line))
					ctx.WriteByte('\n')
				}
			}
			return fmt.Errorf("error building %s: %s\n%s", reflect.ValueOf(n.Component).Elem().Type().Name(), err, ctx)
		}
		return fmt.Errorf("error building %s: %s", reflect.ValueOf(n.Component).Elem().Type().Name(), err)
	}
	return nil
}

func (n *Box) visit(depth int, f func(depth int, n *Box) error) error {
	if n == nil {
		return nil
	}
	if err := f(depth, n); err != nil {
		return err
	}
	for _, c := range n.Children {
		if err := c.visit(depth+1, f); err != nil {
			return err
		}
	}
	return nil
}

func (a *Box) diff(b *Box) error {
	if a == nil && b != nil {
		return fmt.Errorf("a is nil, b is present")
	}
	if a != nil && b == nil {
		return fmt.Errorf("b is nil, a is present")
	}
	if a.Tag != b.Tag {
		return fmt.Errorf("tag mismatch, %s != %s", a.Tag, b.Tag)
	}
	if a.Content != b.Content {
		return fmt.Errorf("content mismatch, %q != %q", a.Content, b.Content)
	}
	if len(a.Children) != len(b.Children) {
		return fmt.Errorf("different number of children, %d != %d", len(a.Children), len(b.Children))
	}
	for i, c := range a.Children {
		if err := c.diff(b.Children[i]); err != nil {
			return err
		}
	}
	return nil
}

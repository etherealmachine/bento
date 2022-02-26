package bento

import (
	"encoding/hex"
	"encoding/xml"
	"image/color"
	"log"
	"regexp"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var (
	sizeSpec = regexp.MustCompile(`(\d+)(em|px)`)
)

type Spacing struct {
	Top, Right, Bottom, Left int `xml:",attr,omitempty"`
}

type Style struct {
	Extends             string            `xml:",omitempty"`
	Attrs               map[string]string `xml:"-"`
	FontName            string            `xml:",omitempty"`
	FontSize            int               `xml:",omitempty"`
	Font                font.Face         `xml:"-"`
	NineSlice           *NineSlice        `xml:"-"`
	Scrollbar           *NineSlice        `xml:"-"`
	Background          *ebiten.Image     `xml:"-"`
	MinWidth, MinHeight int               `xml:",omitempty"`
	MaxWidth, MaxHeight int               `xml:",omitempty"`
	Margin, Padding     *Spacing          `xml:",omitempty"`
	Color               *color.RGBA       `xml:"-"`
}

func (s *Style) Adopt(attrs map[string]string) error {
	if s.Attrs == nil {
		s.Attrs = make(map[string]string)
	}
	for k, v := range attrs {
		if _, exists := s.Attrs[k]; !exists {
			s.Attrs[k] = v
		}
	}
	if margin, exists := s.Attrs["margin"]; exists {
		i := parseSize(margin, s.Font)
		if s.Margin == nil {
			s.Margin = &Spacing{}
		}
		// TODO: multiple margins, space separated
		s.Margin.Top = i
		s.Margin.Right = i
		s.Margin.Bottom = i
		s.Margin.Left = i
	}
	if padding, exists := s.Attrs["padding"]; exists {
		i := parseSize(padding, s.Font)
		if s.Padding == nil {
			s.Padding = &Spacing{}
		}
		// TODO: multiple paddings, space separated
		s.Padding.Top = i
		s.Padding.Right = i
		s.Padding.Bottom = i
		s.Padding.Left = i
	}
	if path := s.Attrs["path"]; path != "" {
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatal(err)
			return err
		}
		s.Background = img
	}
	s.MinWidth = s.parseSize("minWidth")
	s.MaxWidth = s.parseSize("maxWidth")
	s.MinHeight = s.parseSize("minHeight")
	s.MaxHeight = s.parseSize("maxHeight")
	return nil
}

func (s *Style) parseSize(attr string) int {
	if size, exists := s.Attrs[attr]; exists {
		return parseSize(size, s.Font)
	}
	return 0
}

func parseSize(spec string, f font.Face) int {
	matches := sizeSpec.FindStringSubmatch(spec)
	if len(matches) == 3 {
		size, _ := strconv.Atoi(matches[1])
		if matches[2] == "px" {
			return size
		} else {
			return size * text.BoundString(f, "—").Dx()
		}
	}
	return 0
}

func (s *Style) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	for k, v := range s.Attrs {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}
	if s.Color != nil {
		start.Attr = append(start.Attr,
			xml.Attr{
				Name:  xml.Name{Local: "Color"},
				Value: "0x" + hex.EncodeToString([]byte{s.Color.R, s.Color.G, s.Color.B, s.Color.A}),
			})
	}
	type style Style
	return e.EncodeElement((*style)(s), start)
}

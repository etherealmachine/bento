package bento

import (
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var (
	sizeSpec = regexp.MustCompile(`(\d+)(em|px|%)`)
)

type Spacing struct {
	Top, Right, Bottom, Left int `xml:",attr,omitempty"`
}

type Style struct {
	Extends             string
	Attrs               map[string]string
	FontName            string
	FontSize            int
	Font                font.Face
	Border              *NineSlice
	Button              *Button
	Scrollbar           *Scrollbar
	Background          *ebiten.Image
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	Margin, Padding     *Spacing
	Color               *color.RGBA
	Hidden              bool
	Display             bool
}

func (s *Style) adopt(attrs map[string]string) error {
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
	if spec := s.Attrs["border"]; spec != "" {
		a := strings.Split(spec, " ")
		path := a[0]
		var widths, heights [3]int
		for i := 0; i < 3; i++ {
			var err error
			widths[i], err = strconv.Atoi(a[i+1])
			heights[i], err = strconv.Atoi(a[i+4])
			if err != nil {
				return err
			}
		}
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatal(err)
			return err
		}
		s.Border = NewNineSlice(img, widths, heights, 0, 0)
	}
	// TODO
	s.Display = true
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

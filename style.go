package bento

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	bentotext "github.com/etherealmachine/bento/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var (
	sizeSpec = regexp.MustCompile(`(\d+)(em|px|lh)`)
)

type Justification string

const (
	// The items are packed flush to each other toward the start edge of the alignment container in the main axis.
	// [item1 item2              ]
	Start = Justification("start")
	// The items are packed flush to each other toward the end edge of the alignment container in the main axis.
	// [              item1 item2]
	End = Justification("end")
	// The items are packed flush to each other toward the center of the alignment container along the main axis.
	// [       item1 item2       ]
	Center = Justification("center")
	// The items are evenly distributed within the alignment container along the main axis.
	// The spacing between each pair of adjacent items, the main-start edge and the first item, and the main-end edge and the last item, are all exactly the same.
	// [     item1     item2     ]
	Evenly = Justification("evenly")
	// The items are evenly distributed within the alignment container along the main axis.
	// The spacing between each pair of adjacent items is the same.
	// The empty space before the first and after the last item equals half of the space between each pair of adjacent items.
	// [   item1         item2   ]
	Around = Justification("around")
	// The items are evenly distributed within the alignment container along the main axis.
	// The spacing between each pair of adjacent items is the same.
	// The first item is flush with the main-start edge, and the last item is flush with the main-end edge.
	// [item1               item2]
	Between = Justification("between")
)

func (j Justification) Valid() bool {
	return j == "start" || j == "end" || j == "center" || j == "between" || j == "around" || j == "evenly"
}

type Spacing struct {
	Top, Right, Bottom, Left int `xml:",attr,omitempty"`
}

type Style struct {
	Extends             string
	Attrs               map[string]string
	FontName            string
	FontSize            int
	Font                font.Face
	Underline           bool
	Border              *NineSlice
	Button              *[4]*NineSlice
	Scrollbar           *[3][4]*NineSlice
	Input               *[4]*NineSlice
	Image               *ebiten.Image
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	HJust, VJust        Justification
	HGrow, VGrow        int
	Margin, Padding     Spacing
	Color               color.Color
	OffsetX, OffsetY    int
	Float               bool
	Hidden              bool
	Display             bool
	ScaleX, ScaleY      float64
	ZIndex              int
	node                *Box
}

func (s *Style) adopt(node *Box) {
	if s.Attrs == nil {
		s.Attrs = make(map[string]string)
	}
	for k, v := range node.Attrs {
		if _, exists := s.Attrs[k]; !exists {
			s.Attrs[k] = v
		}
	}
	s.node = node
}

func (s *Style) scrollBarWidth() int {
	if s.node.ContentHeight > s.MaxHeight && s.Scrollbar != nil {
		return s.Scrollbar[0][0].Width()
	}
	return 0
}

func (n *Box) styleSize() {
	n.ContentWidth += n.style.scrollBarWidth()
	if n.style.MinWidth > 0 {
		n.ContentWidth = max(n.ContentWidth, n.style.MinWidth)
	}
	if n.style.MinHeight > 0 {
		n.ContentHeight = max(n.ContentHeight, n.style.MinHeight)
	}
	if n.Tag == "p" {
		metrics := n.style.Font.Metrics()
		minHeight := metrics.Height.Round()
		n.ContentHeight = max(n.ContentHeight, minHeight)
	}
	if n.style.MaxHeight > 0 && n.ContentHeight > n.style.MaxHeight {
		n.ContentHeight = n.style.MaxHeight
	}
}

func (s *Style) parseAttributes() error {
	var err error
	if s.Font == nil {
		if s.FontName != "" {
			size := 16
			if s.FontSize > 0 {
				size = s.FontSize
			}
			s.Font = bentotext.Font(s.FontName, size)
		}
	}
	if s.Font == nil {
		if s.FontName, s.FontSize, s.Font, err = parseFont(s.Attrs["font"]); err != nil {
			return fmt.Errorf("error parsing font: %s", err)
		}
	}
	if spec := s.Attrs["underline"]; spec == "true" {
		s.Underline = true
	}
	margin, err := parseSpacing(s.Attrs["margin"], s.Font)
	if err != nil {
		return fmt.Errorf("error parsing margin: %s", err)
	}
	s.Margin = *margin
	padding, err := parseSpacing(s.Attrs["padding"], s.Font)
	if err != nil {
		return fmt.Errorf("error parsing padding: %s", err)
	}
	s.Padding = *padding
	if s.Image == nil {
		if s.Image, err = loadImage(s.Attrs["src"]); err != nil {
			return fmt.Errorf("error parsing image src: %s", err)
		}
	}
	if s.Border == nil {
		if s.Border, err = loadNineSlice(s.Attrs["border"]); err != nil {
			return fmt.Errorf("error parsing border: %s", err)
		}
	}
	if s.MinWidth == 0 {
		if s.MinWidth, err = parseSize(s.Attrs["minWidth"], s.Font); err != nil {
			return fmt.Errorf("error parsing minWidth: %s", err)
		}
	}
	if s.MinHeight == 0 {
		if s.MinHeight, err = parseSize(s.Attrs["minHeight"], s.Font); err != nil {
			return fmt.Errorf("error parsing minHeight: %s", err)
		}
	}
	if s.MaxWidth == 0 {
		if s.MaxWidth, err = parseSize(s.Attrs["maxWidth"], s.Font); err != nil {
			return fmt.Errorf("error parsing maxWidth: %s", err)
		}
	}
	if s.MaxWidth != 0 && s.node.Tag != "p" && s.node.Tag != "textarea" {
		return fmt.Errorf("invalid tag %s: max width can only apply to paragraph (p) or textarea", s.node.Tag)
	}
	if s.MaxHeight == 0 {
		if s.MaxHeight, err = parseSize(s.Attrs["maxHeight"], s.Font); err != nil {
			return fmt.Errorf("error parsing maxHeight: %s", err)
		}
	}
	if s.MaxHeight != 0 && s.node.Tag != "p" && s.node.Tag != "textarea" {
		return fmt.Errorf("invalid tag %s: max height can only apply to paragraph (p) or textarea", s.node.Tag)
	}
	if spec := s.Attrs["justify"]; spec != "" {
		if s.HJust, s.VJust, err = parseJustification(spec); err != nil {
			return fmt.Errorf("error parsing justification: %s", err)
		}
	}
	if s.HJust == "" {
		s.HJust = Start
	}
	if s.VJust == "" {
		s.VJust = Start
	}
	if spec := s.Attrs["grow"]; spec != "" {
		if s.HGrow, s.VGrow, err = parseGrow(spec); err != nil {
			return fmt.Errorf("error parsing grow: %s", err)
		}
	}
	if s.Color == nil {
		if s.Color, err = parseColor(s.Attrs["color"]); err != nil {
			return fmt.Errorf("error parsing color: %s", err)
		}
	}
	if spec := s.Attrs["offset"]; spec != "" {
		if s.OffsetX, s.OffsetY, err = parseOffset(spec); err != nil {
			return fmt.Errorf("error parsing offset: %s", err)
		}
	}
	if spec := s.Attrs["scale"]; spec != "" {
		if s.ScaleX, s.ScaleY, err = parseScale(spec); err != nil {
			return fmt.Errorf("error parsing scale: %s", err)
		}
	} else {
		s.ScaleX, s.ScaleY = 1, 1
	}
	if spec := s.Attrs["zIndex"]; spec != "" {
		s.ZIndex, err = strconv.Atoi(spec)
		if err != nil {
			return fmt.Errorf("error parsing zIndex: %s", err)
		}
	}
	if s.Button == nil {
		if s.Button, err = ParseButton(s.Attrs["btn"]); err != nil {
			return fmt.Errorf("error parsing button: %s", err)
		}
	}
	if s.Scrollbar == nil {
		if s.Scrollbar, err = ParseScrollbar(s.Attrs["scrollbar"]); err != nil {
			return fmt.Errorf("error parsing scrollbar: %s", err)
		}
	}
	if s.Input == nil {
		if s.Input, err = ParseButton(s.Attrs["input"]); err != nil {
			return fmt.Errorf("error parsing input: %s", err)
		}
	}
	s.Float = s.Attrs["float"] == "true"
	s.Hidden = s.Attrs["hidden"] == "true"
	s.Display = s.Attrs["display"] != "false"
	return nil
}

// parse font spec e.g. "NotoSans 16"
func parseFont(spec string) (string, int, font.Face, error) {
	if spec == "" {
		return "NotoSans", 16, bentotext.Font("NotoSans", 16), nil
	}
	a := strings.Split(spec, " ")
	if len(a) != 2 {
		return "", 0, nil, fmt.Errorf("invalid font spec %s", spec)
	}
	size, err := strconv.Atoi(a[1])
	if err != nil {
		return "", 0, nil, err
	}
	return a[1], size, bentotext.Font(a[0], size), nil
}

// parse spacing spec e.g. "24px", "12px 12px", "8px 24px 6px 12px"
// units can be px or em
func parseSpacing(spec string, font font.Face) (*Spacing, error) {
	if spec == "" {
		return &Spacing{0, 0, 0, 0}, nil
	}
	a := strings.Split(spec, " ")
	if len(a) == 1 {
		size, err := parseSize(a[0], font)
		if err != nil {
			return nil, err
		}
		return &Spacing{size, size, size, size}, nil
	}
	if len(a) == 2 {
		x, err := parseSize(a[0], font)
		if err != nil {
			return nil, err
		}
		y, err := parseSize(a[1], font)
		if err != nil {
			return nil, err
		}
		return &Spacing{y, x, y, x}, nil
	}
	if len(a) == 4 {
		top, err := parseSize(a[0], font)
		if err != nil {
			return nil, err
		}
		right, err := parseSize(a[1], font)
		if err != nil {
			return nil, err
		}
		bottom, err := parseSize(a[2], font)
		if err != nil {
			return nil, err
		}
		left, err := parseSize(a[3], font)
		if err != nil {
			return nil, err
		}
		return &Spacing{top, right, bottom, left}, nil
	}
	return nil, fmt.Errorf("invalid spacing spec %s", spec)
}

func parseJustification(spec string) (Justification, Justification, error) {
	if spec == "" {
		return Start, Start, fmt.Errorf("invalid justification spec %q", spec)
	}
	justs := strings.Split(spec, " ")
	var hj, vj Justification
	if len(justs) == 1 {
		hj = Justification(justs[0])
		vj = hj
	} else {
		hj, vj = Justification(justs[0]), Justification(justs[1])
	}
	if !hj.Valid() {
		return Start, Start, fmt.Errorf("invalid justification %s", hj)
	}
	if !vj.Valid() {
		return Start, Start, fmt.Errorf("invalid justification %s", vj)
	}
	return hj, vj, nil
}

func parseGrow(spec string) (int, int, error) {
	if spec == "" {
		return 0, 0, nil
	}
	a := strings.Split(spec, " ")
	var hg, vg int
	var err error
	hg, err = strconv.Atoi(a[0])
	vg = hg
	if err == nil && len(a) == 2 {
		vg, err = strconv.Atoi(a[1])
	} else if len(a) > 2 {
		return 0, 0, fmt.Errorf("too many parameters for grow, expected at most 2: %s", spec)
	}
	return hg, vg, err
}

func parseOffset(spec string) (int, int, error) {
	if spec == "" {
		return 0, 0, nil
	}
	a := strings.Split(spec, " ")
	var x, y int
	var err error
	x, err = strconv.Atoi(a[0])
	y = x
	if err == nil && len(a) == 2 {
		y, err = strconv.Atoi(a[1])
	} else if len(a) > 2 {
		return 0, 0, fmt.Errorf("too many parameters for offset, expected at most 2: %s", spec)
	}
	return x, y, err
}

func parseScale(spec string) (float64, float64, error) {
	if spec == "" {
		return 1, 1, nil
	}
	a := strings.Split(spec, " ")
	var x, y float64
	var err error
	x, err = strconv.ParseFloat(a[0], 32)
	y = x
	if err == nil && len(a) == 2 {
		y, err = strconv.ParseFloat(a[1], 32)
	} else if len(a) > 2 {
		return 0, 0, fmt.Errorf("too many parameters for scale, expected at most 2: %s", spec)
	}
	return x, y, err
}

func loadImage(spec string) (*ebiten.Image, error) {
	if spec == "" {
		return nil, nil
	}
	img, _, err := ebitenutil.NewImageFromFile(spec)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// space-separated string with 2, 4 or 7 items
// path, margin (integer)
// or
// path, widths/heights (3 integers in pixels)
// or
// path, widths (3 integers in pixels), heights (3 integers in pixels)
// e.g. "frame.png 10"
// which, for a 32px sized image is equivaleint to
// "frame.png 10 12 10" and "frame.png 10 12 10 10 12 10"
func loadNineSlice(spec string) (*NineSlice, error) {
	if spec == "" {
		return nil, nil
	}
	img, widths, heights, err := loadImageFromSpec(spec, 1)
	if err != nil {
		return nil, err
	}
	return NewNineSlice(img, *widths, *heights, 0, 0), nil
}

// e.g. "button.png 6"
func ParseButton(spec string) (*[4]*NineSlice, error) {
	if spec == "" {
		return nil, nil
	}
	img, widths, heights, err := loadImageFromSpec(spec, 4)
	if err != nil {
		return nil, err
	}
	var states [4]*NineSlice
	w := widths[0] + widths[1] + widths[2]
	for i := 0; i < 4; i++ {
		states[i] = NewNineSlice(img, *widths, *heights, w*i, 0)
	}
	return &states, nil
}

// e.g. "scrollbar.png 6"
func ParseScrollbar(spec string) (*[3][4]*NineSlice, error) {
	if spec == "" {
		return nil, nil
	}
	img, widths, _, err := loadImageFromSpec(spec, 4)
	if err != nil {
		return nil, err
	}
	var states [3][4]*NineSlice
	w := widths[0] + widths[1] + widths[2]
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			states[i][j] = NewNineSlice(img, *widths, *widths, w*i, w*j)
		}
	}
	return &states, nil
}

func loadImageFromSpec(spec string, frames int) (*ebiten.Image, *[3]int, *[3]int, error) {
	a := strings.Split(spec, " ")
	path := a[0]
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, nil, nil, err
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	var widths, heights [3]int
	if len(a) == 2 {
		margin, err := strconv.Atoi(a[1])
		if err != nil {
			return nil, nil, nil, err
		}
		widths[0] = margin
		widths[2] = margin
		widths[1] = w/frames - 2*margin
		heights[0] = margin
		heights[2] = margin
		heights[1] = h - 2*margin
	} else if len(a) == 4 {
		start, err := strconv.Atoi(a[1])
		if err != nil {
			return nil, nil, nil, err
		}
		middle, err := strconv.Atoi(a[2])
		if err != nil {
			return nil, nil, nil, err
		}
		end, err := strconv.Atoi(a[3])
		if err != nil {
			return nil, nil, nil, err
		}
		widths[0] = start
		widths[1] = middle
		widths[2] = end
		heights[0] = start
		heights[1] = middle
		heights[2] = end
	} else if len(a) == 7 {
		for i := 0; i < 3; i++ {
			var err error
			widths[i], err = strconv.Atoi(a[i+1])
			if err != nil {
				return nil, nil, nil, err
			}
			heights[i], err = strconv.Atoi(a[i+4])
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}
	return img, &widths, &heights, nil
}

func parseColor(spec string) (color.Color, error) {
	c := &color.RGBA{A: 0xff}
	if spec == "" {
		return c, nil
	}
	if spec[0] != '#' {
		return nil, fmt.Errorf("invalid color string %s", spec)
	}

	var err error
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		default:
			err = fmt.Errorf("invalid color string %s", spec)
			return 0
		}
	}

	switch len(spec) {
	case 7:
		c.R = hexToByte(spec[1])<<4 + hexToByte(spec[2])
		c.G = hexToByte(spec[3])<<4 + hexToByte(spec[4])
		c.B = hexToByte(spec[5])<<4 + hexToByte(spec[6])
	case 4:
		c.R = hexToByte(spec[1]) * 17
		c.G = hexToByte(spec[2]) * 17
		c.B = hexToByte(spec[3]) * 17
	default:
		return nil, fmt.Errorf("invalid color string %s", spec)
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func parseSize(spec string, f font.Face) (int, error) {
	matches := sizeSpec.FindStringSubmatch(spec)
	if len(matches) == 3 {
		size, _ := strconv.Atoi(matches[1])
		switch matches[2] {
		case "px":
			return size, nil
		case "em":
			return size * text.BoundString(f, "M").Dx(), nil
		case "lh":
			return size * f.Metrics().Height.Round(), nil
		}
	}
	return 0, nil
}

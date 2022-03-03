package bento

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
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
	sizeSpec = regexp.MustCompile(`(\d+)(em|px|%)`)
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
	Border              *NineSlice
	Button              *Button
	Scrollbar           *Scrollbar
	Image               *ebiten.Image
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	HJust, VJust        Justification
	HGrow, VGrow        int
	Margin, Padding     *Spacing
	Color               color.Color
	Hidden              bool
	Display             bool
	node                *Box
}

func (s *Style) adopt(node *Box) {
	if s.Attrs == nil {
		s.Attrs = make(map[string]string)
	}
	for k, v := range node.attrs {
		if _, exists := s.Attrs[k]; !exists {
			s.Attrs[k] = v
		}
	}
	s.node = node
}

func (s *Style) margin() (int, int, int, int) {
	if s != nil {
		m := s.Margin
		if m != nil {
			return m.Top, m.Right, m.Bottom, m.Left
		}
	}
	return 0, 0, 0, 0
}

func (s *Style) padding() (int, int, int, int) {
	if s != nil {
		p := s.Padding
		if p != nil {
			return p.Top, p.Right, p.Bottom, p.Left
		}
	}
	return 0, 0, 0, 0
}

func (s *Style) justification() (Justification, Justification) {
	hj, vj := Start, Start
	if s != nil {
		if s.HJust != "" {
			hj = s.HJust
		}
		if s.VJust != "" {
			vj = s.VJust
		}
	}
	return hj, vj
}

func (s *Style) growth() (int, int) {
	hg, vg := 0, 0
	if s != nil {
		hg = s.HGrow
		vg = s.VGrow
	}
	return hg, vg
}

func (s *Style) display() bool {
	return s.node.templateAttr("display", true)
}

func (s *Style) hidden() bool {
	return s.node.templateAttr("hidden", false)
}

func (n *Box) styleSize() {
	if n.style == nil {
		return
	}
	if n.style.MinWidth > 0 {
		n.ContentWidth = max(n.ContentWidth, n.style.MinWidth)
	}
	if n.style.MinHeight > 0 {
		n.ContentHeight = max(n.ContentHeight, n.style.MinHeight)
	}
	if n.tag == "p" {
		textHeight := n.TextBounds.Dy()
		metrics := n.style.Font.Metrics()
		minHeight := metrics.Height.Round()
		n.ContentHeight = max(textHeight, minHeight)
	}
	if n.style.MaxHeight > 0 && n.ContentHeight > n.style.MaxHeight {
		n.ContentHeight = n.style.MaxHeight
	}
	if n.style.Button != nil {
		n.style.Button.rect = n.innerRect()
	}
	if n.style.Scrollbar != nil {
		inner := n.innerRect()
		n.style.Scrollbar.x = inner.Max.X
		n.style.Scrollbar.y = inner.Min.Y
		n.style.Scrollbar.height = inner.Dy()
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
	if s.Margin == nil {
		if s.Margin, err = parseSpacing(s.Attrs["margin"], s.Font); err != nil {
			return fmt.Errorf("error parsing margin: %s", err)
		}
	}
	if s.Padding == nil {
		if s.Padding, err = parseSpacing(s.Attrs["padding"], s.Font); err != nil {
			return fmt.Errorf("error parsing padding: %s", err)
		}
	}
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
	if s.MaxWidth != 0 && s.node.tag != "p" {
		return fmt.Errorf("invalid tag %s: max width can only apply to a paragraph (p) tag", s.node.tag)
	}
	if s.MaxHeight == 0 {
		if s.MaxHeight, err = parseSize(s.Attrs["maxHeight"], s.Font); err != nil {
			return fmt.Errorf("error parsing maxHeight: %s", err)
		}
	}
	if s.MaxHeight != 0 && s.node.tag != "p" {
		return fmt.Errorf("invalid tag %s: max height can only apply to a paragraph (p) tag", s.node.tag)
	}
	if spec := s.Attrs["justify"]; spec != "" {
		if s.HJust, s.VJust, err = parseJustification(spec); err != nil {
			return fmt.Errorf("error parsing justification: %s", err)
		}
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
	if s.Button == nil {
		if s.Button, err = parseButton(s.Attrs["btn"]); err != nil {
			return fmt.Errorf("error parsing button: %s", err)
		}
		if s.Button != nil {
			log.Printf("creating button %s for %s", s.node.content, reflect.TypeOf(s.node.component))
			s.Button.onClick = func(id string) {
				m := reflect.ValueOf(s.node.component).MethodByName(s.Attrs["onClick"])
				if !m.IsValid() {
					log.Println(reflect.TypeOf(s.node.component))
					log.Println(reflect.TypeOf(s.node.context))
				}
				m.Call([]reflect.Value{reflect.ValueOf(id)})
			}
		}
	}
	// TODO
	// s.Scrollbar
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
// unints can be px or em
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
		return Start, Start, nil
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
func parseButton(spec string) (*Button, error) {
	if spec == "" {
		return nil, nil
	}
	img, widths, heights, err := loadImageFromSpec(spec, 4)
	if err != nil {
		return nil, err
	}
	return NewButton(img, *widths, *heights, nil), nil
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
		if matches[2] == "px" {
			return size, nil
		} else {
			return size * text.BoundString(f, "—").Dx(), nil
		}
	}
	return 0, nil
}

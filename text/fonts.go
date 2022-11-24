package text

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"

	_ "embed"
)

var (
	fonts map[string]*fontMap = make(map[string]*fontMap)
	//go:embed NotoSans-Regular.ttf
	notoSans []byte
	//go:embed RobotoMono-Regular.ttf
	robotoMono []byte
)

type fontMap struct {
	font  *sfnt.Font
	sizes map[int]font.Face
}

func init() {
	fonts["NotoSans"] = loadFont(notoSans)
	fonts["RobotoMono"] = loadFont(robotoMono)
}

func LoadFontFromFile(name, path string) error {
	def, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fonts[name] = loadFont(def)
	return nil
}

func Font(name string, size int) font.Face {
	fm := fonts[name]
	if fm == nil {
		fm = fonts["NotoSans"]
	}
	f := fm.sizes[size]
	if f == nil {
		return fm.load(size)
	}
	return f
}

func (f *fontMap) load(size int) font.Face {
	const dpi = 72
	face, err := opentype.NewFace(f.font, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	f.sizes[size] = face
	return face
}

func loadFont(def []byte) *fontMap {
	tt, err := opentype.Parse(def)
	if err != nil {
		panic(err)
	}
	return &fontMap{
		font:  tt,
		sizes: make(map[int]font.Face),
	}
}

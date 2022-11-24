package text

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

/*
	Branched from https://github.com/hajimehoshi/ebiten/blob/6443af640115accec7b9a6f90ae579a53217688d/text/text.go
	Copyright 2017 The Ebiten Authors
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
	Package text offers functions to draw texts on an Ebiten's image.
	For the example using a TTF font, see font package in the examples.
*/

type Alignment string

const (
	Start  = Alignment("start")
	End    = Alignment("end")
	Center = Alignment("center")
)

type glyphImageCacheEntry struct {
	image *ebiten.Image
	atime int64
}

var (
	glyphBoundsCache  = map[font.Face]map[rune]fixed.Rectangle26_6{}
	glyphImageCache   = map[font.Face]map[rune]*glyphImageCacheEntry{}
	glyphAdvanceCache = map[font.Face]map[rune]fixed.Int26_6{}
)

func glyphBounds(face font.Face, r rune) fixed.Rectangle26_6 {
	if _, ok := glyphBoundsCache[face]; !ok {
		glyphBoundsCache[face] = map[rune]fixed.Rectangle26_6{}
	}
	if b, ok := glyphBoundsCache[face][r]; ok {
		return b
	}
	b, _, _ := face.GlyphBounds(r)
	glyphBoundsCache[face][r] = b
	return b
}

func glyphAdvance(face font.Face, r rune) fixed.Int26_6 {
	m, ok := glyphAdvanceCache[face]
	if !ok {
		m = map[rune]fixed.Int26_6{}
		glyphAdvanceCache[face] = m
	}

	a, ok := m[r]
	if !ok {
		a, _ = face.GlyphAdvance(r)
		m[r] = a
	}

	return a
}

func fixed26_6ToFloat64(x fixed.Int26_6) float64 {
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

func glyphImage(face font.Face, r rune) *ebiten.Image {
	if _, ok := glyphImageCache[face]; !ok {
		glyphImageCache[face] = map[rune]*glyphImageCacheEntry{}
	}

	if e, ok := glyphImageCache[face][r]; ok {
		e.atime = time.Now().Unix()
		return e.image
	}

	b := glyphBounds(face, r)
	w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
	if w == 0 || h == 0 {
		glyphImageCache[face][r] = &glyphImageCacheEntry{
			image: nil,
			atime: time.Now().Unix(),
		}
		return nil
	}

	if b.Min.X&((1<<6)-1) != 0 {
		w++
	}
	if b.Min.Y&((1<<6)-1) != 0 {
		h++
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	d := font.Drawer{
		Dst:  rgba,
		Src:  image.White,
		Face: face,
	}
	x, y := -b.Min.X, -b.Min.Y
	x, y = fixed.I(x.Ceil()), fixed.I(y.Ceil())
	d.Dot = fixed.Point26_6{X: x, Y: y}
	d.DrawString(string(r))

	img := ebiten.NewImageFromImage(rgba)
	if _, ok := glyphImageCache[face][r]; !ok {
		glyphImageCache[face][r] = &glyphImageCacheEntry{
			image: img,
			atime: time.Now().Unix(),
		}
	}

	return img
}

func drawGlyph(dst *ebiten.Image, face font.Face, r rune, dx, dy fixed.Int26_6, op *ebiten.DrawImageOptions) {
	img := glyphImage(face, r)
	if img == nil {
		return
	}

	b := glyphBounds(face, r)
	op2 := &ebiten.DrawImageOptions{}
	if op != nil {
		*op2 = *op
	}
	op2.GeoM.Reset()
	op2.GeoM.Translate(math.Floor(fixed26_6ToFloat64((dx + b.Min.X))), math.Floor(fixed26_6ToFloat64((dy + b.Min.Y))))
	if op != nil {
		op2.GeoM.Concat(op.GeoM)
	}
	dst.DrawImage(img, op2)
}

// BoundString returns the measured size of a given string using a given font.
// This method will return the exact size in pixels that a string drawn by Draw will be.
// The bound's origin point indicates the dot (period) position.
// This means that if the text consists of one character '.', this dot is rendered at (0, 0).
//
// This is very similar to golang.org/x/image/font's BoundString,
// but this BoundString calculates the actual rendered area considering multiple lines and space characters.
//
// face is the font for text rendering.
// text is the string that's being measured.
//
// Be careful that the passed font face is held by this package and is never released.
// This is a known issue (#498).
func BoundString(face font.Face, text string) image.Rectangle {
	fx, fy := fixed.I(0), fixed.I(0)
	prevR := rune(-1)

	var bounds fixed.Rectangle26_6
	for _, r := range text {
		if r == '\n' {
			continue
		}
		if prevR >= 0 {
			fx += face.Kern(prevR, r)
		}

		b := glyphBounds(face, r)
		if b.Max.X-b.Min.X == 0 {
			b.Max.Y = 1
			b.Max.X = glyphAdvance(face, r)
		}
		b.Min.X += fx
		b.Max.X += fx
		b.Min.Y += fy
		b.Max.Y += fy
		bounds = bounds.Union(b)

		fx += glyphAdvance(face, r)
		prevR = r
	}

	bounds = bounds.Union(glyphBounds(face, 'M'))

	return image.Rect(
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.X))),
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.Y))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.X))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.Y))),
	)
}

// DrawString draws a given text on a given destination image dst.
//
// face is the font for text rendering.
//
// clr is the color for text rendering.
//
// underline draws a line along the baseline.
//
// width, height, ha, and va specify the alignment of the text inside a box of width and height.
//
// cursor will draw a | character at the given boundary between two characters, -1 means no cursor.
//
// op sets the transform used to control placement of the text.
//
// It is OK to call Draw with a same text and a same face at every frame in terms of performance.
// Glyphs used for rendering are cached in least-recently-used way.
//
// Be careful that the passed font face is held by this package and is never released.
// This is a known issue (#498).
func DrawString(dst *ebiten.Image, text string, face font.Face, clr color.Color, underline bool, width, height int, ha Alignment, va Alignment, cursor int, op ebiten.DrawImageOptions) error {
	cr, cg, cb, ca := clr.RGBA()
	if ca == 0 {
		return nil
	}

	mBounds := glyphBounds(face, 'M')
	b := BoundString(face, text)
	w, h := float64(width), float64(height)
	ox, oy := float64(b.Min.X), fixed26_6ToFloat64(mBounds.Min.Y)
	tw, th := float64(b.Dx()), fixed26_6ToFloat64(mBounds.Max.Y)-fixed26_6ToFloat64(mBounds.Min.Y)

	var tx, ty float64

	switch ha {
	case Start:
		tx = -ox
	case Center:
		tx = -ox + w/2 - tw/2
	case End:
		tx = -ox + w - tw
	default:
		return fmt.Errorf("can't handle horizontal alignment %s", ha)
	}

	switch va {
	case Start:
		ty = -oy
	case Center:
		ty = -oy + h/2 - th/2
	case End:
		ty = -oy + h - th
	default:
		return fmt.Errorf("can't handle vertical alignment %s", va)
	}

	op.GeoM.Translate(tx, ty)
	op.ColorM.Scale(float64(cr)/float64(ca), float64(cg)/float64(ca), float64(cb)/float64(ca), float64(ca)/0xffff)

	var dx, dy fixed.Int26_6

	prevR := rune(-1)

	var cx fixed.Int26_6
	if cursor > -1 {
		cx = glyphAdvance(face, '|')
	}
	if cursor == 0 {
		drawGlyph(dst, face, '|', dx-cx/2, dy, &op)
		if len(text) == 0 {
			return nil
		}
	}

	for i, r := range text {
		if r == '\n' {
			continue
		}
		if prevR >= 0 {
			dx += face.Kern(prevR, r)
		}

		drawGlyph(dst, face, r, dx, dy, &op)

		if cursor == i {
			drawGlyph(dst, face, '|', dx-cx/2, dy, &op)
		}

		adv := glyphAdvance(face, r)
		if underline {
			x, y := op.GeoM.Apply(fixed26_6ToFloat64(dx), fixed26_6ToFloat64(dy))
			ebitenutil.DrawLine(dst, x, y+2, x+fixed26_6ToFloat64(adv), y+2, clr)
		}
		dx += adv

		prevR = r
	}

	if cursor >= len(text) {
		drawGlyph(dst, face, '|', dx-cx/2, dy, &op)
	}

	cleanCache(face)
	return nil
}

func BoundParagraph(face font.Face, text string, maxWidth int) image.Rectangle {
	m := face.Metrics()
	lineHeight := m.Height

	fx, fy, mw := fixed.I(0), fixed.I(0), fixed.I(maxWidth)
	prevR := rune(-1)

	var bounds fixed.Rectangle26_6
	for _, r := range text {
		width := glyphAdvance(face, r)

		if prevR >= 0 {
			fx += face.Kern(prevR, r)
		}
		if r == ' ' && mw > 0 && fx+width > mw {
			fx = 0
			fy += lineHeight
		}
		if r == '\n' {
			fx = 0
			fy += lineHeight
			prevR = rune(-1)
		}

		if r == ' ' {
			if fx+width > bounds.Max.X {
				bounds.Max.X = fx + width
			}
		} else {
			b := glyphBounds(face, r)
			b.Min.X += fx
			b.Max.X += fx
			b.Min.Y += fy
			b.Max.Y += fy
			bounds = bounds.Union(b)
		}

		fx += width
		prevR = r
	}

	return image.Rect(
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.X))),
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.Y))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.X))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.Y))),
	)
}

// DrawParagraph draws a paragraph of text.
// It allows configuration of the maximum width and height of the text,
// wrapping words to constrain the width and limiting the lines drawn to
// constrain the height.
//
// DrawParagraph returns the scroll position in [0, 1] if maxHeight and scroll are non-negative
//
// face is the font for text rendering.
//
// clr is the color for text rendering.
//
// underline draws a line along the baseline.
//
// maxWidth and maxHeight (both in pixels), when positive, turn on word wrapping and line limiting
//
// cursor will draw a | character at the given boundary between two characters, when negative no cursor will be drawn
//
// scroll sets the starting line and can be used with maxHeight to draw a scrollable window
//
// op sets the transform used to control placement of the text
//
// The '\n' newline character puts the following text on the next line.
// Line height is based on Metrics().Height of the font.
//
// It is OK to call DrawParagraph with a same text and a same face at every frame in terms of performance.
// Glyphs used for rendering are cached in least-recently-used way.
//
// Be careful that the passed font face is held by this package and is never released.
// This is a known issue (#498).
func DrawParagraph(dst *ebiten.Image, text string, face font.Face, clr color.Color, underline bool, maxWidth, maxHeight, cursor, scroll int, op ebiten.DrawImageOptions) float64 {
	cr, cg, cb, ca := clr.RGBA()
	if ca == 0 {
		return -1
	}

	mw, mh := fixed.I(maxWidth), fixed.I(maxHeight)

	op.ColorM.Scale(float64(cr)/float64(ca), float64(cg)/float64(ca), float64(cb)/float64(ca), float64(ca)/0xffff)

	var dx, dy fixed.Int26_6
	prevR := rune(-1)

	lineHeight := face.Metrics().Height
	mBounds := glyphBounds(face, 'M')
	op.GeoM.Translate(0, float64((mBounds.Max.Y - mBounds.Min.Y).Round()))

	var cx fixed.Int26_6
	if cursor > -1 {
		cx = glyphAdvance(face, '|')
	}
	if cursor == 0 {
		drawGlyph(dst, face, '|', dx-cx/2, dy, &op)
		if len(text) == 0 {
			return -1
		}
	}

	var i int
	var line int
	endLine := -1
	for _, r := range text {
		width := glyphAdvance(face, r)

		if prevR >= 0 {
			dx += face.Kern(prevR, r)
		}
		if r == ' ' && mw > 0 && dx+width > mw {
			dx = 0
			dy += lineHeight
			prevR = rune(-1)
			line++
		}
		if r == '\n' {
			dx = 0
			dy += lineHeight
			prevR = rune(-1)
			line++
			continue
		}

		if line >= scroll {
			oy := fixed.I(scroll).Mul(lineHeight)
			if scroll < 0 {
				oy = 0
			}
			if mh <= 0 || dy-oy+lineHeight <= mh || scroll < 0 {
				drawGlyph(dst, face, r, dx, dy-oy, &op)
				if cursor == i {
					drawGlyph(dst, face, '|', dx-cx/2, dy, &op)
				}
				endLine = line
			}
		}

		adv := glyphAdvance(face, r)
		if underline {
			x, y := op.GeoM.Apply(fixed26_6ToFloat64(dx), fixed26_6ToFloat64(dy))
			ebitenutil.DrawLine(dst, x, y+2, x+fixed26_6ToFloat64(adv), y+2, clr)
		}
		dx += adv
		i++

		prevR = r
	}

	if cursor >= len(text) {
		drawGlyph(dst, face, '|', dx-cx/2, dy, &op)
	}

	cleanCache(face)

	if height := dy.Ceil(); maxHeight > 0 && scroll >= 0 && height > maxHeight {
		if endLine >= line {
			return 1
		}
		return float64(scroll) / float64(line-(endLine-scroll))
	}
	return -1
}

// cacheSoftLimit indicates the soft limit of the number of glyphs in the cache.
// If the number of glyphs exceeds this soft limits, old glyphs are removed.
// Even after clearning up the cache, the number of glyphs might still exceeds the soft limit, but
// this is fine.
const cacheSoftLimit = 512

func cleanCache(face font.Face) {
	if len(glyphImageCache[face]) > cacheSoftLimit {
		for r, e := range glyphImageCache[face] {
			// 60 is an arbitrary number.
			if e.atime < time.Now().Unix()-60 {
				delete(glyphImageCache[face], r)
			}
		}
	}
}

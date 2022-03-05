package bento

var (
	Scrollspeed = 0.1
)

/*
type Scrollarea struct {
	states                                               [3][4]*NineSlice
	position                                             float64
	topBtnState, bottomBtnState, trackState, handleState ButtonState
	buffer                                               *ebiten.Image
}

func (s *Scrollarea) Update(keys []ebiten.Key, box *Box) ([]ebiten.Key, error) {
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	h := s.states[0][0].heights[0] + s.states[0][0].heights[1] + s.states[0][0].heights[2]

	r := box.innerRect()
	s.topBtnState = buttonState(r)
	if s.topBtnState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.position -= Scrollspeed
	}
	trackHeight := r.Dy() - 2*h
	s.trackState = buttonState(image.Rect(r.Min.X, r.Min.Y+h, r.Min.X+w, r.Min.Y+h+trackHeight))
	s.handleState = buttonState(image.Rect(r.Min.X, r.Min.Y+h+int(s.position*float64(trackHeight-h)), r.Min.X+w, r.Min.Y+h+int(s.position*float64(trackHeight-h))+h))
	s.bottomBtnState = buttonState(image.Rect(r.Min.X, r.Max.Y-h, r.Min.X+w, r.Max.X))
	if s.bottomBtnState == Active && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.position += Scrollspeed
	}

	_, dy := ebiten.Wheel()
	s.position += dy * Scrollspeed

	// TODO: bad position can crash during Draw
	if s.position < 0 {
		s.position = 0
	}
	if s.position > 1 {
		s.position = 1
	}
	return keys, nil
}

func (s *Scrollarea) Draw(screen *ebiten.Image, n *Box) {

	n.buffer = ebiten.NewImage(n.TextBounds.Dx(), n.TextBounds.Dy())
	text.DrawParagraph(n.buffer, n.templateContent(), n.style.Font, n.style.Color, 0, 0, n.style.MaxWidth, -n.TextBounds.Min.Y)
	offset := int(scroll.position * float64(n.TextBounds.Dy()))
	cropped := ebiten.NewImageFromImage(n.buffer.SubImage(image.Rect(0, offset, content.Dx(), content.Dy()+offset)))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(content.Min.X), float64(content.Min.Y))
	img.DrawImage(cropped, op)
	n.style.Scrollbar.Draw(img, n)

	r := n.innerRect()
	w := s.states[0][0].widths[0] + s.states[0][0].widths[1] + s.states[0][0].widths[2]
	h := s.states[0][0].heights[0] + s.states[0][0].heights[1] + s.states[0][0].heights[2]
	trackHeight := r.Dy() - 2*h
	s.states[s.topBtnState][0].Draw(screen, r.Min.X, r.Min.Y, w, h)
	s.states[s.trackState][1].Draw(screen, r.Min.X, r.Min.Y+h, w, trackHeight)
	s.states[s.handleState][2].Draw(screen, r.Min.X, r.Min.Y+h+int(s.position*float64(trackHeight-h)), w, h)
	s.states[s.bottomBtnState][3].Draw(screen, r.Min.X, r.Max.Y-h, w, h)
}
*/

package bento

import "fmt"

type Justification int

const (
	Start = Justification(iota)
	End
	Center
	Evenly
	Stretch
)

func ParseJustification(s string) Justification {
	switch s {
	case "start":
		return Start
	case "end":
		return End
	case "center":
		return Center
	case "evenly":
		return Evenly
	case "stretch":
		return Stretch
	default:
		panic(fmt.Errorf("invalid alignment %s", s))
	}
}

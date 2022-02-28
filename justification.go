package bento

import "fmt"

type Justification int

const (
	Start = Justification(iota)
	End
	Center
	Between
	Around
	Evenly
)

func ParseJustification(s string) Justification {
	switch s {
	case "start":
		return Start
	case "end":
		return End
	case "center":
		return Center
	case "between":
		return Between
	case "around":
		return Around
	case "evenly":
		return Evenly
	default:
		panic(fmt.Errorf("invalid alignment %s", s))
	}
}

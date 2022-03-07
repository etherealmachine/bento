package text

import (
	"testing"
)

func TestBoundString(t *testing.T) {
	b1 := BoundString(fonts["NotoSans"].load(16), "Hello World")
	b2 := BoundString(fonts["NotoSans"].load(16), "Hello World ")
	if b1.Dx() == b2.Dx() {
		t.Errorf("incorrect calculation of space")
	}
}

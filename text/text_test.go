package text

import (
	"log"
	"testing"
)

func TestBoundString(t *testing.T) {
	b1 := BoundString(fonts["NotoSans"].load(16), "Hello World")
	b2 := BoundString(fonts["NotoSans"].load(16), "Hello World ")
	log.Println(b1.Dx(), b2.Dx())
	if b1.Dx() == b2.Dx() {
		t.Errorf("bound calculation did not include trailing space")
	}
}

package main

import (
	"flag"
	_ "image/png"
	"log"

	"github.com/etherealmachine/bento/v0"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Demo struct {
	ui bento.Box
}

func (d *Demo) Update() error {
	var keys []ebiten.Key
	keys = inpututil.AppendPressedKeys(keys)
	if _, err := d.ui.Update(keys); err != nil {
		return err
	}
	return nil
}

func (d *Demo) Draw(screen *ebiten.Image) {
	screen.Clear()
	d.ui.Draw(screen)
}

func (d *Demo) Layout(ow, oh int) (int, int) {
	return ow, oh
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ebiten.SetWindowSize(1024, 800)
	ebiten.SetWindowTitle("Bento Demo")
	ui, err := bento.Build(&TextDemo{})
	if err != nil {
		log.Fatal(err)
	}
	d := &Demo{ui: ui}
	if err := ebiten.RunGame(d); err != nil {
		log.Fatal(err)
	}
}

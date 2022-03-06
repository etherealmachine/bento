package main

import (
	"flag"
	_ "image/png"
	"io/ioutil"
	"log"
	"strings"

	"github.com/etherealmachine/bento"
	"github.com/hajimehoshi/ebiten/v2"
)

var paragraphs []string

func init() {
	bs, err := ioutil.ReadFile("loomings.txt")
	if err != nil {
		log.Fatal(err)
	}
	paragraphs = strings.Split(string(bs), "\n")
}

type Game struct {
	ui *bento.Box
}

func (g *Game) Update() error {
	if err := g.ui.Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.ui.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ow, oh
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ebiten.SetWindowSize(1280, 900)
	ebiten.SetWindowTitle("Bento Demo")
	ui, err := bento.Build(&Demo{})
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(&Game{ui: ui}); err != nil {
		log.Fatal(err)
	}
}

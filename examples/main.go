package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/etherealmachine/bento"
	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
	_ "image/png"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

//go:embed loomings.txt
var loomings string
var paragraphs []string

func init() {
	paragraphs = strings.Split(loomings, "\n")
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
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Bento Demo")
	ui, err := bento.Build(&Demo{})
	if err != nil {
		log.Fatal(err)
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if err := ebiten.RunGame(&Game{ui: ui}); err != nil {
		log.Fatal(err)
	}
}

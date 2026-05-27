package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

func (gs *GameState) Cycle() {
	*gs += 1
	if *gs > BarBar {
		*gs = HelloWorld
	}
}

const (
	_ GameState = iota
	HelloWorld
	FooFoo
	BarBar
)

type Game struct {
	keys  []ebiten.Key
	state GameState
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state.Cycle()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == HelloWorld {
		ebitenutil.DebugPrint(screen, "Hello, World!")
	} else if g.state == FooFoo {
		ebitenutil.DebugPrint(screen, "FooFoo!")
	} else if g.state == BarBar {
		ebitenutil.DebugPrint(screen, "BarBar!")
	} else {
		ebitenutil.DebugPrint(screen, "Unknown GameState is found")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{state: HelloWorld}); err != nil {
		log.Fatal(err)
	}
}

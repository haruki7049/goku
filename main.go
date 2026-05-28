package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

func (gs *GameState) Cycle() {
	*gs += 1
	if *gs > GetAdvice {
		*gs = HelloWorld
	}
}

const (
	_ GameState = iota
	HelloWorld
	GetAdvice
)

type Response struct {
	Slip Slip `json:"slip"`
}

type Slip struct {
	ID     int    `json:"id"`
	Advice string `json:"advice"`
}

type Game struct {
	message string
	state   GameState
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state.Cycle()

		if g.state == GetAdvice {
			g.message = "Loading..."
			go g.fetchMessage()
		}
	}
	return nil
}

func (g *Game) fetchMessage() error {
	resp, err := http.Get("https://api.adviceslip.com/advice")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	g.message = response.Slip.Advice
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case HelloWorld:
		ebitenutil.DebugPrint(screen, "Hello, World!\n"+"You can toggle screen mode by Space Key between\nGetAdvice and HelloWorld.")
	case GetAdvice:
		ebitenutil.DebugPrint(screen, g.message+"\n"+"You can toggle screen mode by Space Key between\nGetAdvice and HelloWorld.")
	default:
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

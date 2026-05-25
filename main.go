package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Chart represents a beatmap containing a list of notes.
type Chart struct {
	Notes []NoteData `json:"notes"`
}

// NoteData represents a single note's timing and lane.
type NoteData struct {
	Time int `json:"time"` // Time in milliseconds
	Lane int `json:"lane"` // Lane index (0 to 3)
}

// Note represents an active note in the game.
type Note struct {
	Data NoteData
	Hit  bool
}

// Game represents the main game state.
type Game struct {
	StartTime   time.Time
	Notes       []*Note
	CurrentTime int
	Score       int
	Combo       int
	Keys        []ebiten.Key
}

// NewGame initializes the game and loads the chart.
func NewGame() *Game {
	g := &Game{
		StartTime: time.Now(),
		Keys:      []ebiten.Key{ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK},
	}
	g.loadChart("chart.json")
	return g
}

// loadChart loads notes from a JSON file or uses dummy data if not found.
func (g *Game) loadChart(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// Fallback dummy chart data
		data = []byte(`{"notes": [{"time": 1000, "lane": 0}, {"time": 1500, "lane": 1}, {"time": 2000, "lane": 2}, {"time": 2500, "lane": 3}, {"time": 3000, "lane": 0}]}`)
	}

	var chart Chart
	if err := json.Unmarshal(data, &chart); err != nil {
		log.Fatal(err)
	}

	for _, nd := range chart.Notes {
		g.Notes = append(g.Notes, &Note{Data: nd})
	}
}

// Update handles game logic and input.
func (g *Game) Update() error {
	// Increment time (assuming 60 FPS)
	g.CurrentTime = int(time.Since(g.StartTime).Milliseconds())

	// Handle input
	for i, key := range g.Keys {
		if inpututil.IsKeyJustPressed(key) {
			g.checkHit(i)
		}
	}

	// Check for missed notes
	for _, n := range g.Notes {
		if !n.Hit && g.CurrentTime-n.Data.Time > 200 { // Miss threshold (200ms)
			n.Hit = true
			g.Combo = 0
		}
	}

	return nil
}

// checkHit evaluates if a note is hit in the specified lane.
func (g *Game) checkHit(lane int) {
	for _, n := range g.Notes {
		if !n.Hit && n.Data.Lane == lane {
			diff := n.Data.Time - g.CurrentTime
			if diff < 0 {
				diff = -diff
			}

			if diff <= 100 { // Hit threshold (100ms)
				n.Hit = true
				g.Score += 100
				g.Combo++
				return
			}
		}
	}
}

// Draw renders the game objects on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw lanes
	for i := 0; i < 4; i++ {
		ebitenutil.DrawRect(screen, float64(100+i*60), 0, 50, 480, color.RGBA{50, 50, 50, 255})
	}

	// Draw hit line
	ebitenutil.DrawRect(screen, 100, 400, 230, 5, color.RGBA{255, 0, 0, 255})

	// Draw notes
	for _, n := range g.Notes {
		if !n.Hit {
			// Calculate vertical position (takes 2000ms to reach the hit line)
			y := 400.0 - float64(n.Data.Time-g.CurrentTime)*(400.0/2000.0)
			if y > -50 && y < 480 {
				ebitenutil.DrawRect(screen, float64(100+n.Data.Lane*60), y, 50, 20, color.RGBA{0, 255, 0, 255})
			}
		}
	}

	// Draw UI text
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d\nCombo: %d", g.Score, g.Combo))
}

// Layout determines the window resolution.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Rhythm Game")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

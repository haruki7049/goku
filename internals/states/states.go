package internals

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

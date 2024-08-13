package game

import (
	"fmt"
	"letters/solver"
)

type Game struct {
	g map[int]*solver.Solver
}

func NewGame() *Game {
	return &Game{
		g: make(map[int]*solver.Solver),
	}
}

func (g *Game) AddPlayer(player int, s *solver.Solver) {
	g.g[player] = s
}

func (g *Game) AddNotContains(player int, runes []rune) error {
	if _, ok := g.g[player]; !ok {
		return fmt.Errorf("player %d not found. Please /start game first", player)
	}
	g.g[player].NotContain(runes)
	return nil
}
func (g *Game) AddContains(player int, runes []rune) error {
	if _, ok := g.g[player]; !ok {
		return fmt.Errorf("player %d not found. Please /start game first", player)
	}
	g.g[player].Contains(runes)
	return nil
}

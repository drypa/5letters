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

func (g *Game) AddCorrectPosition(player int, r rune, pos int) error {
	if _, ok := g.g[player]; !ok {
		return fmt.Errorf("player %d not found. Please /start game first", player)
	}
	places := make([]solver.RunePlace, 1)
	places[0] = solver.RunePlace{Rune: r, Pos: pos}
	g.g[player].CorrectRunePlaces(places)
	return nil
}
func (g *Game) AddIncorrectPosition(player int, r rune, pos int) error {
	if _, ok := g.g[player]; !ok {
		return fmt.Errorf("player %d not found. Please /start game first", player)
	}
	places := make([]solver.RunePlace, 1)
	places[0] = solver.RunePlace{Rune: r, Pos: pos}
	g.g[player].IncorrectRunePlaces(places)
	return nil
}

func (g *Game) GetResult(player int) ([]string, error) {
	if _, ok := g.g[player]; !ok {
		return nil, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	return g.g[player].GetSuitable(), nil
}

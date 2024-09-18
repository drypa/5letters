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

func (g *Game) AddNotContains(player int, runes []rune) (int, error) {
	if _, ok := g.g[player]; !ok {
		return -1, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	g.g[player].NotContain(runes)
	return g.getCount(player)
}
func (g *Game) AddContains(player int, runes []rune) (int, error) {
	if _, ok := g.g[player]; !ok {
		return -1, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	g.g[player].Contains(runes)
	return g.getCount(player)
}

func (g *Game) AddCorrectPosition(player int, r rune, pos int) (int, error) {
	if _, ok := g.g[player]; !ok {
		return -1, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	places := make([]solver.RunePlace, 1)
	places[0] = solver.RunePlace{Rune: r, Pos: pos}
	g.g[player].CorrectRunePlaces(places)
	return g.getCount(player)
}
func (g *Game) AddIncorrectPosition(player int, r rune, pos int) (int, error) {
	if _, ok := g.g[player]; !ok {
		return -1, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	places := make([]solver.RunePlace, 1)
	places[0] = solver.RunePlace{Rune: r, Pos: pos}
	g.g[player].IncorrectRunePlaces(places)
	return g.getCount(player)
}

func (g *Game) GetResult(player int) ([]string, error) {
	if _, ok := g.g[player]; !ok {
		return nil, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	return g.g[player].GetSuitable(), nil
}
func (g *Game) getCount(player int) (int, error) {
	if _, ok := g.g[player]; !ok {
		return -1, fmt.Errorf("player %d not found. Please /start game first", player)
	}
	return len(g.g[player].GetSuitable()), nil
}

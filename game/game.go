package game

import (
	"fmt"
	"image/color"
	"my-game/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	player           *Player
	lasers           []*Laser
	meteors          []*Meteor
	meteorSpawnTimer *Timer
	score            int
}

func NewGame() *Game {
	g := &Game{
		meteorSpawnTimer: NewTimer(24),
	}
	player := NewPlayer(g)
	g.player = player
	return g
}

func (g *Game) Update() error {
	g.player.Update()

	for i, l := range g.lasers {
		l.Update()
		if l.position.Y < 0 {
			g.lasers = append(g.lasers[:i], g.lasers[i+1:]...)
		}
	}

	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.isReady() {
		g.meteorSpawnTimer.Reset()
		m := NewMeteor()
		g.meteors = append(g.meteors, m)
	}

	for i, m := range g.meteors {
		m.Update()

		if m.Collider().Intersects(g.player.Collider()) {
			g.Reset()
			break
		}

		if m.position.Y > screenHeight {
			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
		}
	}

	for i, m := range g.meteors {
		for j, l := range g.lasers {
			if m.Collider().Intersects(l.Collider()) {
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
				g.lasers = append(g.lasers[:j], g.lasers[j+1:]...)
				g.score++
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, l := range g.lasers {
		l.Draw(screen)
	}

	for _, m := range g.meteors {
		m.Draw(screen)
	}

	text.Draw(screen, fmt.Sprintf("Pontos: %d", g.score), assets.FontUi, 20, 100, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) AddLasers(laser *Laser) {
	g.lasers = append(g.lasers, laser)
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.meteors = nil
	g.lasers = nil
	g.score = 0
	g.meteorSpawnTimer.Reset()
}

package game

import (
	"my-game/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	image             *ebiten.Image
	position          Vector
	game              *Game
	laserLoadingTimer *Timer
}

func NewPlayer(game *Game) *Player {
	image := assets.PlayerSprite
	bounds := image.Bounds()
	halfW := float64(bounds.Dx()) / 2

	position := Vector{
		X: (screenWidth / 2) - halfW,
		Y: (screenHeight) - (float64(bounds.Dy()) + 20),
	}

	return &Player{
		image:             image,
		position:          position,
		game:              game,
		laserLoadingTimer: NewTimer(12),
	}
}

func (p *Player) Update() {
	speed := 6.0
	positionX := p.position.X

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		positionX -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		positionX += speed
	}

	p.laserLoadingTimer.Update()

	bounds := p.image.Bounds()
	boundX := bounds.Dx()

	if positionX <= 0 {
		positionX = 0
	} else if positionX >= float64(screenWidth-boundX) {
		positionX = float64(screenWidth - boundX)
	}

	p.position.X = positionX

	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.laserLoadingTimer.isReady() {
		p.laserLoadingTimer.Reset()
		halfW := boundX / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + float64(halfW),
			p.position.Y - halfH/2,
		}

		laser := NewLaser(spawnPos)

		p.game.AddLasers(laser)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.image, op)
}

func (p *Player) Collider() Rect {
	bounds := p.image.Bounds()

	return NewRect(p.position.X, p.position.Y, float64(bounds.Dx()), float64(bounds.Dy()))
}

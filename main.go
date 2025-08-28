package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	ebitenVector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quartercastle/vector"
	"image/color"
	"log"
	"math/rand"
)

const (
	Width  = 1280
	Height = 720
)

type Snake struct {
	Parts  []vector.Vector
	Radius int
	Color  color.Color
}

func NewSnake(initVector vector.Vector, length, radius int) *Snake {
	result := Snake{
		Radius: radius,
		Color: color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255}}
	for i := 0; i < length; i++ {
		result.Parts = append(result.Parts, vector.Vector{initVector.X(), initVector.Y() + float64(i*radius)})
	}
	return &result
}

func (s *Snake) DoStep(direction vector.Vector) {
	normalizeDirection := direction.Unit()
	newPart := s.Parts[0].Add(normalizeDirection.Scale(float64(s.Radius)))
	s.Parts = append([]vector.Vector{newPart}, s.Parts[:len(s.Parts)-1]...)
}

type Game struct {
	Dots      []vector.Vector
	CurSnake  *Snake
	GameOver  bool
	ArenaSize int
	WindowW   int
	WindowH   int
}

func NewGame(arenaSize, dotsNum, initSnakeLen, initSnakeRadius, w, h int) *Game {
	result := &Game{
		WindowW:   w,
		WindowH:   h,
		ArenaSize: arenaSize,
		CurSnake:  NewSnake(vector.Vector{}, initSnakeLen, initSnakeRadius),
	}
	for i := 0; i < dotsNum; i++ {
		result.Dots = append(result.Dots, vector.Vector{float64(rand.Intn(arenaSize)), float64(rand.Intn(arenaSize))})
	}
	return result
}

func (g *Game) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorVector := vector.Vector{float64(cursorX - g.WindowW/2), float64(cursorY - g.WindowH/2)}
	g.CurSnake.DoStep(cursorVector)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Фон
	newScreen := ebiten.NewImageFromImage(screen)
	newScreen.Fill(color.Black)

	//Камера
	//leftTopCorner := vector.Vector{
	//	g.CurSnake.Parts[0].X() - float64(g.WindowW/2),
	//	g.CurSnake.Parts[0].Y() - float64(g.WindowH/2)}

	//Змея
	for _, part := range g.CurSnake.Parts {
		ebitenVector.StrokeCircle(newScreen, float32(part.X()), float32(part.Y()),
			float32(g.CurSnake.Radius), 1, g.CurSnake.Color, false)
	}

	//Точки
	for _, dot := range g.Dots {
		ebitenVector.DrawFilledCircle(newScreen, float32(dot.X()), float32(dot.Y()), 2, color.White, false)
	}

	screen.DrawImage(newScreen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("SlitherioClone")
	game := NewGame(3000, 100, 10, 10, Width, Height)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

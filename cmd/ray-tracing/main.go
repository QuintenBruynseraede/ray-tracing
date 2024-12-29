package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/quintenbruynseraede/ray-tracing/internal"
)

const (
	ASPECT_RATIO = 16.0 / 9.0
	IMAGE_WIDTH  = 800
	IMAGE_HEIGHT = int(float64(IMAGE_WIDTH) / ASPECT_RATIO)
)

func main() {
	ebiten.SetWindowSize(IMAGE_WIDTH, IMAGE_HEIGHT)

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", g.Image.At(12, 100))
	fmt.Printf("%v", g.Image.At(100, 200))
	fmt.Printf("%v", g.Image.At(200, 300))
	fmt.Printf("%v", g.Image.At(300, 350))
	internal.SaveScreenshot(g.Image)
}

type Game struct {
	Image  *image.RGBA
	World  *internal.HittableList
	Camera internal.Camera
}

func NewGame() *Game {
	camera := internal.NewCamera(internal.NewVec3(0, 0, 0), IMAGE_WIDTH, IMAGE_HEIGHT)
	world := internal.NewHittableList(
		internal.NewSphere(internal.NewVec3(0, 0, -1), 0.5),
		internal.NewSphere(internal.NewVec3(0, -100.5, -1), 100),
	)

	return &Game{
		Image:  image.NewRGBA(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)),
		Camera: camera,
		World:  world,
	}
}

func (g *Game) Update() error {
	g.Camera.Center = g.Camera.Center.Add(internal.NewVec3(0, 0, 0.01))
	g.Image = g.Camera.Render(g.Image, g.World)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.Image.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return IMAGE_WIDTH, IMAGE_HEIGHT
}

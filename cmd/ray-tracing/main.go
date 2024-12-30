package main

import (
	"fmt"
	"image"
	"image/color"
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
	ebiten.SetTPS(0)

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	internal.SaveScreenshot(g.Image)
}

type Game struct {
	Image  *image.RGBA
	World  *internal.HittableList
	Camera internal.Camera
}

func NewGame() *Game {
	camera := internal.NewCamera(internal.NewVec3(0, 0, 0), IMAGE_WIDTH, IMAGE_HEIGHT)

	material_ground := internal.Lambertian{Albedo: color.RGBA{204, 204, 0, 255}}
	material_center := internal.Lambertian{Albedo: color.RGBA{25, 50, 128, 255}}
	material_left := internal.Dielectric{RefractionIndex: 1.5}
	material_bubble := internal.Dielectric{RefractionIndex: 1.0 / 1.5}
	material_right := internal.Metal{
		Albedo: color.RGBA{205, 153, 50, 255},
		Fuzz:   1.0,
	}

	world := internal.NewHittableList(
		internal.NewSphere(internal.NewVec3(0, -100.5, -1.0), 100, material_ground),
		internal.NewSphere(internal.NewVec3(0, 0, -1.2), 0.5, material_center),
		internal.NewSphere(internal.NewVec3(-1, 0, -1), 0.5, material_left),
		internal.NewSphere(internal.NewVec3(-1, 0, -1), 0.4, material_bubble),
		internal.NewSphere(internal.NewVec3(1, 0, -1), 0.5, material_right),
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

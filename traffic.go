package main

import (
	"image/color"
	"traffic/systems"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	KeyboardScrollSpeed = 400
	EdgeScrollSpeed     = KeyboardScrollSpeed
	EdgeWidth           = 20
	ZoomSpeed           = -0.125
	WorldWidth          = 400
	WorldHeight         = 400
)

type myScene struct{}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk, to allow you to register/queue them
func (*myScene) Preload() {
	engo.Files.Load("textures/city.png")
}

// Setup is called before the main loop starts. It allows you to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	world.AddSystem(common.NewKeyboardScroller(
		KeyboardScrollSpeed, engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis))
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: EdgeScrollSpeed, EdgeMargin: EdgeWidth})
	world.AddSystem(&common.MouseZoomer{ZoomSpeed: ZoomSpeed})

	world.AddSystem(&systems.CityBuildingSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:          "Traffic Manager",
		Width:          WorldWidth,
		Height:         WorldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}

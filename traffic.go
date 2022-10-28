package main

import (
	"image"
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
	// todo needs more intelligence
	HUDHeight = WorldHeight / 4
	HUDWidth  = WorldWidth / 4
)

type myScene struct{}

type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

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

	hud := HUD{BasicEntity: ecs.NewBasic()}
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - HUDHeight},
		Width:    HUDWidth,
		Height:   HUDHeight,
	}
	hudImage := image.NewUniform(color.RGBA{205, 205, 205, 255})
	hudNRGBA := common.ImageToNRGBA(hudImage, 200, 200)
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)

	hud.RenderComponent = common.RenderComponent{
		Drawable: hudTexture,
		Scale:    engo.Point{X: 1, Y: 1},
		Repeat:   common.Repeat,
	}
	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(1)

	world.AddSystem(common.NewKeyboardScroller(
		KeyboardScrollSpeed, engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis))
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: EdgeScrollSpeed, EdgeMargin: EdgeWidth})
	world.AddSystem(&common.MouseZoomer{ZoomSpeed: ZoomSpeed})

	world.AddSystem(&systems.CityBuildingSystem{})

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}
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

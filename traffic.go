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
	EdgeScrollSpeed     = KeyboardScrollSpeed
	EdgeWidth           = 20
	ZoomSpeed           = -0.125
	WorldWidth          = 800
	WorldHeight         = 800
	KeyboardScrollSpeed = (WorldWidth * WorldHeight) / 600
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

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk, to allow you to register/queue them
func (*myScene) Preload() {
	engo.Files.Load("textures/city.png", "tilemap/TrafficMap.tmx")
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

	resource, err := engo.Files.Resource("tilemap/TrafficMap.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	tiles := make([]*Tile, 0)
	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement.Image,
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
				}
				tiles = append(tiles, tile)
			}
		}
	}
	common.CameraBounds = levelData.Bounds()
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			// add HUD to the RenderSystem
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
			// add tiles to the RenderSystem
			for _, tile := range tiles {
				sys.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
			}

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

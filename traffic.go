package main

import (
	"bytes"
	"image"
	"image/color"
	"traffic/globals"
	"traffic/systems"

	"golang.org/x/image/font/gofont/gosmallcaps"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type myScene struct{}

// HUD defines a HUD entity
type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Tile defines a Tile entity
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk, to allow you to register/queue them
func (*myScene) Preload() {
	engo.Files.Load("textures/citySheet.png", "tilemap/TrafficMap.tmx")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
}

// Setup is called before the main loop starts. It allows you to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	world.AddSystem(common.NewKeyboardScroller(
		globals.KeyboardScrollSpeed, engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis))
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: globals.EdgeScrollSpeed, EdgeMargin: globals.EdgeWidth})
	world.AddSystem(&common.MouseZoomer{ZoomSpeed: globals.ZoomSpeed})

	world.AddSystem(&systems.CityBuildingSystem{})
	world.AddSystem(&systems.HUDTextSystem{})
	world.AddSystem(&systems.MoneySystem{})

	hud := HUD{BasicEntity: ecs.NewBasic()}
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - globals.HUDHeight},
		Width:    globals.HUDWidth,
		Height:   globals.HUDHeight,
	}
	hudImage := image.NewUniform(color.RGBA{205, 205, 205, 255})
	hudNRGBA := common.ImageToNRGBA(hudImage, 200, 200)
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)

	hud.RenderComponent = common.RenderComponent{
		Repeat:   common.Repeat,
		Drawable: hudTexture,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(globals.HUDZ)

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

	// add everything to the RenderSystem
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
		Width:          globals.WorldWidth,
		Height:         globals.WorldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}

package systems

import (
	"fmt"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type City struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type CityBuildingSystem struct {
	world *ecs.World

	mouseTracker MouseTracker
}

// Remove is called whenever an Entity is removed from the World, in order to remove it from this sytem as well
func (*CityBuildingSystem) Remove(ecs.BasicEntity) {}

func (cb *CityBuildingSystem) New(w *ecs.World) {
	cb.world = w
	fmt.Println("CityBuildingSystem was added to the scene")

	cb.mouseTracker.BasicEntity = ecs.NewBasic()
	cb.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&cb.mouseTracker.BasicEntity, &cb.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

// Update is ran every frame, with `dt` being the time in seconds since the last frame
func (cb *CityBuildingSystem) Update(dt float32) {
	if engo.Input.Button("AddCity").JustPressed() {
		fmt.Println("Player pressed F1")

		city := City{BasicEntity: ecs.NewBasic()}
		city.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X: cb.mouseTracker.MouseX,
				Y: cb.mouseTracker.MouseY},
			Width:  30,
			Height: 64,
		}

		texture, err := common.LoadedSprite("textures/city.png")
		if err != nil {
			log.Println("Unable to load texture: " + err.Error())
		}

		city.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{X: 0.1, Y: 0.1},
		}

		for _, system := range cb.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&city.BasicEntity, &city.RenderComponent, &city.SpaceComponent)
			}
		}
	}
}

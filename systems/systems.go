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
	world        *ecs.World
	mouseTracker MouseTracker
}

func (cb *CityBuildingSystem) Remove(e ecs.BasicEntity) {}

func (cb *CityBuildingSystem) Update(dt float32) {
	if engo.Input.Button("AddCity").JustPressed() {
		fmt.Println("The gamer pressed F1")

		city := createCity(cb.mouseTracker.MouseX, cb.mouseTracker.MouseY)

		for _, system := range cb.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&city.BasicEntity, &city.RenderComponent, &city.SpaceComponent)
			}
		}

	}
}

func createCity(x, y float32) City {
	city := City{BasicEntity: ecs.NewBasic()}
	city.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: x, Y: y},
		Width:    30,
		Height:   64,
	}
	texture, err := common.LoadedSprite("textures/city.png")
	if err != nil {
		log.Printf("Unable to load texture: %v\n", err)
	}
	city.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 0.1, Y: 0.1},
	}
	return city
}

func (cb *CityBuildingSystem) New(world *ecs.World) {
	fmt.Println("CityBuildingSystem was added to the Scene")

	cb.world = world
	cb.mouseTracker.BasicEntity = ecs.NewBasic()
	cb.mouseTracker.MouseComponent = common.MouseComponent{
		Track: true,
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&cb.mouseTracker.BasicEntity, &cb.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

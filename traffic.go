package main

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type City struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type myScene struct{}

func (s *myScene) Type() string { return "myGame" }

func (s *myScene) Preload() {
	log.Println(engo.Files.GetRoot())
	engo.Files.Load("textures/city.png")
}

func (s *myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})

	city := City{BasicEntity: ecs.NewBasic()}
	city.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 10, Y: 10},
		Width:    303,
		Height:   641,
	}

	texture, err := common.LoadedSprite("textures/city.png")
	if err != nil {
		log.Printf("Unable to load texture: %v", err)
	}
	city.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 1, Y: 1},
	}

	common.SetBackground(color.White)

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&city.BasicEntity, &city.RenderComponent, &city.SpaceComponent)
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}
	engo.Run(opts, &myScene{})
}

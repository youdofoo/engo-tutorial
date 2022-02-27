package main

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/youdofoo/engo-tutorial/systems"
)

type myScene struct{}

func (s *myScene) Type() string { return "myGame" }

func (s *myScene) Preload() {
	log.Println(engo.Files.GetRoot())
	engo.Files.Load("textures/city.png")
}

func (s *myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	kbs := common.NewKeyboardScroller(400, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: 400, EdgeMargin: 20})
	world.AddSystem(&common.MouseZoomer{ZoomSpeed: -0.125})

	world.AddSystem(&systems.CityBuildingSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          400,
		Height:         400,
		StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}

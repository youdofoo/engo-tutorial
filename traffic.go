package main

import (
	"image"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/youdofoo/engo-tutorial/systems"
)

type HUD struct {
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
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})

	kbs := common.NewKeyboardScroller(400, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: 400, EdgeMargin: 20})
	world.AddSystem(&common.MouseZoomer{ZoomSpeed: -0.125})

	hud := HUD{BasicEntity: ecs.NewBasic()}
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: engo.WindowHeight() - 200},
		Width:    200,
		Height:   200,
	}
	hudImage := image.NewUniform(color.RGBA{R: 205, G: 205, B: 205, A: 255})
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
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}

	world.AddSystem(&systems.CityBuildingSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          800,
		Height:         800,
		StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}

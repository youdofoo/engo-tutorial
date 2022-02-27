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

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type myScene struct{}

func (s *myScene) Type() string { return "myGame" }

func (s *myScene) Preload() {
	log.Println(engo.Files.GetRoot())
	engo.Files.Load("textures/city.png", "tilemap/TrafficMap.tmx")
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

	resource, err := engo.Files.Resource("tilemap/TrafficMap.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level
	tiles := make([]*Tile, 0)
	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
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
	common.CameraBounds = levelData.Bounds()
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
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

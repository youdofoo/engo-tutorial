package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var Spritesheet *common.Spritesheet
var cities = [...][12]int{
	{
		99, 100, 101,
		454, 269, 455,
		415, 195, 416,
		452, 306, 453,
	},
	{
		99, 100, 101,
		268, 269, 270,
		268, 269, 270,
		305, 306, 307,
	},
	{
		75, 76, 77,
		446, 261, 447,
		446, 261, 447,
		444, 298, 445,
	},
	{
		75, 76, 77,
		407, 187, 408,
		407, 187, 408,
		444, 298, 445,
	},
	{
		75, 76, 77,
		186, 150, 188,
		186, 150, 188,
		297, 191, 299,
	},
	{
		83, 84, 85,
		413, 228, 414,
		411, 191, 412,
		448, 302, 449,
	},
	{
		83, 84, 85,
		227, 228, 229,
		190, 191, 192,
		301, 302, 303,
	},
	{
		91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 946, 947,
	},
	{
		91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 803, 947,
	},
	{
		91, 92, 93,
		238, 239, 240,
		238, 239, 240,
		312, 313, 314,
	},
}

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

	usedTiles []int
	elapsed   float32
	buildTime float32
	built     int
}

func (cb *CityBuildingSystem) Remove(e ecs.BasicEntity) {}

func (cb *CityBuildingSystem) Update(dt float32) {
	cb.elapsed += dt
	if cb.elapsed >= cb.buildTime {
		cb.generateCity()
		cb.built++
		cb.elapsed = 0
		cb.updateBuildTime()
	}
}

func (cb *CityBuildingSystem) generateCity() {
	x := rand.Intn(18)
	y := rand.Intn(18)
	t := x + y*18

	for cb.isTileUsed(t) {
		if len(cb.usedTiles) > 300 {
			break
		}
		x = rand.Intn(18)
		y = rand.Intn(18)
		t = x * y * 18
	}
	cb.usedTiles = append(cb.usedTiles, t)

	city := rand.Intn(len(cities))
	cityTiles := make([]*City, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			xx := float32((x+1)*64 + 8 + i*16)
			yy := float32((y+1)*64 + 8 + j*16)
			tile := createCityTile(xx, yy, Spritesheet.Cell(cities[city][i+j*3]))
			cityTiles = append(cityTiles, tile)
		}
	}

	for _, system := range cb.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range cityTiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}

func createCityTile(x, y float32, d common.Drawable) *City {
	city := &City{BasicEntity: ecs.NewBasic()}
	city.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: x, Y: y},
	}
	city.RenderComponent = common.RenderComponent{
		Drawable: d,
	}
	city.RenderComponent.SetZIndex(1)
	return city
}

func (cb *CityBuildingSystem) isTileUsed(tile int) bool {
	for _, t := range cb.usedTiles {
		if tile == t {
			return true
		}
	}
	return false
}

func (cb *CityBuildingSystem) updateBuildTime() {
	switch {
	case cb.built < 2:
		cb.buildTime = 10 + 5*rand.Float32()
	case cb.built < 5:
		cb.buildTime = 60 + 30*rand.Float32()
	case cb.built < 10:
		cb.buildTime = 30 + 60*rand.Float32()
	case cb.built < 20:
		cb.buildTime = 30 + 35*rand.Float32()
	case cb.built < 25:
		cb.buildTime = 30 + 30*rand.Float32()
	default:
		cb.buildTime = 20 + 20*rand.Float32()
	}
}

func (cb *CityBuildingSystem) New(world *ecs.World) {
	fmt.Println("CityBuildingSystem was added to the Scene")
	rand.Seed(time.Now().UnixNano())

	Spritesheet = common.NewSpritesheetWithBorderFromFile("textures/citySheet.png", 16, 16, 1, 1)

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

package graphUI

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/iniconfig"
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/internal/common"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/world/constant"

)

type (
	GraphUI struct {
		actor.Base
		Config iniconfig.Config
		l      sync.Mutex

		entities map[string]*Entity
		tiles    map[int64]*Tile
		grids    []*inner.Vec3F
		lastEID  string
		Path     []*inner.Vec3F
		clear    int
	}

	Tile struct {
		ExclusiveBound *inner.Bound
		ActualBound    *inner.Bound
		Count          int64
	}

	Entity struct {
		EID       string
		X         float64
		Y         float64
		Duplicate bool
		HP        int64

		Path []*inner.Vec3F
	}
)

func (s *GraphUI) OnInit() {
	s.entities = make(map[string]*Entity)
	s.tiles = make(map[int64]*Tile)
	s.grids = make([]*inner.Vec3F, 0)
	uiActor = s
}

func (s *GraphUI) random() {
	for i := 0; i < 50; i++ {
		time.Sleep(10 * time.Millisecond)
		randpos := &inner.Vec3F{X: float64(rand.Intn(constant.MAX_WORLD_WIDTH-100) + 1), Y: float64(rand.Intn(constant.MAX_WORLD_HEIGHT-100) + 1)}
		s.lastEID = tools.UUID()
		//randpos := &inner.Vec3F{X: float64(click.X), Y: float64(click.Y)}
		tempa := areaNum(int(randpos.X), int(randpos.Y))
		s.Send(common.WorldName(int32(tempa)),
			&inner.U2GAddEntity{
				AreaNum: int32(tempa),
				EID:     s.lastEID,
				RealPos: randpos,
			})

		path := []*inner.Vec3F{}
		for i := 0; i < 200; i++ {
			path = append(path, &inner.Vec3F{X: float64(rand.Intn(constant.MAX_WORLD_WIDTH-100) + 1), Y: float64(rand.Intn(constant.MAX_WORLD_HEIGHT-100) + 1)})
		}
		s.Send(common.WorldName(int32(tempa)), &inner.U2GMoveEntity{EID: s.lastEID, Path: path, Speed: float64(rand.Intn(1) + 20)})
	}
}

func (s *GraphUI) OnStop() bool {
	return true
}

func areaNum(x, y int) int {
	return (y/constant.SERVER_AREA_WIDTH)*3 + (x / constant.SERVER_AREA_HEIGHT)
}
func area(eid string) int {
	arr := strings.Split(eid, "_")
	i, _ := strconv.Atoi(arr[1])
	return i
}
func SEREID(eid string) string {
	arr := strings.Split(eid, "_")
	return arr[0]
}
func foramtEid(eid string, areaNum int64) string {
	return fmt.Sprintf("%v_%v", eid, areaNum)
}

func (s *GraphUI) OnHandleMessage(sourceId, targetId string, v interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	switch msg := v.(type) {
	case *inner.G2UCountEntity:
		if _, ok := s.tiles[msg.WorldId]; ok {
			s.tiles[msg.WorldId].Count = int64(msg.Count)
		}
		//delete(s.entities, foramtEid(msg.EID, int64(msg.AreaId)))
	case *inner.G2UDelEntity:
		delete(s.entities, foramtEid(msg.EID, int64(msg.GetAreaId())))
	case *inner.G2UUpdateEntity:
		id := foramtEid(msg.EID, int64(msg.GetAreaId()))
		if _, ok := s.entities[id]; !ok {
			s.entities[id] = &Entity{
				EID:  id,
				Path: make([]*inner.Vec3F, 0),
			}
		}
		s.entities[id].X = msg.RealPos.X
		s.entities[id].Y = msg.RealPos.Y
		s.entities[id].HP = msg.Hp
		s.entities[id].Duplicate = msg.Duplicate
	case *inner.G2UArea:
		s.tiles[msg.WorldId] = &Tile{
			ExclusiveBound: msg.ExclusiveBound,
			ActualBound:    msg.ActualBound,
		}
		s.grids = append(s.grids, msg.Point...)
	}
	return
}

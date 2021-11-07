package system

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/world/constant"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/ecs/entity"
	"github.com/wwj31/ecsDemo/world/interfaces"
	"github.com/wwj31/ecsDemo/world/utils"
	"reflect"
)

type GridTuple struct {
	posComp  *component.Position
	areaComp *component.AreaInfo
}

func (s *GridTuple) Init(comps map[ecs.ComponentType]ecs.IComponent) {
	var ok bool
	s.posComp, ok = comps[s.posComp.Type()].(*component.Position)
	expect.True(ok)
	s.areaComp, ok = comps[s.areaComp.Type()].(*component.AreaInfo)
	expect.True(ok)
}

// 格子系统
type GridSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewGridSys(w interfaces.IWorld) ecs.ISystem {
	ns := &GridSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(GRID_SYSTEM, w.Runtime(), reflect.TypeOf((*GridTuple)(nil))),
	}
	cross := w.Runtime().SingleComponent(component.CROSS_AREA_COMP).(*component.CrossArea)
	x, y := utils.Pos(int(w.AreaId()))
	area := cross.Areas[x][y]
	gridComp := w.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	buildGrid(area.ExtendBound, area.ActualBound, gridComp)
	return ns
}

func buildGrid(extendBound component.Bound, actualBound component.Bound, g *component.Grids) {
	for i := extendBound.X; i < extendBound.X+extendBound.Width; i = i + constant.GRID_SIZE {
		for j := extendBound.Y; j < extendBound.Y+extendBound.Height; j = j + constant.GRID_SIZE {
			pos := tools.Vec3f{X: i, Y: j}
			expect.True(utils.InBound(extendBound, pos))
			grid := &component.Grid{
				Entities: make(map[string]*ecs.Entity),
				Watchers: make(map[int64]bool),
			}
			//centerPos := tools.Add(pos, tools.Vec3f{X: constant.GRID_SIZE / 2, Y: constant.GRID_SIZE / 2})
			//if utils.InBound(actualBound, centerPos) { //真是区域的数据才能同步
			//	grid.NeedSync = true
			//}
			gk := utils.GetGridKey(tools.Vec3f{X: i, Y: j})
			g.Grids[gk] = grid
		}
	}
}

func (s *GridSys) UpdateFrame(float64) {
	moverComp := s.Runtime().SingleComponent(component.MOVER_COMP).(*component.MoverSet)
	SpawnComp := s.Runtime().SingleComponent(component.SPAWN_COMP).(*component.Spawn)
	deadComp := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
	for _, ent := range SpawnComp.Newcomes {
		s.NewEntity(ent.Ent.Id())
	}
	for _, ent := range deadComp.Deads {
		s.RemoveEntity(ent.Id())
	}
	for _, ent := range moverComp.Movers {
		s.EntityMove(ent)
	}
}

func (s *GridSys) EssentialComp() uint64 {
	return component.POS_COMP.ComponentType() | component.AREA_COMP.ComponentType()
}

func (s *GridSys) NewEntity(eid string) {
	tuple := s.GetTuple(eid).(*GridTuple)
	expect.True(tuple != nil, log.Fields{"eid": eid})

	gridKey := utils.GetGridKey(tuple.posComp.Pos)
	gridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	utils.AddEntityInGrid(s.world, gridComp, gridKey, eid)
	appears, _, _ := utils.SyncGrids(-1, -1, tuple.posComp.Pos.X, tuple.posComp.Pos.Y)
	if len(appears) > 0 {
		utils.SyncData(s.world, eid, appears, entity.TotalSyncData...) // 新实体，全量同步
	}
}

func (s *GridSys) RemoveEntity(eid string) {
	tuple := s.GetTuple(eid).(*GridTuple)
	expect.True(tuple != nil, log.Fields{"eid": eid})
	gridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	gridKey := utils.GetGridKey(tuple.posComp.OldPos)
	utils.DelEntityInGrid(s.world, gridComp, gridKey, eid)
	_, disappears, _ := utils.SyncGrids(tuple.posComp.OldPos.X, tuple.posComp.OldPos.Y, -1, -1)
	if len(disappears) > 0 {
		utils.SyncData(s.world, eid, disappears) //同步删除数据
	}
}

func (s *GridSys) EntityMove(eid string) {
	tuple := s.GetTuple(eid).(*GridTuple)
	expect.True(tuple != nil, log.Fields{"eid": eid})

	oldGridKey := utils.GetGridKey(tuple.posComp.OldPos)
	newGridKey := utils.GetGridKey(tuple.posComp.Pos)
	gridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	utils.MoveEntityInGrid(s.world, gridComp, oldGridKey, newGridKey, eid)
	// 副本实体，由实体真实区域同步，本区域忽略
	if tuple.areaComp.AreaId != s.world.AreaId() {
		return
	}

	appears, disappears, _ := utils.SyncGrids(tuple.posComp.OldPos.X, tuple.posComp.OldPos.Y, tuple.posComp.Pos.X, tuple.posComp.Pos.Y)
	gridLog.KVs(log.Fields{"actorId": s.world.GetID(), "eid": eid, "pos": tuple.posComp, "oldGridKey": oldGridKey, "newGridKey": newGridKey, "appears": appears, "disappears": disappears}).Yellow().Debug("grid move")
	if len(appears) > 0 {

		utils.SyncData(s.world, eid, appears, entity.TotalSyncData...) //进格子
	}
	if len(disappears) > 0 {
		utils.SyncData(s.world, eid, disappears) //出格子
	}
	//if len(modifyArr) > 0 {
	//	syncComp.AddSyncData(&component.SyncData{SyncGrid: modifyArr, Message:})
	//}
}

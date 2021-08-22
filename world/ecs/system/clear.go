package system

import (
	"github.com/wwj31/dogactor/ecs"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/interfaces"
)

// 实体输入系统
type ClearSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewClearSys(w interfaces.IWorld) ecs.ISystem {
	ns := &ClearSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(CLEAR_SYSTEM, w.Runtime(), nil),
	}
	return ns
}

func (s *ClearSys) UpdateFrame(float64) {
	inputComp := s.Runtime().SingleComponent(component.INPUT_COMP).(*component.InputSet)
	moveComp := s.Runtime().SingleComponent(component.MOVER_COMP).(*component.MoverSet)
	syncComp := s.Runtime().SingleComponent(component.SYNCINFO_COMP).(*component.SyncInfo)
	spawnComp := s.Runtime().SingleComponent(component.SPAWN_COMP).(*component.Spawn)

	// todo .... 考虑优化，减少gc
	if len(inputComp.Inputs) > 0 {
		inputComp.Inputs = make(map[string]*component.InputTuple)
	}
	if len(moveComp.Movers) > 0 {
		moveComp.Movers = make([]string, 0)
	}
	if len(syncComp.Data) > 0 {
		syncComp.Data = make([]*component.SyncData, 0)
	}
	if len(spawnComp.Newcomes) > 0 {
		spawnComp.Newcomes = make(map[string]*component.NewEntity)
	}
}
func (s *ClearSys) EssentialComp() uint64 {
	return component.FORBID_ENTITY
}

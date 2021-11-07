package system

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/ecsDemo/internal/common"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

// 创建实体系统
type SpawnSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewSpawnSys(w interfaces.IWorld) ecs.ISystem {
	ms := &SpawnSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(SPAWN_SYSTEM, w.Runtime(), nil),
	}
	return ms
}

func (s *SpawnSys) UpdateFrame(float64) {
	spawnComp := s.Runtime().SingleComponent(component.SPAWN_COMP).(*component.Spawn)
	for eid, ent := range spawnComp.Newcomes {
		s.world.Send(common.Center_Actor, &inner.W2CenterUpdateEntity{EId: eid, AreaId: s.world.AreaId()})
		s.world.Runtime().AddEntity(ent.Ent)
		ent.Ent.SetComponent(ent.Comps...)
		spawnLog.KVs(log.Fields{"eid": eid}).Debug("spawn entity")
	}
}
func (s *SpawnSys) EssentialComp() uint64 {
	return component.FORBID_ENTITY
}

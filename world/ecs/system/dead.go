package system

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/log"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/interfaces"
)

// 删除实体系统
type DeadSys struct {
	*ecs.SystemBase
	world interfaces.IWorld
}

func NewDeadSys(w interfaces.IWorld) ecs.ISystem {
	ms := &DeadSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(DEAD_SYSTEM, w.Runtime(), nil),
	}
	return ms
}

func (s *DeadSys) UpdateFrame(float64) {
	deads := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
	localAreaId := s.world.AreaId()
	for eid, ent := range deads.Deads {
		AreaComp := ent.GetComponent(component.AREA_COMP).(*component.AreaInfo)
		s.Runtime().DeleteEntity(eid)
		delete(deads.Deads, eid)

		s.world.Send(common.Center_Actor, &inner.W2CenterDeleteEntity{EId: ent.Id()})
		dealLog.KVs(log.Fields{"actorId": s.world.GetID(), "eid": eid, "AreaComp.AreaId ": AreaComp.AreaId, "AreaComp.DuplicatedAreas": AreaComp.DuplicatedAreas}).Debug("del entity ")
		if AreaComp.AreaId == localAreaId && AreaComp.DuplicatedAreas != nil {
			msg := &inner.W2WDelDuplicate{EId: eid}
			for n, _ := range AreaComp.DuplicatedAreas {
				s.world.Send(common.WorldName(int32(n)), msg)
			}
		}
	}
}
func (s *DeadSys) EssentialComp() uint64 {
	return component.FORBID_ENTITY
}

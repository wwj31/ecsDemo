package system

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/fsm"
	"github.com/wwj31/dogactor/log"
	"github.com/vmihailenco/msgpack"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/utils"
)

// 添加邻服实体副本
func (s *CrossAreaSys) W2WAddDuplicate(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2WAddDuplicate)
	if s.Runtime().GetEntity(data.Entity.EId) != nil {
		crossLog.Debug("跨区域不用添加")
		return
	}
	ent := ecs.NewEntity(data.Entity.EId)
	s.Runtime().AddEntity(ent) // 副本实体，不进spawn，收到消息直接添加
	s.updateData(ent, data.GetEntity())

	// 添加到地块上
	posComp := ent.GetComponent(component.POS_COMP).(*component.Position)
	gridKey := utils.GetGridKey(posComp.Pos)
	gridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
	utils.AddEntityInGrid(s.world, gridComp, gridKey, data.Entity.EId)
	crossLog.KVs(log.Fields{"eid": ent.Id(), "Target": s.world.GetID(), "Source": sourceId}).Debug("add dunplicate")
}

// 删除邻服实体副本
func (s *CrossAreaSys) W2WDelDuplicate(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2WDelDuplicate)
	tuple := s.GetTuple(data.EId).(*CrossTuple)

	if tuple != nil {
		logf := log.Fields{"eid": data.EId, "Target": s.world.GetID(), "Source": sourceId, "ent areaId": tuple.areaComp.AreaId, "world": s.world.AreaId()}
		if tuple.areaComp.AreaId != s.world.AreaId() {
			// 副本应该立刻删除
			s.Runtime().DeleteEntity(data.EId)

			// 地块上删除
			gridKey := utils.GetGridKey(tuple.posComp.Pos)
			gridComp := s.Runtime().SingleComponent(component.GRID_COMP).(*component.Grids)
			utils.DelEntityInGrid(s.world, gridComp, gridKey, data.EId)
			crossLog.KVs(logf).Debug("del dunplicate")
		} else {
			crossLog.KVs(logf).Error("why can recive??")
		}
	} else {
		crossLog.KVs(log.Fields{"eid": data.EId, "Target": s.world.GetID(), "Source": sourceId}).Warn("repeated del msg???")
	}
}

// 更新邻服实体副本
func (s *CrossAreaSys) W2WUpdateDuplicate(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2WUpdateDuplicate)
	ent := s.Runtime().GetEntity(data.Entity.EId)
	if ent == nil {
		crossLog.KVs(log.Fields{"eid": data.Entity.EId, "target": s.world.GetID(), "source": sourceId}).Warn("can not find entity")
		return
	}
	s.updateData(ent, data.GetEntity())

	tuple := s.GetTuple(data.Entity.EId).(*CrossTuple)
	crossAreaComp := s.Runtime().SingleComponent(component.CROSS_AREA_COMP).(*component.CrossArea)
	inputComp := s.Runtime().SingleComponent(component.INPUT_COMP).(*component.InputSet)
	deadComp := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
	s.checkPos(ent.Id(), crossAreaComp, inputComp, tuple.areaComp, deadComp, tuple.posComp)
	crossLog.KVs(log.Fields{"eid": ent.Id(), "Target": s.world.GetID(), "Source": sourceId}).Debug("up dunplicate")
}

func (s *CrossAreaSys) updateData(entity *ecs.Entity, entInfo *inner.EntityInfo) {
	entity.ResetComponent()

	comps := map[ecs.ComponentType]ecs.IComponent{}
	arrcomp := []ecs.IComponent{}
	for k, v := range entInfo.Data {
		t := component.ComponentType(k)
		comps[t] = component.COMPONENTS[t]()
		msgpack.Unmarshal(v, comps[t])
		arrcomp = append(arrcomp, comps[t])
	}
	entity.SetComponent(arrcomp...)

	expect.True(comps[component.AREA_COMP] != nil)
	areaComp := comps[component.AREA_COMP].(*component.AreaInfo)
	// 恢复战斗状态机
	var fightingComp *component.Fighting
	if comps[component.FIGHTING_COMP] != nil {
		fightingComp = comps[component.FIGHTING_COMP].(*component.Fighting)
	}
	if areaComp.AreaId == s.world.AreaId() && fightingComp != nil {
		fightingComp.Fight = fsm.New()
		fightingComp.Fight.Add(NewFightState(s.world, entity.Id()))
		fightingComp.Fight.ForceState(TURN)
		fightLog.KVs(log.Fields{"actorId": s.world.GetID(), "eid": entity.Id()}).White().Debug("recover fight")
	}
}

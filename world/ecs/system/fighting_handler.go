package system

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"ecsDemo/world/ecs/component"
	"ecsDemo/world/utils"
)

// 攻击消息
func (s *FightingSys) W2WAttackReq(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2WAttackReq)
	logfields := log.Fields{"attacker": data.AttackerEId, "defender": data.DefenderEId, "actorId": s.world.GetID()}
	defender := s.Runtime().GetEntity(data.DefenderEId)
	if defender == nil {
		//fightLog.KVs(logfields).Warn("找不到攻击目标")
		return
	}

	areaInfo := defender.GetComponent(component.AREA_COMP).(*component.AreaInfo)
	if areaInfo.AreaId != s.world.AreaId() {
		fightLog.KVs(logfields).KV("ent areaId", areaInfo.AreaId).Warn("attack msg, not entity for this world")
		s.world.Send(common.WorldName(areaInfo.AreaId), msg)
		return
	}

	// 收到被攻击的消息，还没进入战斗，创建战斗相关组件
	if defender.GetComponent(component.FIGHTING_COMP) == nil {
		fightLog.KVs(logfields).Debug("enter turn for be attacked")
		newFightComp := component.NewFighting(NewFightState(s.world, defender.Id()), utils.NewTurnCount())
		defender.SetComponent(newFightComp)
		newFightComp.Fight.Switch(TURN)
	}
	fightComp := defender.GetComponent(component.FIGHTING_COMP).(*component.Fighting)
	expect.True(fightComp.Fight.State() > 0, logfields)
	fightComp.Fight.Handle(StateData{sourceId: sourceId, pbMsg: msg})
}

// 反击消息
func (s *FightingSys) W2WAttackResp(sourceId string, msg interface{}, gateSession string) {
	data := msg.(*inner.W2WAttackResp)
	attacker := s.Runtime().GetEntity(data.AttackerEId)
	if attacker == nil {
		fightLog.KVs(log.Fields{"eid": data.DefenderEId}).Warn("找不到反击目标")
		return
	}
	logfields := log.Fields{"attacker": data.AttackerEId, "defender": data.DefenderEId, "actorId": s.world.GetID()}
	areaInfo := attacker.GetComponent(component.AREA_COMP).(*component.AreaInfo)
	if areaInfo.AreaId != s.world.AreaId() {
		fightLog.KVs(logfields).KV("ent areaId", areaInfo.AreaId).Warn("counter msg,not entity for this world")
		s.world.Send(common.WorldName(areaInfo.AreaId), msg)
		return
	}
	if attacker.GetComponent(component.FIGHTING_COMP) == nil {
		fightLog.KVs(log.Fields{"defender": data.DefenderEId, "attacker": data.AttackerEId}).Warn("目标已经脱离战斗，不能接受反击")
		return
	}

	fightComp := attacker.GetComponent(component.FIGHTING_COMP).(*component.Fighting)
	expect.True(fightComp.Fight.State() > 0, logfields)
	fightComp.Fight.Handle(StateData{sourceId: sourceId, pbMsg: msg})
}

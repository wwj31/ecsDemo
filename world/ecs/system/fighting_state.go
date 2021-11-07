package system

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/fsm"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/internal/common"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/internal/message"
	"github.com/wwj31/ecsDemo/world/constant"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/ecs/entity"
	"github.com/wwj31/ecsDemo/world/interfaces"
	"github.com/wwj31/ecsDemo/world/utils"
)

const TURN = 1

type FightState struct {
	interfaces.IWorld
	eid string
}

func NewFightState(w interfaces.IWorld, eid string) *FightState {
	return &FightState{
		IWorld: w,
		eid:    eid,
	}
}

func (s *FightState) State() int {
	return TURN
}

// 进入状态，回合开始
func (s *FightState) Enter(*fsm.FSM) {
	Ent := s.Runtime().GetEntity(s.eid)
	expect.True(Ent != nil)
	if Ent.GetComponent(component.FIGHTING_COMP) == nil {
		// 战斗结束会删除FightComp，这里是正常情况
		return
	}

	fightingComp := Ent.GetComponent(component.FIGHTING_COMP).(*component.Fighting)
	attriComp := Ent.GetComponent(component.ATTRIBUTE_COMP).(*component.Attribute)
	fightingComp.TurnNumber++

	//if attriComp.Target == "" {
	//	attriComp.Target = _lasttarget(fightingComp)
	//}

	logf := log.Fields{"actorId": s.GetID(), "eid": Ent.Id(), "target": attriComp.Target, "turnNumber": fightingComp.TurnNumber}
	fightLog.KVs(logf).Green().Debug("startTurn")

	if attriComp.Target != "" {
		targetEnt := s.Runtime().GetEntity(attriComp.Target)
		if targetEnt != nil {
			deadComp := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
			if _, ok := deadComp.Deads[attriComp.Target]; !ok {
				// 攻击距离是否足够
				entPos := Ent.GetComponent(component.POS_COMP).(*component.Position).Pos
				targetPos := targetEnt.GetComponent(component.POS_COMP).(*component.Position).Pos
				dist := tools.Distance(entPos, targetPos)
				if dist < constant.ATTACK_DIST {
					attMsg := &inner.W2WAttackReq{
						AttackerEId:  Ent.Id(),
						DefenderEId:  attriComp.Target,
						TurnNumber:   fightingComp.TurnNumber,
						AttackerAttr: attriComp.InnerPB(),
					}
					targetAreaId := targetEnt.GetComponent(component.AREA_COMP).(*component.AreaInfo).AreaId

					s.Send(common.WorldName(targetAreaId), attMsg)
					fightLog.KVs(logf).KVs(log.Fields{"dist": dist}).Yellow().Debug("attack target")
				} else {
					fightLog.KVs(logf).KVs(log.Fields{"dist": dist}).Red().Debug("out of ATTACK_DIST")
				}
			}
		} else {
			fightLog.KVs(logf).Debug("can not found targetEId")
		}
	} else {
		fightLog.KVs(log.Fields{"actorId": s.GetID(), "eid": Ent.Id()}).Debug("entity not attack target")
	}
	// 把上一回合提前收到的攻击消息，丢给系统处理
	turn := _currentTurn(fightingComp)
	for attackerEId, msg := range turn.AttMsgs {
		attacker := s.Runtime().GetEntity(attackerEId)
		if attacker != nil {
			comp := attacker.GetComponent(component.AREA_COMP).(*component.AreaInfo)
			s.System().Send(common.WorldName(comp.AreaId), s.GetID(), "", msg)
		}
	}
}

// 退出状态，回合结算
func (s *FightState) Leave(*fsm.FSM) {
	Ent := s.Runtime().GetEntity(s.eid)
	expect.True(Ent != nil)
	fightingComp := Ent.GetComponent(component.FIGHTING_COMP).(*component.Fighting)
	seTurn := _currentTurn(fightingComp)
	attriComp := Ent.GetComponent(component.ATTRIBUTE_COMP).(*component.Attribute)
	attriComp.Sets[message.FIGHT_ATTR_HP.Int32()] -= seTurn.TotalDamage
	attriComp.Sets[message.FIGHT_ATTR_HP.Int32()] -= seTurn.Counter

	logf := log.Fields{
		"actorId":     s.GetID(),
		"starget":     attriComp.Target,
		"attacker":    seTurn.Attackers,
		"TotalDamage": seTurn.TotalDamage,
		"Counter":     seTurn.Counter,
		"eid":         Ent.Id(),
		"att":         attriComp.Sets,
		"turnNumber":  fightingComp.TurnNumber,
	}

	defer func() {
		areaComp := Ent.GetComponent(component.AREA_COMP).(*component.AreaInfo)
		utils.SyncAreaDuplicate(entity.AreaDuplicate(s, s.eid), areaComp.DuplicatedAreas, s)
	}()

	fightLog.KVs(logf).Green().Debug("settleTurn")

	// todo ....... 血量为0删除实体，这里没有需求，先简易处理
	if attriComp.Sets[message.FIGHT_ATTR_HP.Int32()] <= 0 {
		attriComp.Sets[message.FIGHT_ATTR_HP.Int32()] = 0
		dead := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
		dead.Deads[Ent.Id()] = Ent
	}

	if s.fightOver(Ent, fightingComp) {
		attriComp.Target = ""
		// todo ....... 生产战报
		Ent.DelComponent(fightingComp)
		fightLog.KVs(logf).Red().Debug("fight over")
	}
}

type StateData struct {
	sourceId string
	pbMsg    interface{}
}

// 回合过程中，处理消息
func (s *FightState) Handle(fms *fsm.FSM, v interface{}) {
	switch data := v.(type) {
	case StateData:
		switch data.pbMsg.(type) {
		case *inner.W2WAttackReq: // 受到某单位攻击
			s.W2WAttackReq(data.sourceId, data.pbMsg)
		case *inner.W2WAttackResp: // 受到某单位反击
			s.W2WAttackResp(data.sourceId, data.pbMsg)
		}
	default:
		log.KVs(log.Fields{"type": data}).Error("undefined type")
	}
}

// 攻击消息
func (s *FightState) W2WAttackReq(sourceId string, pbMsg interface{}) {
	msg := pbMsg.(*inner.W2WAttackReq)

	Ent := s.Runtime().GetEntity(s.eid)
	expect.True(Ent != nil)
	deffightComp, ok := Ent.GetComponent(component.FIGHTING_COMP).(*component.Fighting)
	expect.True(ok, log.Fields{"eid": msg.DefenderEId})
	expect.True(int(deffightComp.TurnNumber) <= len(deffightComp.Turns), log.Fields{"eid": msg.DefenderEId, "TurnNumber": deffightComp.TurnNumber, "len(turns)": len(deffightComp.Turns)})

	turn := _currentTurn(deffightComp)
	// 一个回合内，一个攻击者只能攻击一次，多出来的次数，放入下一个回合
	if _, ok = turn.Attackers[msg.AttackerEId]; ok {
		nextturn := _nextTurn(deffightComp)
		attmsgs := nextturn.AttMsgs
		expect.True(attmsgs[msg.AttackerEId] == nil, log.Fields{"opps! 这尼玛是服务器卡了吧 turnNumber": deffightComp.TurnNumber})
		attmsgs[msg.AttackerEId] = msg
		return
	}

	damage, counter := _fighting(msg.AttackerAttr, Ent, turn.TotalDamage)
	turn.Attackers[msg.AttackerEId] = damage
	turn.TotalDamage += damage
	attResp := &inner.W2WAttackResp{
		AttackerEId: msg.AttackerEId,
		DefenderEId: msg.DefenderEId,
		TurnNumber:  msg.TurnNumber,
		DamageVal:   damage,
		CounterVal:  counter,
	}

	s.Send(sourceId, attResp)
	fightLog.KVs(log.Fields{
		"actorId":  s.GetID(),
		"attacker": msg.AttackerEId,
		"defender": msg.DefenderEId,
		"damage":   damage,
		"counter":  counter,
	}).White().Debug("be attacked")
}

// 反击消息
func (s *FightState) W2WAttackResp(sourceId string, pbMsg interface{}) {
	msg := pbMsg.(*inner.W2WAttackResp)

	Ent := s.Runtime().GetEntity(s.eid)
	expect.True(Ent != nil)
	attacker := Ent

	fightComp, ok := attacker.GetComponent(component.FIGHTING_COMP).(*component.Fighting)
	errlog := log.Fields{"attacker": msg.AttackerEId, "defender": msg.DefenderEId, "atter turnnum": fightComp.TurnNumber, "msg turnnum": msg.TurnNumber, "len(turns)": len(fightComp.Turns), "actorId": s.GetID()}
	expect.True(ok, errlog)
	expect.True(fightComp.TurnNumber == msg.TurnNumber, errlog)
	expect.True(int(fightComp.TurnNumber) <= len(fightComp.Turns), errlog)

	turn := _currentTurn(fightComp)
	turn.TargetDamage = msg.DamageVal
	turn.Counter = msg.CounterVal
	fightLog.KVs(log.Fields{"attacker": msg.AttackerEId, "defender": msg.DefenderEId, "Damage": msg.DamageVal, "Counter": msg.CounterVal, "actorId": s.GetID()}).White().Debug("Counter")
}

////////////////////////////////////////////////////// 内部函数 ////////////////////////////////////////////////////////////
func (s *FightState) fightOver(ent *ecs.Entity, fightingComp *component.Fighting) bool {
	dead := s.Runtime().SingleComponent(component.DEAD_COMP).(*component.DeadEntities)
	if dead.Deads[ent.Id()] != nil {
		return true
	}
	// 本回合没有受到任何攻击，且 没有攻击任何人，战斗结束
	turn := _currentTurn(fightingComp)
	nextTurn := _nextTurn(fightingComp)
	if turn.TotalDamage == 0 &&
		turn.TargetDamage == 0 && turn.Counter == 0 &&
		len(nextTurn.AttMsgs) == 0 {
		return true
	}
	return false
}

/*
伤害、反击计算
attackerAttr 攻击者战斗属性
defender	防御者entity
totalDamage 当前受到的总伤害

dval		攻击伤害
cval		反击伤害
*/
func _fighting(attackerAttr *inner.FightAttr, defender *ecs.Entity, totalDamage int64) (dval, cval int64) {
	defAttrComp := defender.GetComponent(component.ATTRIBUTE_COMP).(*component.Attribute)
	attacker_ATT := attackerAttr.Sets[message.FIGHT_ATTR_ATT.Int32()]
	attacker_HP := attackerAttr.Sets[message.FIGHT_ATTR_HP.Int32()]
	bearHP := tools.If(defAttrComp.HP() > totalDamage, defAttrComp.HP()-totalDamage, int64(0)).(int64) // 本次攻击，最多还能承受多少伤害
	dval = tools.If(attacker_ATT < bearHP, attacker_ATT, bearHP).(int64)                               // 本次攻击，最多能造成多少伤害
	cval = tools.If(defAttrComp.ATT() < attacker_HP, defAttrComp.ATT(), attacker_HP).(int64)           // 本次反击，最多能造成多少伤害
	return attacker_ATT, defAttrComp.ATT()
}

// 获得本回合，没有就创建
func _currentTurn(fightingComp *component.Fighting) *component.TurnData {
	expect.True(fightingComp.TurnNumber > 0, log.Fields{"turnNumber": fightingComp.TurnNumber, "len(turns)": len(fightingComp.Turns)})
	if int(fightingComp.TurnNumber) >= len(fightingComp.Turns) {
		fightingComp.Turns = append(fightingComp.Turns, component.NewTurnData())
	}
	return fightingComp.Turns[fightingComp.TurnNumber-1]
}

// 获得下一个回合，没有就创建
func _nextTurn(fightingComp *component.Fighting) *component.TurnData {
	expect.True(fightingComp.TurnNumber >= 0)
	expect.True(int(fightingComp.TurnNumber) <= len(fightingComp.Turns))
	if int(fightingComp.TurnNumber) == len(fightingComp.Turns) {
		fightingComp.Turns = append(fightingComp.Turns, component.NewTurnData())
	}

	return fightingComp.Turns[fightingComp.TurnNumber]
}

//// 获得上一回合的目标
//func _lasttarget(fightingComp *component.Fighting) string {
//	lastNum := fightingComp.TurnNumber - 1
//	if lastNum <= 0 {
//		return ""
//	}
//	expect.True(len(fightingComp.Turns) >= int(lastNum-1))
//	return fightingComp.Turns[lastNum-1].Target
//}

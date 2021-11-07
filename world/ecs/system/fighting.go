package system

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/log"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
	"github.com/wwj31/ecsDemo/world/utils"
	"reflect"
)

type FightTuple struct {
	posComp   *component.Position
	areaComp  *component.AreaInfo
	fightComp *component.Fighting
	attriComp *component.Attribute
}

func (s *FightTuple) Init(comps map[ecs.ComponentType]ecs.IComponent) {
	var ok bool
	s.posComp, ok = comps[s.posComp.Type()].(*component.Position)
	expect.True(ok)

	s.areaComp, ok = comps[s.areaComp.Type()].(*component.AreaInfo)
	expect.True(ok)

	s.fightComp, ok = comps[s.fightComp.Type()].(*component.Fighting)
	expect.True(ok)

	s.attriComp, ok = comps[s.attriComp.Type()].(*component.Attribute)
	expect.True(ok)

}

// 战斗系统
type (
	FightingSys struct {
		*ecs.SystemBase
		world interfaces.IWorld
	}
)

func NewFightingSys(w interfaces.IWorld) ecs.ISystem {
	ns := &FightingSys{
		world:      w,
		SystemBase: ecs.NewSystemBase(FIGHTING_SYSTEM, w.Runtime(), reflect.TypeOf((*FightTuple)(nil))),
	}

	ns.world.RegistMsg((*inner.W2WAttackReq)(nil), ns.W2WAttackReq)
	ns.world.RegistMsg((*inner.W2WAttackResp)(nil), ns.W2WAttackResp)
	return ns
}

func (s *FightingSys) EssentialComp() uint64 {
	return component.POS_COMP.ComponentType() |
		component.AREA_COMP.ComponentType() |
		component.FIGHTING_COMP.ComponentType() |
		component.ATTRIBUTE_COMP.ComponentType()
}

func (s *FightingSys) UpdateFrame(float64) {
	//每个战斗实体，先处理战斗输入.通过计数器判断是否结算
	fightInputComp := s.Runtime().SingleComponent(component.INPUT_COMP).(*component.InputSet)
	s.Range(func(eid string, t ecs.ITuple) bool {
		tuple, ok := t.(*FightTuple)
		expect.True(ok, log.Fields{"eid": eid})
		// 副本实体不处理
		if tuple.areaComp.AreaId != s.world.AreaId() {
			return true
		}

		input := fightInputComp.Inputs[eid]
		if input != nil && input.Attack != nil {
			// 有新的目标，放入下一回合中，如果没进战斗,直接进入战斗开始第1回合，否则就等待下回合开始攻击目标
			//nextTurn := _nextTurn(tuple.fightComp)
			tuple.attriComp.Target = input.Attack.TargetEId
			if tuple.fightComp.Fight.State() == -1 {
				fightLog.Debug("有攻击目标，战斗回合开始！")
				tuple.fightComp.Fight.Switch(TURN)
			}
		}
		tuple.fightComp.TurnCount--
		if tuple.fightComp.TurnCount == 0 {
			tuple.fightComp.TurnCount = utils.NewTurnCount()
			tuple.fightComp.Fight.Switch(TURN)
			// 每回合结算后，同步一次实体数据
			utils.SyncData(s.world, eid, utils.AroundGrid(tuple.posComp.Pos), component.ATTRIBUTE_COMP) // 战斗
		}
		return true
	})
	fightLog.KVs(log.Fields{"count": s.world.FC(), "actorId": s.world.GetID()}).Debug("FRAME")
}

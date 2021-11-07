package component

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/fsm"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
)

type (
	Fighting struct {
		ecs.ComponentBase
		TurnNumber int32       // 回合数
		Turns      []*TurnData // 历史回合数据
		TurnCount  int64       // 下个回合开始需要执行帧计数
		Fight      *fsm.FSM    `msgpack:"-"` // 战斗状态
	}
	TurnData struct {
		AttMsgs     map[string]*inner.W2WAttackReq // 本回合开始前收到的攻击消息
		Attackers   map[string]int64               // 本回合攻击过我的人 map[攻击者EID]Damage
		TotalDamage int64                          // 本回合总伤害

		TargetDamage int64 // 本回合对目标造成的伤害
		Counter      int64 // 本回合目标的反击伤害
	}
)

func NewFighting(state fsm.StateHandler, tc int64) *Fighting {
	m := &Fighting{TurnNumber: 0, Turns: make([]*TurnData, 0), TurnCount: tc}
	m.Fight = fsm.New()
	m.Fight.Add(state)
	return m
}
func NewTurnData() *TurnData {
	m := &TurnData{Attackers: make(map[string]int64), AttMsgs: make(map[string]*inner.W2WAttackReq)}
	return m
}

func (s *Fighting) Type() ecs.ComponentType {
	return FIGHTING_COMP
}

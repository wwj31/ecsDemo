package component

import "github.com/wwj31/dogactor/ecs"

type (
	Player struct {
		ecs.ComponentBase
		RID         int64 // 关联的角色ID
		GateSession string
	}
)

func NewPlayer(rid int64, gateSession string) *Player {
	m := &Player{RID: rid, GateSession: gateSession}
	return m
}

func (s *Player) Type() ecs.ComponentType {
	return PLAY_COMP
}

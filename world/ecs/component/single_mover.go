package component

import (
	"github.com/wwj31/dogactor/ecs"
)

// 保存一帧所有移动过的实体
type (
	MoverSet struct {
		ecs.ComponentBase
		Movers []string
	}
)

func NewMoverSet() *MoverSet {
	v := &MoverSet{Movers: make([]string, 0)}
	return v
}

func (s *MoverSet) Type() ecs.ComponentType {
	return MOVER_COMP
}

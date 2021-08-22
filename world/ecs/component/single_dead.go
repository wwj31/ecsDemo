package component

import (
	"github.com/wwj31/dogactor/ecs"
)

/*
	管理实体的销毁
*/
type (
	DeadEntities struct {
		ecs.ComponentBase
		Deads map[string]*ecs.Entity
	}
)

func NewDeadEntity() *DeadEntities {
	v := &DeadEntities{
		Deads: make(map[string]*ecs.Entity),
	}
	return v
}

func (s *DeadEntities) Type() ecs.ComponentType {
	return DEAD_COMP
}

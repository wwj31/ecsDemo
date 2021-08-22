package component

import (
	"github.com/wwj31/dogactor/ecs"
)

/*
	管理实体的创建和销毁
*/
type (
	Spawn struct {
		ecs.ComponentBase
		Newcomes map[string]*NewEntity
	}
	NewEntity struct {
		Ent   *ecs.Entity
		Comps []ecs.IComponent
	}
)

func NewSpawnEntity() *Spawn {
	v := &Spawn{
		Newcomes: make(map[string]*NewEntity),
	}
	return v
}

func (s *Spawn) Type() ecs.ComponentType {
	return SPAWN_COMP
}

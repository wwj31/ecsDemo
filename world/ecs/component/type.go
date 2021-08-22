package component

import (
	"github.com/wwj31/dogactor/ecs"
	"math"
)

type ComponentType uint64

func (s ComponentType) ComponentType() uint64 {
	return uint64(s)
}

const (
	MOVE_COMP ComponentType = 1 << iota
	AREA_COMP
	POS_COMP
	ATTRIBUTE_COMP
	PLAY_COMP
	FIGHTING_COMP

	EMAX ComponentType = 1 << 62 // 最多62个
)

// 单列组件
const (
	CROSS_AREA_COMP ComponentType = math.MaxInt64 - iota
	SPAWN_COMP
	DEAD_COMP
	INPUT_COMP
	GRID_COMP
	MOVER_COMP
	SYNCINFO_COMP
)

const EVERYONES = 0
const FORBID_ENTITY = math.MaxInt64

var COMPONENTS = map[ComponentType]func() ecs.IComponent{
	MOVE_COMP:      func() ecs.IComponent { return &Move{} },
	AREA_COMP:      func() ecs.IComponent { return &AreaInfo{} },
	POS_COMP:       func() ecs.IComponent { return &Position{} },
	ATTRIBUTE_COMP: func() ecs.IComponent { return &Attribute{} },
	PLAY_COMP:      func() ecs.IComponent { return &Player{} },
	FIGHTING_COMP:  func() ecs.IComponent { return &Fighting{} },
}

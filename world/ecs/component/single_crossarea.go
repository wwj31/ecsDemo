package component

import (
	"ecsDemo/world/constant"
	"github.com/wwj31/dogactor/ecs"
)

/*
       ┌ ── ── ── ── ── ── ── ── ── ── ── ─┐
       │ExtendBound                        │
       │    ┌ ── ── ── ── ── ── ── ── ┐    │
       │    │ActualBound              │    │
       │    │    ┌── ── ── ── ── ┐    │    │
       │    │    │               │    │    │
 Area =│    │    │ExclusiveBound │    │    │
       │    │    │               │    │    │
       │    │    │               │    │    │
       │    │    └── ── ── ── ── ┘    │    │
       │    │                         │    │
       │    └ ── ── ── ── ── ── ── ── ┘    │
       │                                   │
       └ ── ── ── ── ── ── ── ── ── ── ── ─┘
*/
type Bound struct {
	X, Y          float64 // 左上坐标
	Width, Height float64 // 宽高
}
type (
	// 区域服管理边界
	Area struct {
		ExclusiveBound Bound // 独享的区域
		ActualBound    Bound // 实际管理的区域
		ExtendBound    Bound // 扩展区域
	}

	CrossArea struct {
		ecs.ComponentBase
		Areas [constant.SERVER_SPLIT_AREA][constant.SERVER_SPLIT_AREA]*Area
	}
)

func NewCrossArea() *CrossArea {
	v := &CrossArea{}
	return v
}

func (s *CrossArea) Type() ecs.ComponentType {
	return CROSS_AREA_COMP
}

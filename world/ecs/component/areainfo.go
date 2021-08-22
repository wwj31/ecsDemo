package component

import "github.com/wwj31/dogactor/ecs"

type AreaInfo struct {
	ecs.ComponentBase
	AreaId          int32        // 实体所在区域Id
	DuplicatedAreas map[int]bool // 其他区域副本
}

func NewAreaInfo(aid int32) *AreaInfo {
	m := &AreaInfo{AreaId: aid, DuplicatedAreas: make(map[int]bool)}
	return m
}

func (s *AreaInfo) Type() ecs.ComponentType {
	return AREA_COMP
}

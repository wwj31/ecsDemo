package component

import (
	"github.com/wwj31/dogactor/ecs"
)

type (
	SyncInfo struct {
		ecs.ComponentBase
		Data []*SyncData
	}
)

func NewSyncGridsData() *SyncInfo {
	v := &SyncInfo{Data: make([]*SyncData, 0)}
	return v
}

func (s *SyncInfo) Type() ecs.ComponentType {
	return SYNCINFO_COMP
}

type SyncData struct {
	EId       string
	SyncGrids []int32 // 同步数据给地块里的所有玩家
	Comp      map[ecs.ComponentType]bool
}

func NewSyncData() *SyncData {
	return &SyncData{Comp: make(map[ecs.ComponentType]bool)}
}

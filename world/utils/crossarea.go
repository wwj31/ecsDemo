package utils

import (
	"ecsDemo/internal/common"
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/world/constant"
	"ecsDemo/world/interfaces"
)

func SyncAreaDuplicate(entData *inner.EntityInfo, areaId map[int]bool, world interfaces.IWorld) {
	msg := &inner.W2WUpdateDuplicate{Entity: entData}
	for id, _ := range areaId {
		world.Send(common.WorldName(int32(id)), msg)
	}
}

func Pos(i int) (x, y int) {
	return i % constant.SERVER_SPLIT_AREA, i / constant.SERVER_SPLIT_AREA
}

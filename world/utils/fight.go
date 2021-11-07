package utils

import "github.com/wwj31/ecsDemo/world/constant"

func NewTurnCount() int64 {
	return constant.FIGHTING_FRAME / constant.FRAME_RATE
}

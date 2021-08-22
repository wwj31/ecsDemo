package utils

import (
	"ecsDemo/world/constant"
)

func NewTurnCount() int64 {
	return constant.FIGHTING_FRAME / constant.FRAME_RATE
}

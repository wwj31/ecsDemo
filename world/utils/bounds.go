package utils

import (
	"github.com/wwj31/dogactor/tools"
	"ecsDemo/world/ecs/component"
)

func InBound(bound component.Bound, f tools.Vec3f) bool {
	if bound.X <= f.X && f.X <= bound.X+bound.Width &&
		bound.Y <= f.Y && f.Y <= bound.Y+bound.Height {
		return true
	}
	return false
}

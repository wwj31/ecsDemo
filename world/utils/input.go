package utils

import (
	"github.com/wwj31/ecsDemo/world/ecs/component"
	"github.com/wwj31/ecsDemo/world/interfaces"
)

// 拿不到就创建
func GetDataWithBuild(w interfaces.IWorld, eid string) *component.InputTuple {
	inputComp := w.Runtime().SingleComponent(component.INPUT_COMP).(*component.InputSet)
	if _, ok := inputComp.Inputs[eid]; !ok {
		inputComp.Inputs[eid] = &component.InputTuple{}
	}
	return inputComp.Inputs[eid]
}

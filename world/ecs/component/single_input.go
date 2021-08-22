package component

import (
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/tools"
)

// 保存帧间实体的输入操作
type (
	InputSet struct {
		ecs.ComponentBase
		Inputs map[string]*InputTuple
	}

	// 实体的所有输入数据
	InputTuple struct {
		Move   *MoveInputData
		Attack *AttackInputData
		// todo ...other input data
	}

	// 移动操作相关数据
	MoveInputData struct {
		Path  []tools.Vec3f
		Speed float64
	}

	// 攻击操作相关数据
	AttackInputData struct {
		TargetEId string
	}
)

func NewInputSet() *InputSet {
	v := &InputSet{Inputs: make(map[string]*InputTuple)}
	return v
}

func (s *InputSet) Type() ecs.ComponentType {
	return INPUT_COMP
}

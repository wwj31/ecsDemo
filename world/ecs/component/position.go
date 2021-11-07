package component

import (
	"github.com/golang/protobuf/proto"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/internal/message"
	"github.com/wwj31/ecsDemo/internal/msgtools"
)

type Position struct {
	ecs.ComponentBase
	Pos    tools.Vec3f // 当前位置
	OldPos tools.Vec3f // 上一帧位置
}

func NewPosition(pos tools.Vec3f) *Position {
	m := &Position{Pos: pos, OldPos: tools.Invalid()}
	return m
}

func (s *Position) Type() ecs.ComponentType {
	return POS_COMP
}

func (s *Position) SyncData() *message.SyncData {
	byt, err := proto.Marshal(&message.Position{
		Pos: msgtools.Vec3f2Msg(s.Pos),
	})
	expect.Nil(err)
	return &message.SyncData{
		Type: message.SyncDataType_POSITION,
		Data: byt,
	}
}

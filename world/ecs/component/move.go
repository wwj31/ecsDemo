package component

import (
	"ecsDemo/internal/message"
	"ecsDemo/internal/msgtools"
	"github.com/golang/protobuf/proto"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/wwj31/dogactor/tools"
)

type Option func(*Move)
type Move struct {
	ecs.ComponentBase
	OriginPos tools.Vec3f // 起点位置
	DestPos   tools.Vec3f // 终点位置

	Velocity float64     // 速度值 m/s
	Offset   tools.Vec3f // 速度向量(每帧偏移值)
	Distance float64     // 起点到终点总长度
	Lastf    float64     // 上一次偏移的帧率

	Path []tools.Vec3f // 路径
}

func NewMove(op ...Option) *Move {
	m := &Move{}
	for _, f := range op {
		f(m)
	}
	return m
}

func Pos(v tools.Vec3f) Option {
	return func(l *Move) {
		l.OriginPos = v
		l.DestPos = v
	}
}

func (s *Move) Type() ecs.ComponentType {
	return MOVE_COMP
}

func (s *Move) SyncData() *message.SyncData {
	byt, err := proto.Marshal(&message.Velocity{
		Velocity: s.Velocity,
		Path:     msgtools.Vec3f2Msg_Arr(s.Path),
	})
	expect.Nil(err)
	return &message.SyncData{
		Type: message.SyncDataType_VELOCITY,
		Data: byt,
	}
}

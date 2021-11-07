package interfaces

import (
	"github.com/wwj31/dogactor/actor"
	"github.com/wwj31/dogactor/ecs"
	"github.com/golang/protobuf/proto"
)

type IWorld interface {
	actor.Actor
	RegistMsg(msg proto.Message, f RegistFun)
	AreaId() int32
	Runtime() *ecs.Runtime
	FC() int64
}

type RegistFun func(sourceId string, msg interface{}, gateSession string)

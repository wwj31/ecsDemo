package msgtools

import (
	"github.com/wwj31/dogactor/tools"
	"github.com/wwj31/ecsDemo/internal/inner_message/inner"
	"github.com/wwj31/ecsDemo/internal/message"
)

func Vec3f2Inner(v tools.Vec3f) *inner.Vec3F {
	return &inner.Vec3F{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}

func Vec3f2Inner_Arr(arr []tools.Vec3f) []*inner.Vec3F {
	v3f := make([]*inner.Vec3F, 0, len(arr))
	for _, v := range arr {
		v3f = append(v3f, Vec3f2Inner(v))
	}
	return v3f
}

func Inner2Vec3f(v *inner.Vec3F) tools.Vec3f {
	return tools.Vec3f{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}
func Inner2Vec3f_Arr(arr []*inner.Vec3F) []tools.Vec3f {
	v3f := make([]tools.Vec3f, 0, len(arr))
	for _, v := range arr {
		v3f = append(v3f, Inner2Vec3f(v))
	}
	return v3f
}

func Msg2Vec3f(v *message.Vec3F) tools.Vec3f {
	return tools.Vec3f{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}
func Msg2Vec3f_Arr(arr []*message.Vec3F) []tools.Vec3f {
	v3f := make([]tools.Vec3f, 0, len(arr))
	for _, v := range arr {
		v3f = append(v3f, Msg2Vec3f(v))
	}
	return v3f
}

func Vec3f2Msg(v tools.Vec3f) *message.Vec3F {
	return &message.Vec3F{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}
func Vec3f2Msg_Arr(arr []tools.Vec3f) []*message.Vec3F {
	v3f := make([]*message.Vec3F, 0, len(arr))
	for _, v := range arr {
		v3f = append(v3f, Vec3f2Msg(v))
	}
	return v3f
}

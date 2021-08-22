package component

import (
	"ecsDemo/internal/inner_message/inner"
	"ecsDemo/internal/message"
	"github.com/wwj31/dogactor/ecs"
	"github.com/wwj31/dogactor/expect"
	"github.com/golang/protobuf/proto"
	"math/rand"

)

/*
	trun-based game,一个回合内，属性始终保持不变，回合结束结算一次
*/
type Attribute struct {
	ecs.ComponentBase
	Sets   map[int32]int64
	Target string // 攻击目标
}

func NewAttribute() *Attribute {
	m := &Attribute{Sets: make(map[int32]int64)}
	m.Sets[message.FIGHT_ATTR_HP.Int32()] = 100
	m.Sets[message.FIGHT_ATTR_ATT.Int32()] = int64(rand.Intn(10) + 1)
	return m
}

func (s *Attribute) Type() ecs.ComponentType {
	return ATTRIBUTE_COMP
}

func (s *Attribute) HP() int64 {
	return s.Sets[message.FIGHT_ATTR_HP.Int32()]
}
func (s *Attribute) ATT() int64 {
	return s.Sets[message.FIGHT_ATTR_ATT.Int32()]
}

func (s *Attribute) InnerPB() *inner.FightAttr {
	return &inner.FightAttr{
		Sets: s.Sets,
	}
}

func (s *Attribute) SyncData() *message.SyncData {
	byt, err := proto.Marshal(&message.Attri{
		Sets: s.Sets,
	})
	expect.Nil(err)
	return &message.SyncData{
		Type: message.SyncDataType_ATTRIBUTE,
		Data: byt,
	}
}

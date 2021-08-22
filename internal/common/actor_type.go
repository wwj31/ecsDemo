package common

import (
	"fmt"
	"github.com/wwj31/dogactor/log"
	"regexp"
)

// 静态Actor类型,全局唯一
const (
	Login_Actor  = "Login"
	Center_Actor = "CenterWorld"
)

// 动态增删的actor,会有多个 game1、game2
const (
	Game_Actor  = "Game"
	Gate_Actor  = "Gate"
	World_Actor = "World"
	DB_Actor    = "DB"
	MYSQL_Actor = "Mysql"
)

func GameName(id int32) string {
	return fmt.Sprintf("%v%v_Actor", Game_Actor, id)
}

func WorldName(id int32) string {
	return fmt.Sprintf("%v%v_Actor", World_Actor, id)
}

func GateName(id int32) string {
	return fmt.Sprintf("%v%v_Actor", Gate_Actor, id)
}

// 匹配actor类型 按照固定格式匹配
func IsActorOf(t, actorId string) bool {
	str := fmt.Sprintf("%v([0-9]+)_Actor", t)
	match, e := regexp.MatchString(str, actorId)
	if e != nil {
		log.KV("actorId", actorId).ErrorStack(3, "error")
		return false
	}
	return match
}

func DBName(prex string, id int32) string {
	return fmt.Sprintf("%v%v%v_Actor", prex, DB_Actor, id)
}

func MysqlName(prex string, id int32) string {
	return fmt.Sprintf("%v%v%v_Actor", prex, MYSQL_Actor, id)
}

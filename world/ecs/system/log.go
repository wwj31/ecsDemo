package system

import (
	"github.com/wwj31/dogactor/log"
)

var (
	moveLog  = log.New(log.TAG_INFO_I)
	crossLog = log.New(log.TAG_INFO_I)
	fightLog = log.New(log.TAG_WARN_I)
	aiLog    = log.New(log.TAG_WARN_I)
	gridLog  = log.New(log.TAG_WARN_I)
	spawnLog = log.New(log.TAG_INFO_I)
	dealLog  = log.New(log.TAG_INFO_I)
	uiLog    = log.New(log.TAG_WARN_I)
	syncLog  = log.New(log.TAG_INFO_I)
)

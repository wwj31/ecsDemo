package interfaces

import "github.com/wwj31/ecsDemo/internal/message"

type IMsgComponent interface {
	SyncData() *message.SyncData
}

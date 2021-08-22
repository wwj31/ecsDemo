package interfaces

import "ecsDemo/internal/message"

type IMsgComponent interface {
	SyncData() *message.SyncData
}

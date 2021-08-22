package interfaces

type IEventSystem interface {
	HandlerActorEvent(event interface{})
}

package infrastructure

type EventDispatcher interface {
	Enqueue(data []byte) error
	Close()
}
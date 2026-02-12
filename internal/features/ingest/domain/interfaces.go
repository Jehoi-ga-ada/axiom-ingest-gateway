package domain

type EventDispatcher interface {
	Enqueue(data []byte) error
	Close()
}
package awakengine

import "github.com/DrJosh9000/vec"

type EventType int

const (
	EventNone = EventType(iota)
	EventMouseDown
	EventMouseUp
)

type Event struct {
	Type EventType
	Pos  vec.I2
}

package reader

import (
	"encoding/binary"
	"fmt"
)

//Event https://www.kernel.org/doc/Documentation/input/joystick-api.txt
type Event struct {
	Time   uint32
	Value  int16
	Type   uint8
	Number uint8
}

func ParseEvent(msg []byte) (Event, error) {
	if len(msg) != 8 {
		return Event{}, fmt.Errorf("expected message size is 8, got %d", len(msg))
	}
	event := Event{}

	event.Time = binary.BigEndian.Uint32(msg[0:4])

	event.Value = int16(msg[5])
	event.Value = (event.Value << 8) | int16(msg[4])

	event.Type = uint8(msg[6])
	event.Number = uint8(msg[7])

	return event, nil
}

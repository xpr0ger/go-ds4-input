package reader

import (
	"encoding/binary"
	"fmt"
)

// ErrorBufferLength incorrect buffer length
type ErrorBufferLength struct{ error }

//NewErrorBufferLength ErrorMessageLength's constructor
func NewErrorBufferLength(err error) error { return ErrorBufferLength{error: err} }

//InputEvent https://www.kernel.org/doc/Documentation/input/joystick-api.txt
type InputEvent struct {
	Time   uint32
	Value  int16
	Type   uint8
	Number uint8
}

// ParseInputEvent raw event that is emitted by /dev/input/jsX
func ParseInputEvent(buf []byte) (InputEvent, error) {
	if len(buf) != 8 {
		return InputEvent{}, NewErrorBufferLength(fmt.Errorf("expected buffer size is 8, got %d", len(buf)))
	}
	inputEvent := InputEvent{}

	inputEvent.Time = binary.BigEndian.Uint32(buf[0:4])

	inputEvent.Value = int16(buf[5])
	inputEvent.Value = (inputEvent.Value << 8) | int16(buf[4])

	inputEvent.Type = uint8(buf[6])
	inputEvent.Number = uint8(buf[7])

	return inputEvent, nil
}

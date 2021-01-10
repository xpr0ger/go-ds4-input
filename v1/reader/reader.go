package reader

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
)

type (
	// ErrorOpenPath failed to open file device descriptor
	ErrorOpenPath struct{ error }
	// ErrorSetDeadline failed to set deadline
	ErrorSetDeadline struct{ error }
	// ErrorReadInputEvent failed to read event
	ErrorReadInputEvent struct{ error }
	// ErrorParseInputEvent failed to parce input event
	ErrorParseInputEvent struct{ error }
	// ErrorReadDeadline deadline deadline exited
	ErrorReadDeadline struct{ error }
)

// NewErrorOpenPath ErrorOpenPath's constructor
func NewErrorOpenPath(err error) error { return ErrorOpenPath{error: err} }

// NewErrorSetDeadline ErrorSetDeadline's constructor
func NewErrorSetDeadline(err error) error { return ErrorSetDeadline{error: err} }

// NewErrorReadInputEvent ErrorReadEvent's constructor
func NewErrorReadInputEvent(err error) error { return ErrorReadInputEvent{error: err} }

// NewErrorParseInputEvent ErrorParseEvent's constructor
func NewErrorParseInputEvent(err error) error { return ErrorParseInputEvent{error: err} }

// NewErrorReadDeadline ErrorReadDeadline constructor
func NewErrorReadDeadline(err error) error { return ErrorReadDeadline{error: err} }

// InputEventReader allows to read /dev/input/jsX events
type InputEventReader struct {
	deviceDescriptor IOReaderWithDeadline
	readTimeout      time.Duration
	eventsCh         chan InputEvent
	errorsCh         chan error
}

//IOReaderWithDeadline provider reader interface with deadline setter
type IOReaderWithDeadline interface {
	io.Reader
	SetReadDeadline(deadline time.Time) error
}

// NewInputEventReader InputEventReader's constructor
func NewInputEventReader(deviceDescriptor IOReaderWithDeadline, readTimeout time.Duration) *InputEventReader {
	inputEventReader := &InputEventReader{
		readTimeout:      readTimeout,
		deviceDescriptor: deviceDescriptor,
		eventsCh:         make(chan InputEvent),
		errorsCh:         make(chan error),
	}

	return inputEventReader
}

// Close Closes all open channels and descriptors to avoid memory leaks call this function at the end of the program
func (i *InputEventReader) Close() error {
	close(i.eventsCh)
	close(i.errorsCh)
	return nil
}

// GetChannels provides the input events channel and errors channel. A user must read from both channels to prevent an
// event processing block. Before use channels - call the ListenEvents function
func (i *InputEventReader) GetChannels() (<-chan InputEvent, <-chan error) {
	return i.eventsCh, i.errorsCh
}

// ListenEvents starts to listen to the input events
func (i *InputEventReader) ListenEvents(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			//TODO: logging
			fmt.Println("listening to input events was stopped, cause: ", r)
		}
	}()

	// Buffer for the input events
	buf := make([]byte, 8)

eventLoop:
	for {
		select {
		case <-ctx.Done():
			break eventLoop
		default:
			//Sets a deadline for the read operation to prevent blocking on read, it requires to poll ctx.Done channel
			err := i.deviceDescriptor.SetReadDeadline(time.Now().Add(i.readTimeout))
			if err != nil {
				i.errorsCh <- errors.Wrap(NewErrorSetDeadline(err), "failed to set read deadline")
				continue
			}

			_, err = i.deviceDescriptor.Read(buf)
			if err != nil {
				switch err.(type) {
				case *os.PathError: // Check for read deadline
					err = NewErrorReadDeadline(err)
				default:
					err = NewErrorReadInputEvent(err)
				}

				i.errorsCh <- errors.Wrap(err, "failed to read event")
				continue
			}

			event, err := ParseInputEvent(buf)
			if err != nil {
				i.errorsCh <- errors.Wrap(NewErrorParseInputEvent(err), "failed to parse message")
				continue
			}

			i.eventsCh <- event
		}
	}
}

package reader

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
)

func NewEventReader(ctx context.Context, devicePath string, events chan Event, errorsCh chan error) error {
	fp, err := os.Open(devicePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", devicePath)
	}

	go readInputEvents(ctx, fp, events, errorsCh)

	return nil
}

func readInputEvents(ctx context.Context, inputSocket *os.File, events chan Event, errorsCh chan error) {
	buf := make([]byte, 8)
eventLoop:
	for {
		select {
		case <-ctx.Done():
			inputSocket.Close()
			break eventLoop
		default:
			err := inputSocket.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
			if err != nil {
				errorsCh <- errors.Wrap(err, "failed to set read deadline")
				continue
			}

			_, err = inputSocket.Read(buf)

			if err != nil {
				errorsCh <- errors.Wrap(err, "failed to read event")
				continue
			}

			event, err := ParseEvent(buf)
			if err != nil {
				errorsCh <- errors.Wrap(err, "failed to parse message")
				continue
			}

			events <- event
		}
	}

}

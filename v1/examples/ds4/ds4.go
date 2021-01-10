package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/xpr0ger/go-js-input/v1/controller"
	"github.com/xpr0ger/go-js-input/v1/ds4"
	"github.com/xpr0ger/go-js-input/v1/reader"
)

func main() {

	// Create context with deadline
	ctx, cancelFn := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	// Open input device
	fp, err := os.Open("/dev/input/js0")
	if err != nil {
		println(err.Error())
		return
	}
	defer fp.Close()

	// Creates input event reader
	inputEventReader := reader.NewInputEventReader(fp, time.Duration(time.Millisecond*10))

	// Create gamepad specific event's adapter
	ds4EventAdapter := ds4.NewDS4EventAdapter(inputEventReader)

	// Starts listening to input events
	defer ds4EventAdapter.Close()

	// Starts listening to input events
	go ds4EventAdapter.ListenEvents(ctx)

	// Processing CTRL+C signal
	signalCh := make(chan os.Signal)
	defer close(signalCh)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Gets channels with events and errors. To prevent blocking, both channels must be read
	gamepadEventsCh, errorsCh := ds4EventAdapter.GetChannels()

mainLoop:
	for {
		select {
		// Reads gamepad events
		case event, open := <-gamepadEventsCh:
			if !open {
				continue
			}

			fmt.Printf("%#v\n", event)
		// Reads error events
		case err, open := <-errorsCh:
			if !open {
				continue
			}
			cause := errors.Cause(err)
			switch cause.(type) {
			// Ignores some of error types
			case reader.ErrorReadDeadline, controller.ErrorUnknownGamePadEvent:
				continue
			}
			cancelFn()
		case <-ctx.Done():
			break mainLoop
		case <-signalCh:
			cancelFn()
			break mainLoop
		}
	}
}

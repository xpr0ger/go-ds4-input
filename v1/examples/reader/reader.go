package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
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

	// Initializing event reader
	inputEventReader := reader.NewInputEventReader(fp, time.Duration(time.Millisecond*10))

	// Closes event reader's channels and file descriptors at the end
	defer inputEventReader.Close()

	// Starts listening to input events
	go inputEventReader.ListenEvents(ctx)

	// Processing CTRL+C signal
	signalCh := make(chan os.Signal)
	defer close(signalCh)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Gets channels with events and errors. To prevent blocking, both channels must be read
	eventsCh, errorsCh := inputEventReader.GetChannels()

mainLoop:
	for {
		select {
		// Reads input events
		case event, open := <-eventsCh:
			if !open {
				continue
			}

			fmt.Printf("%#v\n", event)
		// Read error events
		case err, open := <-errorsCh:
			if !open {
				continue
			}
			cause := errors.Cause(err)
			switch cause.(type) {
			// Ignores read deadline errors
			case reader.ErrorReadDeadline:
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

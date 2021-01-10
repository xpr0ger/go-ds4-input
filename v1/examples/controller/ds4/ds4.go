package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xpr0ger/go-js-input/v1/controller"
	"github.com/xpr0ger/go-js-input/v1/ds4"
	"github.com/xpr0ger/go-js-input/v1/reader"
)

func main() {

	// Create context with deadline
	ctx, cancelFn := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	//Open input device
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

	// Creates gamepad state
	gamepadState := controller.NewGampadState(ds4EventAdapter)

	// Closes all related channels
	defer gamepadState.Close()

	// Starts listening to input events
	go gamepadState.ListenEvents(ctx)

	// Processing CTRL+C signal
	signalCh := make(chan os.Signal)
	defer close(signalCh)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

mainLoop:
	for {
		select {
		// Reads update events
		case _, open := <-gamepadState.GetEventUpdatedChannel():
			if !open {
				continue
			}

			fmt.Printf("%#v\n", gamepadState.GetLastEvent())
			// Emulate event's read lag
			time.Sleep(time.Second)

		case <-ctx.Done():
			break mainLoop
		case <-signalCh:
			cancelFn()
			break mainLoop
		}
	}
}

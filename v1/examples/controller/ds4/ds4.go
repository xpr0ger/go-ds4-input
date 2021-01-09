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
	events := make(chan reader.Event)
	errorsCh := make(chan error)

	ctx, cancelFn := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	err := reader.NewEventReader(ctx, "/dev/input/js0", events, errorsCh)
	if err != nil {
		println(err.Error())
		return
	}

	ds4Events := ds4.NewDS4EventAdapter(ctx, events, errorsCh)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	state := controller.NewGampadState(ctx, ds4Events)

mainLoop:
	for {
		select {
		case <-state.GetEventUpdatedChannel():
			fmt.Println(state.GetLastEvent())
			time.Sleep(time.Second)
		case err := <-errorsCh:
			cause := errors.Cause(err)
			switch cause.(type) {
			case *os.PathError:
				continue
			case controller.UnknownGamePadEvent:
				continue
			}

			println(err)
			cancelFn()
		case <-ctx.Done():
			break mainLoop
		case <-signalCh:
			cancelFn()
			break mainLoop
		}
	}

	close(events)
	close(errorsCh)
	close(signalCh)
}

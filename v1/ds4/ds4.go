package ds4

import (
	"context"
	"fmt"

	"github.com/xpr0ger/go-js-input/v1/controller"
	"github.com/xpr0ger/go-js-input/v1/reader"
)

//InputEventProvider provides interface to get events channels
type InputEventProvider interface {
	ListenEvents(ctx context.Context)
	GetChannels() (<-chan reader.InputEvent, <-chan error)
	Close() error
}

//DS4EventAdapter input events adapter for the DualShock 4 gamepad
type DS4EventAdapter struct {
	eventsProvider InputEventProvider

	dPadUp    bool
	dPadRight bool
	dPadDown  bool
	dPadLeft  bool

	gamePadEventsCh chan controller.GamePadEvent
	errorsCh        chan error
}

//NewDS4EventAdapter DS4EventAdapter's constructor
func NewDS4EventAdapter(eventsProvider InputEventProvider) *DS4EventAdapter {
	adapter := &DS4EventAdapter{
		eventsProvider:  eventsProvider,
		gamePadEventsCh: make(chan controller.GamePadEvent),
		errorsCh:        make(chan error),
	}
	return adapter
}

// Close Closes all open channels and descriptors to avoid memory leaks call this function at the end of the program
func (d *DS4EventAdapter) Close() error {
	close(d.gamePadEventsCh)
	close(d.errorsCh)
	return d.eventsProvider.Close()
}

// GetChannels provides the gamepad specific input events channel and errors channel. A user must read from both channels to prevent an
// event processing block. Before use channels - call the ListenEvents function
func (d *DS4EventAdapter) GetChannels() (<-chan controller.GamePadEvent, <-chan error) {
	return d.gamePadEventsCh, d.errorsCh
}

// ListenEvents starts to listen to the input events and convert them to gamepad specific
func (d *DS4EventAdapter) ListenEvents(ctx context.Context) {
	go d.eventsProvider.ListenEvents(ctx)
	eventCh, errorsCh := d.eventsProvider.GetChannels()
eventLoop:
	for {
		select {
		case <-ctx.Done():
			break eventLoop
		case rawEvent, open := <-eventCh:
			if !open {
				continue
			}

			event, err := d.ConvertEvent(rawEvent)
			if err != nil {
				d.errorsCh <- err
				continue
			}
			d.gamePadEventsCh <- event
		case err, open := <-errorsCh: // Redirect errors
			if !open {
				continue
			}

			d.errorsCh <- err
		}
	}
}

//ConvertEvent convert input event to gamepad specific event
func (d *DS4EventAdapter) ConvertEvent(event reader.InputEvent) (controller.GamePadEvent, error) {
	gamePadEvent := controller.GamePadEvent{EventTime: event.Time}

	//DPad Up
	if event.Type == 2 && event.Number == 7 && event.Value < 0 {
		gamePadEvent.ButtonType = ButtonDPadUp
		gamePadEvent.Value = 1
		d.dPadUp = true
		return gamePadEvent, nil
	}

	if event.Type == 2 && event.Number == 7 && event.Value == 0 && d.dPadUp {
		gamePadEvent.ButtonType = ButtonDPadUp
		gamePadEvent.Value = 0
		d.dPadUp = false
		return gamePadEvent, nil
	}

	// DPad Right
	if event.Type == 2 && event.Number == 6 && event.Value > 0 {
		gamePadEvent.ButtonType = ButtonDPadRight
		gamePadEvent.Value = 1
		d.dPadRight = true
		return gamePadEvent, nil
	}

	if event.Type == 2 && event.Number == 6 && event.Value == 0 && d.dPadRight {
		gamePadEvent.ButtonType = ButtonDPadRight
		gamePadEvent.Value = 0
		d.dPadRight = false
		return gamePadEvent, nil
	}

	// DPad Down
	if event.Type == 2 && event.Number == 7 && event.Value > 0 {
		gamePadEvent.ButtonType = ButtonDPadDown
		gamePadEvent.Value = 1
		d.dPadDown = true
		return gamePadEvent, nil
	}

	if event.Type == 2 && event.Number == 7 && event.Value == 0 && d.dPadDown {
		gamePadEvent.ButtonType = ButtonDPadDown
		gamePadEvent.Value = 0
		d.dPadDown = false
		return gamePadEvent, nil
	}

	//DPad Left
	if event.Type == 2 && event.Number == 6 && event.Value < 0 {
		gamePadEvent.ButtonType = ButtonDPadLeft
		gamePadEvent.Value = 1
		d.dPadLeft = true
		return gamePadEvent, nil
	}

	if event.Type == 2 && event.Number == 6 && event.Value == 0 && d.dPadLeft {
		gamePadEvent.ButtonType = ButtonDPadLeft
		gamePadEvent.Value = 0
		d.dPadLeft = false
		return gamePadEvent, nil
	}

	//Triangle
	if event.Type == 1 && event.Number == 2 {
		gamePadEvent.ButtonType = ButtonTriangle
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//Circle
	if event.Type == 1 && event.Number == 1 {
		gamePadEvent.ButtonType = ButtonCircle
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//Cross
	if event.Type == 1 && event.Number == 0 {
		gamePadEvent.ButtonType = ButtonCross
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//Square
	if event.Type == 1 && event.Number == 3 {
		gamePadEvent.ButtonType = ButtonSquare
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//Share
	if event.Type == 1 && event.Number == 8 {
		gamePadEvent.ButtonType = ButtonShare
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	// ButtonOption
	if event.Type == 1 && event.Number == 9 {
		gamePadEvent.ButtonType = ButtonOption
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//L1
	if event.Type == 1 && event.Number == 4 {
		gamePadEvent.ButtonType = ButtonL1
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//L2
	if event.Type == 1 && event.Number == 6 {
		gamePadEvent.ButtonType = ButtonL2
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//L2 Force
	if event.Type == 2 && event.Number == 5 {
		gamePadEvent.ButtonType = ButtonL2Force
		gamePadEvent.Value = d.fixTriggerValue(int(event.Value))
		return gamePadEvent, nil
	}

	//L3
	if event.Type == 1 && event.Number == 11 {
		gamePadEvent.ButtonType = ButtonL3
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//L3 Vertical
	if event.Type == 2 && event.Number == 1 {
		gamePadEvent.ButtonType = ButtonL3Vertical
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//L3 Horizontal
	if event.Type == 2 && event.Number == 0 {
		gamePadEvent.ButtonType = ButtonL3Horizontal
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//R1
	if event.Type == 1 && event.Number == 5 {
		gamePadEvent.ButtonType = ButtonR1
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//R2
	if event.Type == 1 && event.Number == 7 {
		gamePadEvent.ButtonType = ButtonR2
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//R2 Force
	if event.Type == 2 && event.Number == 2 {
		gamePadEvent.ButtonType = ButtonR2Force
		gamePadEvent.Value = d.fixTriggerValue(int(event.Value))
		return gamePadEvent, nil
	}

	//R3
	if event.Type == 1 && event.Number == 12 {
		gamePadEvent.ButtonType = ButtonR3
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//R3 Vertical
	if event.Type == 2 && event.Number == 4 {
		gamePadEvent.ButtonType = ButtonL3Vertical
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	//R3 Horizontal
	if event.Type == 2 && event.Number == 3 {
		gamePadEvent.ButtonType = ButtonR3Horizontal
		gamePadEvent.Value = int(event.Value)
		return gamePadEvent, nil
	}

	return gamePadEvent, controller.NewErrorUnknownGamePadEvent(fmt.Errorf("unknown gamepad event %#v", event))
}

//Max value 65534
func (d *DS4EventAdapter) fixTriggerValue(val int) int {
	correctionValue := 0xffff / 2
	return val + correctionValue

}

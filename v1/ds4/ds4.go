package ds4

import (
	"context"
	"fmt"

	"github.com/xpr0ger/go-js-input/v1/controller"
	"github.com/xpr0ger/go-js-input/v1/reader"
)

type DS4EventAdapter struct {
	dPadUp    bool
	dPadRight bool
	dPadDown  bool
	dPadLeft  bool
}

func NewDS4EventAdapter(ctx context.Context, rawEvents chan reader.Event, errorsCh chan error) chan controller.GamePadEvent {
	adapter := &DS4EventAdapter{}
	events := make(chan controller.GamePadEvent)
	go adapter.Run(ctx, rawEvents, events, errorsCh)
	return events
}

func (d *DS4EventAdapter) Run(ctx context.Context, rawEvents chan reader.Event, events chan controller.GamePadEvent, errorsCh chan error) {
eventLoop:
	for {
		select {
		case <-ctx.Done():
			break eventLoop
		case rawEvent := <-rawEvents:
			event, err := d.ConvertEvent(rawEvent)
			if err != nil {
				errorsCh <- err
				continue
			}
			events <- event
		}
	}
}

func (d *DS4EventAdapter) ConvertEvent(event reader.Event) (controller.GamePadEvent, error) {
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

	return gamePadEvent, controller.NewUnknownGamePadEvent(fmt.Errorf("unknown gamepad event %#v", event))
}

func (d *DS4EventAdapter) fixTriggerValue(val int) int {
	correctionValue := 0xffff/2 + 1
	return val + correctionValue

}

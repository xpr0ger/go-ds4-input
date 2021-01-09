package controller

import "context"

type ButtonType int

type GamePadEvent struct {
	EventTime  uint32
	ButtonType ButtonType
	Value      int
}

type UnknownGamePadEvent struct {
	error
}

func NewUnknownGamePadEvent(err error) error {
	return UnknownGamePadEvent{error: err}
}

type GamePadState struct {
	updatedCh chan struct{}
	lastEvent GamePadEvent
	state     map[ButtonType]int
}

func NewGampadState(ctx context.Context, eventsCh chan GamePadEvent) *GamePadState {
	gamePadState := &GamePadState{
		updatedCh: make(chan struct{}),
		state:     make(map[ButtonType]int),
	}

	go gamePadState.eventLoop(ctx, eventsCh)

	return gamePadState
}

func (g *GamePadState) GetEventUpdatedChannel() <-chan struct{} {
	return g.updatedCh
}

func (g *GamePadState) GetLastEvent() GamePadEvent {
	return g.lastEvent
}

func (g *GamePadState) GetButtonState(buttonType ButtonType) int {
	return g.state[buttonType]
}

func (g *GamePadState) eventLoop(ctx context.Context, eventsCh chan GamePadEvent) {
eventLoop:
	for {
		select {
		case <-ctx.Done():
			close(g.updatedCh)
			break eventLoop
		case event := <-eventsCh:
			g.lastEvent = event
			g.state[event.ButtonType] = event.Value
			select {
			case g.updatedCh <- struct{}{}:
			default:
			}
		}
	}
}

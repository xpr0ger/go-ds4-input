package controller

import (
	"context"
	"sync"
)

type ButtonType int

// GamepadEventProvider provides interface to get gamepad specific events channels
type GamepadEventProvider interface {
	ListenEvents(ctx context.Context)
	GetChannels() (<-chan GamePadEvent, <-chan error)
	Close() error
}

// GamePadEvent describes singe gamepad's event
type GamePadEvent struct {
	EventTime  uint32
	ButtonType ButtonType
	Value      int
}

// ErrorUnknownGamePadEvent emits on unknown gamepad event
type ErrorUnknownGamePadEvent struct {
	error
}

// NewErrorUnknownGamePadEvent - ErrorUnknownGamePadEvent's constructor
func NewErrorUnknownGamePadEvent(err error) error { return ErrorUnknownGamePadEvent{error: err} }

// GamePadState stores a current gamepad state
type GamePadState struct {
	lastError            error
	gamepadEventProvider GamepadEventProvider
	updatedCh            chan struct{}
	lock                 sync.RWMutex
	lastEvent            GamePadEvent
	state                map[ButtonType]int
}

// NewGampadState GampadState's constructor
func NewGampadState(gamepadEventProvider GamepadEventProvider) *GamePadState {
	gamePadState := &GamePadState{
		gamepadEventProvider: gamepadEventProvider,
		updatedCh:            make(chan struct{}),
		state:                make(map[ButtonType]int),
	}

	return gamePadState
}

// GetEventUpdatedChannel on every update this channel gets notification
func (g *GamePadState) GetEventUpdatedChannel() <-chan struct{} {
	return g.updatedCh
}

// GetLastEvent gets the latest gamepad event
func (g *GamePadState) GetLastEvent() GamePadEvent {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.lastEvent
}

// GetButtonState get specific button state see ds4/buttons.go
func (g *GamePadState) GetButtonState(buttonType ButtonType) int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.state[buttonType]
}

// Close Closes all open channels and descriptors to avoid memory leaks call this function at the end of the program
func (g *GamePadState) Close() error {
	close(g.updatedCh)
	return g.gamepadEventProvider.Close()
}

// ListenEvents starts to listen to the input events
func (g *GamePadState) ListenEvents(ctx context.Context) {
	go g.gamepadEventProvider.ListenEvents(ctx)
	gamepadEventsCh, errorsCh := g.gamepadEventProvider.GetChannels()

eventLoop:
	for {
		select {
		case <-ctx.Done():
			break eventLoop
		case event, open := <-gamepadEventsCh:
			if !open {
				continue
			}
			g.lock.Lock()
			g.lastEvent = event
			g.state[event.ButtonType] = event.Value
			g.lock.Unlock()

			select {
			case g.updatedCh <- struct{}{}:
			default:
			}
		case err, open := <-errorsCh:
			if !open {
				continue
			}

			g.lastError = err
		}
	}
}

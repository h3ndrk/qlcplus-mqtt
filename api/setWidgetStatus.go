package api

import (
	"fmt"
	"strconv"
)

// SetWidgetStatusButton sets the state of the button widget with the given `<id>` to the given `<pressed>` state
func (a *API) SetWidgetStatusButton(id int, pressed bool) error {
	if pressed {
		return a.setWidgetStatus(id, "1")
	}

	return a.setWidgetStatus(id, "0")
}

// SetWidgetStatusSlider sets the state of the slider widget with the given `<id>` to the given `<value>` state
func (a *API) SetWidgetStatusSlider(id int, value int) error {
	return a.setWidgetStatus(id, strconv.Itoa(value))
}

// SetWidgetStatusAudioTrigger sets the state of the audio trigger widget with the given `<id>` to the given `<active>` state
func (a *API) SetWidgetStatusAudioTrigger(id int, active bool) error {
	if active {
		return a.setWidgetStatus(id, "1")
	}

	return a.setWidgetStatus(id, "0")
}

// SetWidgetStatusCueListPlay sets the state of the cue list widget with the given `<id>` to the "PLAY" state
func (a *API) SetWidgetStatusCueListPlay(id int) error {
	return a.setWidgetStatus(id, "PLAY")
}

// SetWidgetStatusCueListStop sets the state of the cue list widget with the given `<id>` to the "STOP" state
func (a *API) SetWidgetStatusCueListStop(id int) error {
	return a.setWidgetStatus(id, "STOP")
}

// SetWidgetStatusCueListPrev sets the state of the cue list widget with the given `<id>` to the "PREV" state
func (a *API) SetWidgetStatusCueListPrev(id int) error {
	return a.setWidgetStatus(id, "PREV")
}

// SetWidgetStatusCueListNext sets the state of the cue list widget with the given `<id>` to the "NEXT" state
func (a *API) SetWidgetStatusCueListNext(id int) error {
	return a.setWidgetStatus(id, "NEXT")
}

// SetWidgetStatusCueListStep sets the state of the cue list widget with the given `<id>` to the "STEP" state
func (a *API) SetWidgetStatusCueListStep(id int, index int) error {
	return a.setWidgetStatus(id, fmt.Sprintf("STEP|%d", index))
}

// SetWidgetStatusFramePrevPage sets the state of the frame/solo frame widget with the given `<id>` to the "PREV_PG" state
func (a *API) SetWidgetStatusFramePrevPage(id int) error {
	return a.setWidgetStatus(id, "PREV_PG")
}

// SetWidgetStatusFrameNextPage sets the state of the frame/solo frame widget with the given `<id>` to the "NEXT_PG" state
func (a *API) SetWidgetStatusFrameNextPage(id int) error {
	return a.setWidgetStatus(id, "NEXT_PG")
}

// setWidgetStatus sets the state of the function with the given `<id>` to the given `<state>`
func (a *API) setWidgetStatus( /*ctx context.Context,*/ id int, value string) error {
	// - [ ] `<id>[|<value>]`: Sets given `<value>` of widget with given `<id>`
	//     - `Button`: `<value>` equals to `1` for press and `0` for release
	//     - `Slider`: `<value>` is in range of slider
	//     - `AudioTrigger`: `<value>` equals to `1` for active and `0` for inactive
	//     - `CueList`: `<value>` is `<command>[|<index>]`:
	//         - `PLAY`: Plays cue
	//         - `STOP`: Stops cue
	//         - `PREV`: Selects previous cue
	//         - `NEXT`: Selects next cue
	//         - `STEP|<index>`: Selects cue with given `<index>`
	//     - `Frame`, `SoloFrame`: `<value>` is one of `NEXT_PG` (next page) or `PREV_PG` (previous page)

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|setFunctionStatus|%d|%s", id, value)))
	if err != nil {
		return err
	}

	return nil
}

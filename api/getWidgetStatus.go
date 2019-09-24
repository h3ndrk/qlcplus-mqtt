package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// WidgetStatus represents the status of a widget
type WidgetStatus struct {
	ButtonActive     bool // only valid if widget type is `Button`: whether button is active
	ButtonMonitoring bool // only valid if widget type is `Button`: whether button is in monitoring state
	SliderValue      int  // only valid if widget type is `Slider`: slider value in range `[0,255]`
	CueListIndex     int  // only valid if widget type is `CueList`: current cue list index
	CueListPlaying   bool // only valid if widget type is `CueList`: whether cue list is playing currently
}

// GetWidgetStatus returns the state of the widget with the given `<id>`
func (a *API) GetWidgetStatus( /*ctx context.Context,*/ id int) (WidgetStatus, error) {
	// - `QLC+API|getWidgetStatus|<id>`: Returns the state of the widget with the given `<id>`
	//     - answer for type `Button`: `QLC+API|getWidgetStatus|<value>` with `<value>` equals to `255` when active, `127` when monitoring or `0` when inactive
	//     - answer for type `Slider`: `QLC+API|getWidgetStatus|<value>` with `<value>` in the range `[0,255]`
	//     - answer for type playing `CueList`: `QLC+API|getWidgetStatus|PLAY|<cueIndex>` with `<cueIndex>` equals to the current cue index
	//     - answer for type stopped `CueList`: `QLC+API|getWidgetStatus|STOP`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|getWidgetStatus|%d", id)))
	if err != nil {
		return WidgetStatus{}, err
	}

	// Receive message
	msg, err := a.receiveText()
	if err != nil {
		return WidgetStatus{}, err
	}

	// Unmarshal message
	parts := strings.Split(string(msg), "|")
	if len(parts) < 2 {
		return WidgetStatus{}, errors.New("Invalid amount of parts")
	}
	if parts[0] != "QLC+API" || parts[1] != "getWidgetStatus" {
		return WidgetStatus{}, errors.New("Unexpected response")
	}

	s := WidgetStatus{}

	// `Button` type
	if len(parts) == 3 {
		if parts[2] == "255" {
			s.ButtonActive = true
		} else if parts[2] == "127" {
			s.ButtonMonitoring = true
		}
	}
	// `Slider` type
	if len(parts) == 3 {
		value, err := strconv.Atoi(parts[2])
		if err == nil {
			s.SliderValue = value
		}
	}
	// `CueList` type
	if len(parts) == 4 && parts[2] == "PLAY" {
		value, err := strconv.Atoi(parts[3])
		if err == nil {
			s.CueListIndex = value
			s.CueListPlaying = true
		}
	}

	return s, nil
}

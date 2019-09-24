package api

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// WidgetType is one of `Scene`, `Chaser`, `EFX`, `Collection`, `Script`, `RGBMatrix`, `Show`, `Sequence`, `Audio`, `Video`
type WidgetType string

// GetWidgetType returns the type of the widget with the given `<id>`
func (a *API) GetWidgetType( /*ctx context.Context,*/ id int) (WidgetType, error) {
	// - `QLC+API|getWidgetType|<id>`: Returns the type of the function with the given `<id>` (or `Unknown`)
	//     - answer: `QLC+API|getWidgetType|<type>` with `<type>` one of `Button`, `Slider`, `XYPad`, `Frame`, `SoloFrame`, `SpeedDial`, `CueList`, `Label`, `AudioTriggers`, `Animation`, `Clock`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|getWidgetType|%d", id)))
	if err != nil {
		return "", err
	}

	// Receive message
	msg, err := a.receiveText()
	if err != nil {
		return "", err
	}

	// Unmarshal message
	parts := strings.Split(string(msg), "|")
	if len(parts) < 3 {
		return "", errors.New("Invalid amount of parts")
	}
	if parts[0] != "QLC+API" || parts[1] != "getWidgetType" {
		return "", errors.New("Unexpected response")
	}
	if parts[2] == "Unknown" {
		return "", errors.New("Undefined response received")
	}

	return WidgetType(parts[2]), nil
}

package api

import (
	"strings"

	"github.com/pkg/errors"
)

// IsProjectLoaded queries and clears/disables project loaded state
func (a *API) IsProjectLoaded( /*ctx context.Context,*/ ) (bool, error) {
	// - `QLC+API|isProjectLoaded`: Query and clear/disable project loaded state
	//     - answer: `QLC+API|isProjectLoaded|<state>` with `<state>` equals to `true` or `false`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte("QLC+API|isProjectLoaded"))
	if err != nil {
		return false, err
	}

	// Receive message
	msg, err := a.receiveText()
	if err != nil {
		return false, err
	}

	// Unmarshal message
	parts := strings.Split(string(msg), "|")
	if len(parts) < 3 {
		return false, errors.New("Invalid amount of parts")
	}
	if parts[0] != "QLC+API" || parts[1] != "isProjectLoaded" {
		return false, errors.New("Unexpected response")
	}

	return parts[2] == "true", nil
}

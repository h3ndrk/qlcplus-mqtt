package api

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// GetWidgetsNumber returns the number of widgets
func (a *API) GetWidgetsNumber( /*ctx context.Context,*/ ) (int, error) {
	// - `QLC+API|getWidgetsNumber`: Returns the number of widgets
	//     - answer: `QLC+API|getWidgetsNumber|<number>`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte("QLC+API|getWidgetsNumber"))
	if err != nil {
		return 0, err
	}

	// Receive message
	msg, err := a.receiveText()
	if err != nil {
		return 0, err
	}

	// Unmarshal message
	parts := strings.Split(string(msg), "|")
	if len(parts) < 3 {
		return 0, errors.New("Invalid amount of parts")
	}
	if parts[0] != "QLC+API" || parts[1] != "getWidgetsNumber" {
		return 0, errors.New("Unexpected response")
	}

	number, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert number of widgets")
	}

	return number, nil
}

package api

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// FunctionStatus is true when the function is running, false otherwise
type FunctionStatus bool

// GetFunctionStatus returns the state of the function with the given `<id>`
func (a *API) GetFunctionStatus( /*ctx context.Context,*/ id int) (FunctionStatus, error) {
	// - `QLC+API|getFunctionStatus|<id>`: Returns the state of the function with the given `<id>`
	//     - answer: `QLC+API|getFunctionStatus|<state>` with `<state>` equals to `Running` or `Stopped` (or `Undefined`)

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|getFunctionStatus|%d", id)))
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
	if parts[0] != "QLC+API" || parts[1] != "getFunctionStatus" {
		return false, errors.New("Unexpected response")
	}
	if parts[2] == "Undefined" {
		return false, errors.New("Undefined response received")
	}

	return parts[2] == "Running", nil
}

package api

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// FunctionType is one of `Scene`, `Chaser`, `EFX`, `Collection`, `Script`, `RGBMatrix`, `Show`, `Sequence`, `Audio`, `Video`
type FunctionType string

// GetFunctionType returns the type of the function with the given `<id>` (or `Undefined`)
func (a *API) GetFunctionType( /*ctx context.Context,*/ id int) (FunctionType, error) {
	// - `QLC+API|getFunctionType|<id>`: Returns the type of the function with the given `<id>` (or `Undefined`)
	//     - answer: `QLC+API|getFunctionType|<type>` with `<type>` one of `Scene`, `Chaser`, `EFX`, `Collection`, `Script`, `RGBMatrix`, `Show`, `Sequence`, `Audio`, `Video`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|getFunctionType|%d", id)))
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
	if parts[0] != "QLC+API" || parts[1] != "getFunctionType" {
		return "", errors.New("Unexpected response")
	}
	if parts[2] == "Undefined" {
		return "", errors.New("Undefined response received")
	}

	return FunctionType(parts[2]), nil
}

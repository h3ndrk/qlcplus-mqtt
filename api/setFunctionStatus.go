package api

import (
	"fmt"
)

// SetFunctionStatus sets the state of the function with the given `<id>` to the given `<state>`
func (a *API) SetFunctionStatus( /*ctx context.Context,*/ id int, status FunctionStatus) error {
	// - `QLC+API|setFunctionStatus|<id>|<state>`: Sets the state of the function with the given `<id>` to the given `<state>` (`0` is `Stopped`, `1` is `Running`)

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	state := 0
	if status {
		state = 1
	}

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|setFunctionStatus|%d|%d", id, state)))
	if err != nil {
		return err
	}

	return nil
}

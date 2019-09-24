package api

import (
	"fmt"
)

// SetChannelValue sets the channel described with given `<absAddr>` (absolute address) to given `<value>` in range `[0,255]`
func (a *API) SetChannelValue( /*ctx context.Context,*/ absAddr int, value int) error {
	// - `CH|<absAddr>|<value>`: Sets the channel described with given `<absAddr>` to given `<value>` in range `[0,255]`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("CH|%d|%d", absAddr, value)))
	if err != nil {
		return err
	}

	return nil
}

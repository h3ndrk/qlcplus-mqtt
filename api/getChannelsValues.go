package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ChannelValue represents one channel
type ChannelValue struct {
	Addr  int
	Value int
	Type  string
}

// ChannelsValues is the list of channels
type ChannelsValues []ChannelValue

// GetChannelsValues returns the channel values of a given `<universe>` (starting at zero) beginning at given `<startAddr>` (starting at zero) and returning one or given `<count>` channels
func (a *API) GetChannelsValues( /*ctx context.Context,*/ universe int, startAddr int, count int) (ChannelsValues, error) {
	// - `QLC+API|getChannelsValues|<universe>|<startAddr>|[<count>]`: Returns the channel values of a given `<universe>` (starting at one) beginning at given `<startAddr>` (starting at one) and returning one or given `<count>` channels
	//     - answer: `QLC+API|getChannelsValues|[<channel>|<value>|<type>]...` with `<channel>` starting at one, value in range `[0,255]` and `<type>` the fixture type

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte(fmt.Sprintf("QLC+API|getChannelsValues|%d|%d|%d", universe+1, startAddr+1, count)))
	if err != nil {
		return nil, err
	}

	// Receive message
	msg, err := a.receiveText()
	if err != nil {
		return nil, err
	}

	// Unmarshal message
	parts := strings.Split(string(msg), "|")
	if len(parts) < 2 || (len(parts)-2)%3 != 0 {
		return nil, errors.New("Invalid amount of parts")
	}
	if parts[0] != "QLC+API" || parts[1] != "getChannelsValues" {
		return nil, errors.New("Unexpected response")
	}

	c := ChannelsValues{}
	for i := 2; i < len(parts); i += 3 {
		Addr, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to convert part #%d", i)
		}
		Value, err := strconv.Atoi(parts[i+1])
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to convert part #%d", i+1)
		}

		c = append(c, ChannelValue{Addr: Addr, Value: Value, Type: parts[i+2]})
	}

	return c, nil
}

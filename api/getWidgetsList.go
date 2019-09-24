package api

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// WidgetIDCaption represents one widget
type WidgetIDCaption struct {
	ID      int
	Caption string
}

// WidgetsList is the list of widgets
type WidgetsList []WidgetIDCaption

// GetWidgetsList returns all widgets
func (a *API) GetWidgetsList( /*ctx context.Context,*/ ) (WidgetsList, error) {
	// - `QLC+API|getWidgetsList`: Returns all widgets
	//     - answer: `QLC+API|getWidgetsList|[<id>|<caption>]...`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte("QLC+API|getWidgetsList"))
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
	if len(parts) < 2 || (len(parts)-2)%2 != 0 {
		return nil, errors.New("Invalid amount of parts")
	}
	if parts[0] != "QLC+API" || parts[1] != "getWidgetsList" {
		return nil, errors.New("Unexpected response")
	}

	f := WidgetsList{}
	for i := 2; i < len(parts); i += 2 {
		ID, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to convert part #%d", i)
		}

		f = append(f, WidgetIDCaption{ID: ID, Caption: parts[i+1]})
	}

	return f, nil
}

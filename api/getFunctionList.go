package api

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// FunctionIDName represents one function
type FunctionIDName struct {
	ID   int
	Name string
}

// FunctionList is the list of functions
type FunctionList []FunctionIDName

// GetFunctionList returns all functions
func (a *API) GetFunctionList( /*ctx context.Context,*/ ) (FunctionList, error) {
	// - `QLC+API|getFunctionsList`: Returns all functions
	//     - answer: `QLC+API|getFunctionsList|[<id>|<name>]...`

	// lock API
	a.lock.Lock()
	defer a.lock.Unlock()

	// Send message
	err := a.writeText([]byte("QLC+API|getFunctionsList"))
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

	f := FunctionList{}
	for i := 2; i < len(parts); i += 2 {
		ID, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to convert part #%d", i)
		}

		f = append(f, FunctionIDName{ID: ID, Name: parts[i+1]})
	}

	return f, nil
}

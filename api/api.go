package api

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// API represents an API connection to QLC+
type API struct {
	lock sync.Mutex // API lock (allow only one concurrent call)
	ws   *websocket.Conn
}

// NewAPI constructs a new API connection to QLC+
func NewAPI(wsURL string) (*API, error) {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to websocket: %s", wsURL)
	}

	return &API{ws: ws}, nil
}

// Close disconnects an API connection to QLC+
func (a *API) Close() {
	a.ws.Close()
}

func (a *API) receiveText() ([]byte, error) {
	t, msg, err := a.ws.ReadMessage()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read websocket message")
	}
	if t != websocket.TextMessage {
		return nil, errors.Errorf("Unexpected websocket message type: %d", t)
	}

	return msg, nil
}

func (a *API) writeText(msg []byte) error {
	return a.ws.WriteMessage(websocket.TextMessage, msg)
}

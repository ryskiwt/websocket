package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Conn represents a WebSocket connection.
// It is a thin wrapper of gorilla/websocket.Conn implementing net.Conn.
type Conn struct {
	*websocket.Conn
}

// Dial opens a new client connection to a WebSocket.
func Dial(url string, header http.Header) (ws *Conn, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Conn: conn,
	}, nil
}

// WrapUpgrade wraps gorilla/websocket.Upgrade.
func WrapUpgrade(conn *websocket.Conn, err error) (*Conn, error) {
	return &Conn{
		Conn: conn,
	}, err
}

// Read reads data from the connection.
func (ws *Conn) Read(bs []byte) (n int, err error) {
	_, rd, err := ws.NextReader()
	if err != nil {
		return 0, err
	}
	return rd.Read(bs)
}

// Write writes data to the connection.
func (ws *Conn) Write(bs []byte) (n int, err error) {
	if err := ws.WriteMessage(websocket.BinaryMessage, bs); err != nil {
		return 0, err
	}
	return len(bs), nil
}

// SetDeadline sets the read and write deadlines associated with the connection.
func (ws *Conn) SetDeadline(t time.Time) error {
	if err := ws.SetReadDeadline(t); err != nil {
		return err
	}
	if err := ws.SetWriteDeadline(t); err != nil {
		return err
	}
	return nil
}

package websocket

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Conn represents a WebSocket connection.
// It is a thin wrapper of gorilla/websocket.Conn implementing net.Conn.
type Conn struct {
	*websocket.Conn
	pr     io.ReadCloser
	pw     io.WriteCloser
	ctx    context.Context
	cancel context.CancelFunc
}

// Dial opens a new client connection to a WebSocket.
func Dial(url string, header http.Header) (ws *Conn, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}
	return newConn(conn), nil
}

// WrapUpgrade wraps gorilla/websocket.Upgrade.
func WrapUpgrade(conn *websocket.Conn, err error) (*Conn, error) {
	return newConn(conn), err
}

// Read reads data from the connection.
func (ws *Conn) Read(bs []byte) (n int, err error) {
	return ws.pr.Read(bs)
}

// Write writes data to the connection.
func (ws *Conn) Write(bs []byte) (n int, err error) {
	if err := ws.WriteMessage(websocket.BinaryMessage, bs); err != nil {
		return 0, err
	}
	return len(bs), nil
}

// Close closes connection.
func (ws *Conn) Close() error {
	ws.cancel()
	if err := ws.Conn.CloseHandler()(websocket.CloseNormalClosure, ""); err != nil {
		return err
	}
	return ws.Conn.Close()
}

// SetDeadline sets the read and write deadlines associated with the connection.
func (ws *Conn) SetDeadline(t time.Time) error {
	if err := ws.SetWriteDeadline(t); err != nil {
		return err
	}
	if err := ws.SetReadDeadline(t); err != nil {
		return err
	}
	return nil
}

func newConn(conn *websocket.Conn) *Conn {
	ctx, cancel := context.WithCancel(context.Background())
	pr, pw := io.Pipe()
	wsc := Conn{
		Conn:   conn,
		pw:     pw,
		pr:     pr,
		ctx:    ctx,
		cancel: cancel,
	}

	// receive routine
	go func() {
		for {
			_, rd, err := conn.NextReader()
			if err != nil {
				select {
				case <-ctx.Done():
					pw.Close()
					return
				default:
				}
				pw.CloseWithError(err)
				return
			}
			if _, err := io.Copy(pw, rd); err != nil {
				select {
				case <-ctx.Done():
					pw.Close()
					return
				default:
				}
				pw.CloseWithError(err)
				return
			}
		}
	}()

	return &wsc
}

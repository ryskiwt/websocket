package websocket_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ryskiwt/websocket"
)

func TestConn_ReadWrite(t *testing.T) {

	//
	// Dial
	//

	conn, err := websocket.Dial("ws://127.0.0.1:9001/echo", nil)
	if err != nil {
		t.Fatalf("%#v\n", err)
	}

	msg := []byte("test message")
	if _, err := conn.Write(msg); err != nil {
		t.Fatalf("%#v\n", err)
	}

	bs := make([]byte, 1024)
	n, err := conn.Read(bs)
	if err != nil && err != io.EOF {
		t.Fatalf("%#v\n", err)
	}

	if exp, act := msg, bs[:n]; !bytes.Equal(exp, act) {
		t.Fatalf("exp=%#v, act=%#v\n", exp, act)
	}
}

func TestConn_ReadAfterClose(t *testing.T) {

	conn, err := websocket.Dial("ws://127.0.0.1:9001/echo", nil)
	if err != nil {
		t.Fatalf("%#v\n", err)
	}

	if err := conn.Close(); err != nil {
		t.Errorf("%#v\n", err)
	}

	bs := make([]byte, 1024)
	if _, err = conn.Read(bs); err != nil && err != io.EOF {
		t.Errorf("%#v\n", err)
	}
}

func TestConn_WriteAfterClose(t *testing.T) {

	conn, err := websocket.Dial("ws://127.0.0.1:9001/echo", nil)
	if err != nil {
		t.Fatalf("%#v\n", err)
	}

	if err := conn.Close(); err != nil {
		t.Errorf("%#v\n", err)
	}

	msg := []byte("test message")
	if _, err := conn.Write(msg); err == nil {
		t.Fatalf("%#v\n", err)
	}
}

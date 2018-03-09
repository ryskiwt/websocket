# WebSocket

A thin wrapper of [gorilla/websocket](https://github.com/gorilla/websocket). This library just adds net.Conn interface to gorilla/websocket.Conn.

[![GoDoc](https://godoc.org/github.com/ryskiwt/websocket?status.svg)](https://godoc.org/github.com/ryskiwt/websocket)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://github.com/ryskiwt/websocket/blob/master/LICENSE)

### Usage

```go
package main

import (
	"log"
	"net/http"

	gorilla "github.com/gorilla/websocket"
	"github.com/ryskiwt/websocket"
)

func main() {

	var conn *websocket.Conn
	var err error

	conn, err = websocket.Dial("ws://127.0.0.1/socket", nil)
	if err != nil {
		log.Fatalf("%#v\n", err)
	}
	defer conn.Close()

	var w http.ResponseWriter
	var r *http.Request
	conn, err = websocket.WrapUpgrade(gorilla.Upgrade(w, r, nil, 1024, 1024))
	if err != nil {
		log.Fatalf("%#v\n", err)
	}
	defer conn.Close()

}
```

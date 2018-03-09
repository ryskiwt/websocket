# websocket

A thin wrapper of websocket/gorilla implementing net.Conn

# Usage

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

package websocket_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"
	"github.com/ryskiwt/websocket"
)

func TestMain(m *testing.M) {

	// initialize

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", handlerFunc)
	s := http.Server{
		Addr:    "127.0.0.1:9001",
		Handler: mux,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	<-time.After(200 * time.Millisecond)

	code := m.Run()

	// clean up
	s.Shutdown(context.Background())

	os.Exit(code)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {

	upgrader := gorilla.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
	}

	//
	// WrapUpgrade
	//

	conn, err := websocket.WrapUpgrade(upgrader.Upgrade(w, r, nil))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		if _, err := io.Copy(conn, conn); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

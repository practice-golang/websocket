package main // import "ws"

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	_ "embed"

	ws "golang.org/x/net/websocket"
)

//go:embed html/index.html
var html string

var (
	Address      string = "localhost"
	Port         string = "1323"
	HttpProtocol string = "http"
)

func wsController(w http.ResponseWriter, req *http.Request) {
	ws.Handler(func(conn *ws.Conn) {
		defer conn.Close()

		handler := ws.Message

		fin := false
		for !fin {
			// Read
			message := ""
			err := handler.Receive(conn, &message)
			if err != nil {
				fmt.Println("Receive error:", err)
				fin = true
			}

			// Print received message
			if len(message) > 0 {
				fmt.Printf("%s\n", message)
			}

			// Write
			response := "Hello, Client! You sent me: " + message
			err = handler.Send(conn, response)
			if err != nil {
				fmt.Println("Send error:", err)
				fin = true
			}
		}
	}).ServeHTTP(w, req)
}

func main() {
	listenURI := Address + ":" + Port
	uri := HttpProtocol + "://" + Address + ":" + Port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})
	http.HandleFunc("/ws", wsController)

	fmt.Printf("Server is running at %s\n", uri)

	switch os := runtime.GOOS; os {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "linux":
		exec.Command("xdg-open", uri).Start()
	case "darwin":
		exec.Command("open", uri).Start()
	default:
		fmt.Printf("%s: unsupported platform", os)
	}

	http.ListenAndServe(listenURI, nil)
}

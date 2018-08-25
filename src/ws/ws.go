package main

import (
	"fmt"
	"net/http"

	wsgor "github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	wsbase "golang.org/x/net/websocket"
)

var (
	upgrader = wsgor.Upgrader{}
)

func wsBasic(c echo.Context) error {
	wsbase.Handler(func(ws *wsbase.Conn) {
		defer ws.Close()

		fin := false
		for fin == false {
			// Write
			err := wsbase.Message.Send(ws, "Hello, Client! I'm net/websocket")
			if err != nil {
				c.Logger().Error(err)
				fin = true
			}

			// Read
			msg := ""
			err = wsbase.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				fin = true
			}

			// Print or Do Something
			if len(msg) > 0 {
				fmt.Printf("%s\n", msg)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func wsGorilla(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	fin := false
	for fin == false {
		// Write
		err := ws.WriteMessage(wsgor.TextMessage, []byte("Hello, Client! I'm gorilla/websocket"))
		if err != nil {
			c.Logger().Error(err)
			fin = true
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			fin = true
		}

		// Print or Do Something
		if len(msg) > 0 {
			fmt.Printf("%s\n", msg)
		}
	}

	return nil
}

func main() {
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Static("/", "./index.html")

	e.GET("/ws", wsBasic)
	e.GET("/ws-gorilla", wsGorilla)

	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var connNum int = 0
var connPool [10]*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	var err error

	connNum++
	connPool[connNum], err = upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, message, err := connPool[connNum].ReadMessage()
			if err != nil {
				panic(err)
			}

			messageText := string(message)
			messageText = fmt.Sprintf("#%d: %s", connNum, messageText)

			for i := 0; i <= connNum; i++ {
				connPool[connNum].WriteMessage(websocket.TextMessage, []byte(messageText))
			}

			fmt.Println(messageText)
		}
	}()
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	http.ListenAndServe(":3000", nil)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var connNum int = 0
var connPool [1000]*websocket.Conn

type Message struct {
	Num  int
	Type string
	Text string
}

const MESSAGE_ACK = "ACK"
const MESSAGE_NORMAL = "NORMAL"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	msg := make(chan bool)

	var err error
	var conn *websocket.Conn

	conn, err = upgrader.Upgrade(w, r, nil)

	connNum++
	var inx int
	inx = connNum
	connPool[inx] = conn

	if err != nil {
		panic(err)
	}

	connPool[inx].WriteJSON(Message{
		Num:  inx,
		Type: MESSAGE_ACK,
		Text: "",
	})

	go func() {
		for {
			_, message, err := connPool[inx].ReadMessage()
			if err != nil {
				panic(err)
			}

			messageText := string(message)
			messageText = fmt.Sprintf("#%d: %s", inx, messageText)

			for i := 1; i <= connNum; i++ {
				connPool[i].WriteJSON(Message{
					Num:  inx,
					Type: MESSAGE_NORMAL,
					Text: messageText,
				})
			}

			fmt.Println(messageText)
		}

		msg <- true
	}()

	<-msg
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	http.ListenAndServe(":3000", nil)
}

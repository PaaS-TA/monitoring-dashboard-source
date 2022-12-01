package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Client struct {
	conn      *websocket.Conn
	broadcast chan []byte
}

type config struct {
	division string
}

var wg sync.WaitGroup

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", socketHandler)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) readFunc() {
	for {
		_, message, err := c.conn.ReadMessage() /*클라이언트에서 받은 websocket메세지 ReadMessage메소드로 읽음.*/

		fmt.Println("read" + string(message))

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			break
		}
		c.broadcast <- message /*broadcast채널로 대입*/

	}
	wg.Done()
}

func (c *Client) writeFunc() {

	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.broadcast: /*broadcast채널에서 데이터 꺼내기*/
			fmt.Println("write" + string(message))

			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.broadcast)
			for i := 0; i < n; i++ {
				w.Write(<-c.broadcast)
			}

			if err := w.Close(); err != nil {
				return
			}

		}
	}
	wg.Done()
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	wg.Add(2)

	conn, err := upgrader.Upgrade(w, r, nil) /*Websocket프로토콜로 전환*/
	defer conn.Close()

	client := &Client{conn: conn, broadcast: make(chan []byte, 256)}

	if err != nil {
		log.Printf("upgrader.Upgrade: %v", err)
		return
	}

	go client.writeFunc()
	go client.readFunc()

	wg.Wait()

}

// read함수 , write함수 분리
// 분리 후 goroutine 으로 호출
// 채널 맞추고 송수신 확인

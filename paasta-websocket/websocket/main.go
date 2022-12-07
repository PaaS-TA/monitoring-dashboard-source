package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type JsonData struct {
	Type  string `json:"type"`
	Time  string `json:"time"`
	Usage string `json:"usage"`
	Index int    `json:"Index"`
}

type JsonDataArray struct {
	array []JsonData
}

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
	http.Handle("/", http.FileServer(http.Dir("static/public")))
	http.HandleFunc("/ws", socketHandler)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) readFunc() {

	JsonData1 := JsonData{"cpu", "1669892640", "2", 0}
	JsonData2 := JsonData{"cpu", "1669892700", "3", 1}
	JsonData3 := JsonData{"cpu", "1669892760", "4", 2}
	JsonData4 := JsonData{"cpu", "1669892820", "5", 3}
	JsonData5 := JsonData{"cpu", "1669892880", "3", 4}
	JsonData6 := JsonData{"cpu", "1669892940", "10", 5}
	JsonData7 := JsonData{"cpu", "1669893000", "2", 6}
	JsonData8 := JsonData{"cpu", "1669893060", "5", 7}
	JsonData9 := JsonData{"cpu", "1669893120", "1", 8}
	JsonData10 := JsonData{"cpu", "1669893180", "4", 9}

	JsonDataArray := []JsonData{
		JsonData1, JsonData2, JsonData3,
		JsonData4, JsonData5, JsonData6,
		JsonData7, JsonData8, JsonData9, JsonData10}

	for i := 0; i <= len(JsonDataArray); i++ {
		b, _ := json.Marshal(JsonDataArray)
		c.broadcast <- b
	}

	/*for _, data := range JsonDataArray {
		b, _ := json.Marshal(data)
		fmt.Println("read: ", string(b))

		c.broadcast <- b
	}*/

	/*for {
		_, message, err := c.conn.ReadMessage() 클라이언트에서 받은 websocket메세지 ReadMessage메소드로 읽음.

		fmt.Println("read" + string(message))

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			break
		}
		c.broadcast <- message broadcast채널로 대입
	}*/
	wg.Done()
}

func (c *Client) writeFunc() {

	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message := <-c.broadcast: /*broadcast채널에서 데이터 꺼내기*/

			ticker := time.NewTicker(time.Millisecond * 3000)

			time.Sleep(time.Millisecond * 1600)
			ticker.Stop()

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			fmt.Println("write: " + string(message))

			/*n := len(c.broadcast)
			for i := 0; i < n; i++ {
				w.Write(<-c.broadcast)
			}*/

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

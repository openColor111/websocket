package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有CORS请求
		return true
	},
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	//var wg sync.WaitGroup
	//wg.Add(2)
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		log.Println("websocket upgrade failed:", err)
		return
	}
	defer conn.Close()
	for {
		//创建阻塞接收信息
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error:", err)
			break
		}
		log.Printf("received: %s", message)
		//创建不间断发送信息
		go func() {
			for {
				t := time.Now().Format("2006-01-02 15:04:05")
				// 这里可以处理接收到的消息，或者直接返回
				message = []byte("hello client " + t)
				err = conn.WriteMessage(mt, message)
				if err != nil {
					log.Println("error:", err)
					break
				}
				time.Sleep(time.Second * 5)
			}
		}()
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocketConnection)
	log.Println("Starting WebSocket server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 本番では制限する
	},
}

type FrameMessage struct {
	FrameID int    `json:"frameId"`
	Data    string `json:"data"`
}

type AckMessage struct {
	Type    string `json:"type"`
	FrameID int    `json:"frameId"`
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connected!")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var frame FrameMessage
		if err := json.Unmarshal(msg, &frame); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		log.Printf("Received frameId: %d\n", frame.FrameID)

		ack := AckMessage{
			Type:    "ack",
			FrameID: frame.FrameID,
		}
		ackBytes, _ := json.Marshal(ack)
		conn.WriteMessage(websocket.TextMessage, ackBytes)
	}
}

func main() {
	r := gin.Default()
	r.GET("/ws", handleWebSocket)

	log.Println("Starting server at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed:", err)
	}
}

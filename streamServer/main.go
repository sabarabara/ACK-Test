/*
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
*/
/*
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // CORS許可（開発用）
	},
}

func main() {
	router := gin.Default()

	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		for {
			// メッセージを受け取る
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("📩 Received: %s", message)

			// JSONの中から frameId を抽出（超簡易パーサー）
			// {"frameId":1,"data":"xxx"} という前提
			frameId := extractFrameId(string(message))

			// 偶数のframeIdはACKを送らない
			if frameId%2 == 0 {
				log.Printf("❌ Skipping ACK for frameId %d", frameId)
				continue
			}

			// ACK送信
			ackMsg := fmt.Sprintf(`{"type":"ack","frameId":%d}`, frameId)
			err = conn.WriteMessage(mt, []byte(ackMsg))
			if err != nil {
				log.Println("Write error:", err)
				break
			}
			log.Printf("✅ Sent ACK for frameId %d", frameId)
		}
	})

	router.Run(":8080")
}

// 文字列から frameId を抽出（簡易実装）
func extractFrameId(msg string) int {
	start := len(`{"frameId":`)
	end := len(msg)
	for i := start; i < len(msg); i++ {
		if msg[i] == ',' {
			end = i
			break
		}
	}
	idStr := msg[start:end]
	id, _ := strconv.Atoi(idStr)
	return id
}
*/

package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

// クライアントから受信するフレーム構造体（camelCaseに注意！）
type Frame struct {
    ClientId int    `json:"clientId"`
    FrameId  int    `json:"frameId"`
    Data     string `json:"data"`
}

// クライアントに返すACKメッセージ
type Ack struct {
    FrameId int `json:"frameId"`
}

// WebSocketアップグレーダー
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func main() {
    r := gin.Default()

    r.GET("/ws", func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Println("Failed to set websocket upgrade:", err)
            return
        }
        defer conn.Close()

        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                log.Println("Read error:", err)
                break
            }

            var frame Frame
            if err := json.Unmarshal(msg, &frame); err != nil {
                log.Println("Unmarshal error:", err)
                continue
            }

            log.Printf("📩 Received: %+v\n", frame)

            // ACKをスキップするかどうか（frameId % 5 == 0 でスキップ）
            if frame.FrameId%5 == 0 {
                log.Printf("❌ Skipping ACK for frameId %d\n", frame.FrameId)
                continue
            }

            ack := Ack{FrameId: frame.FrameId}
            ackJson, _ := json.Marshal(ack)
            if err := conn.WriteMessage(websocket.TextMessage, ackJson); err != nil {
                log.Println("Write ACK error:", err)
                break
            }
            log.Printf("✅ Sent ACK for frameId %d\n", frame.FrameId)
        }
    })

    log.Println("✅ Server started on :8080")
    r.Run(":8080")
}

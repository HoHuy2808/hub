package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	UserID string // ID của người dùng
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Thêm CheckOrigin để bỏ qua lỗi CORS khi test từ localhost:5500 sang localhost:2808
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// cài đặt giới hạn
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("read error:", err)
			}
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub đóng channel
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// Định kỳ gửi Ping để giữ kết nối không bị ngắt
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs xử lý request nâng cấp lên WebSocket
func ServeWs(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy userID từ URL query parameter (ví dụ: ?userID=1)
		// Trong thực tế, bạn sẽ lấy từ JWT token (c.MustGet("userID"))
		userID := c.Query("userID")
		if userID == "" {
			// Bắt buộc phải có userID để kết nối trong mô hình mạng xã hội
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userID in query parameter"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Lỗi nâng cấp WebSocket:", err)
			return
		}

		// Tạo một Client mới
		client := &Client{
			UserID: userID,
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
		}

		// Đăng ký Client này với Hub
		client.hub.register <- client

		// Bật 2 Goroutine chạy song song để ngóng tin nhắn (Đọc) và chờ lệnh (Ghi)
		go client.WritePump()
		go client.ReadPump()
	}
}

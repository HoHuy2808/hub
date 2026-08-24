package websocket

// NotificationMessage định nghĩa cấu trúc tin nhắn để gửi cho 1 user cụ thể
type NotificationMessage struct {
	TargetUserID string
	Data         []byte
}

type Hub struct {
	// Danh sách client theo dạng map[UserID]map[*Client]bool
	// Cho phép 1 user kết nối từ nhiều thiết bị (web, điện thoại)
	clients map[string]map[*Client]bool

	// Kênh gửi tin nhắn đích danh
	SendToUser chan *NotificationMessage

	// Kênh broadcast nếu muốn gửi thông báo hệ thống cho toàn bộ user
	broadcast chan []byte

	register   chan *Client
	unregister chan *Client
}

// Hàm khởi tạo một Hub mới
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		SendToUser: make(chan *NotificationMessage),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Nếu user này chưa có danh sách thiết bị nào, tạo mới
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true

		case client := <-h.unregister:
			if connections, ok := h.clients[client.UserID]; ok {
				if _, ok := connections[client]; ok {
					delete(connections, client)
					close(client.send)

					// Nếu user không còn mở thiết bị nào nữa thì xóa luôn key UserID để nhẹ RAM
					if len(connections) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}

		case msg := <-h.SendToUser:
			// Lấy danh sách thiết bị đang online của User này
			if connections, ok := h.clients[msg.TargetUserID]; ok {
				for client := range connections {
					select {
					case client.send <- msg.Data:
					default:
						close(client.send)
						delete(connections, client)
						if len(connections) == 0 {
							delete(h.clients, msg.TargetUserID)
						}
					}
				}
			}

		case message := <-h.broadcast:
			// Gửi thông báo hệ thống cho toàn bộ người dùng đang online
			for _, connections := range h.clients {
				for client := range connections {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(connections, client)
					}
				}
			}
		}
	}
}

package messages

type Message struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

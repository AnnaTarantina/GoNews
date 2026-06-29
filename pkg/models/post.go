package models

import "time"

// Post представляет публикацию из RSS-ленты
type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	PubTime time.Time `json:"pub_time"`
	Link    string    `json:"link"`
	Source  string    `json:"source"`
}

// Pagination представляет данные для постраничной навигации
type Pagination struct {
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
}

// NewsResponse представляет ответ API с новостями и пагинацией
type NewsResponse struct {
	Posts      []Post     `json:"posts"`
	Pagination Pagination `json:"pagination"`
}

package account

import "time"

type Account struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Title      string    `json:"title"`
	IsDefault  bool      `json:"is_default"`
	CreateTime time.Time `json:"create_time"`
}

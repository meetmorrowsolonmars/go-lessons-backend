package user

import "time"

type User struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	CreateTime time.Time `json:"create_time"`
}

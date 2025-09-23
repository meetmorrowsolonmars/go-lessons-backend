package model

import "time"

type Account struct {
	ID         int64
	UserID     int64
	Title      string
	IsDefault  bool
	CreateTime time.Time
}

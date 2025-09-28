package operation

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Operation struct {
	ID          uuid.UUID       `json:"id"`
	UserID      int64           `json:"user_id"`
	AccountID   int64           `json:"account_id"`
	Type        string          `json:"type"`
	CategoryID  int64           `json:"category_id,omitempty"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	CreateTime  time.Time       `json:"create_time"`
}

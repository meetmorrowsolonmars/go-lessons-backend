package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Operation struct {
	ID          uuid.UUID
	UserID      int64
	AccountID   int64
	Type        OperationType
	CategoryID  int64
	Amount      decimal.Decimal
	Description string
	CreateTime  time.Time
}

type OperationType string

const (
	OperationTypeDebit  OperationType = "debit"
	OperationTypeCredit OperationType = "credit"
)

package models

import (
    "github.com/shopspring/decimal"
    "github.com/google/uuid"
)

type Transaction struct {
    BaseModel
    TransactionID       string          `gorm:"uniqueIndex;size:255;not null" json:"transaction_id"`
    UserID              uuid.UUID       `gorm:"type:uuid;not null;index" json:"-"`
    SubscriptionModelID uuid.UUID       `gorm:"type:uuid;not null" json:"-"`
    Amount              decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount"`
    Currency            string          `gorm:"size:3;not null;default:'INR'" json:"currency"`
    PaymentGateway      string          `gorm:"size:255;not null" json:"payment_gateway"`
    PaymentGatewayTxnID string          `gorm:"size:255;not null" json:"payment_gateway_txn_id"`
    Status              string          `gorm:"size:20;not null;default:'pending';comment:pending,success,failed" json:"status"`
    PaymentMethod       string          `gorm:"size:50" json:"payment_method"`
    Description         string          `gorm:"size:500" json:"description"`
    
    // Relationships
    User              User              `gorm:"foreignKey:UserID;references:ID"`
    SubscriptionModel SubscriptionModel `gorm:"foreignKey:SubscriptionModelID;references:ID"`
}

func (Transaction) TableName() string {
    return "transactions"
}


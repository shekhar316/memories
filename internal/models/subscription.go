package models

import (
    "time"
    "github.com/shopspring/decimal"
    "github.com/google/uuid"
)

type SubscriptionModel struct {
    BaseModel
    Name             string          `gorm:"size:255;not null" json:"name"`
    Description      string          `gorm:"size:500" json:"description"`
    Price            decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
    QuotaGB          int             `gorm:"not null" json:"quota_gb"`
    ValidityMonths   int             `gorm:"not null;comment:subscription validity in months" json:"validity_months"`
    IsAIEnabled      bool            `gorm:"default:false" json:"is_ai_enabled"`
    
    // Relationships
    UserSubscriptions []UserSubscription `gorm:"foreignKey:SubscriptionModelID;references:ID"`
    Transactions      []Transaction      `gorm:"foreignKey:SubscriptionModelID;references:ID"`
}

func (SubscriptionModel) TableName() string {
    return "subscription_models"
}

type UserSubscription struct {
    BaseModel
    UserID              uuid.UUID        `gorm:"type:uuid;not null;index" json:"-"`
    SubscriptionModelID uuid.UUID        `gorm:"type:uuid;not null" json:"-"`
    StartDate           time.Time        `gorm:"not null" json:"start_date"`
    EndDate             time.Time        `gorm:"not null" json:"end_date"`
    UsedStorageGB       decimal.Decimal  `gorm:"type:decimal(10,3);default:0" json:"used_storage_gb"`
    
    // Relationships
    User              User              `gorm:"foreignKey:UserID;references:ID"`
    SubscriptionModel SubscriptionModel `gorm:"foreignKey:SubscriptionModelID;references:ID"`
}

func (UserSubscription) TableName() string {
    return "user_subscriptions"
}


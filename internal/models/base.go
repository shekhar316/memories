package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
    IsActive  int            `gorm:"default:1;comment:1=active,0=inactive" json:"is_active"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID if not set
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
    if b.ID == uuid.Nil {
        b.ID = uuid.New()
    }
    return nil
}

func (b *BaseModel) IsUserActive() bool {
    return b.IsActive == 1
}

func (b *BaseModel) ActivateUser() {
    b.IsActive = 1
}

func (b *BaseModel) DeactivateUser() {
    b.IsActive = 0
}


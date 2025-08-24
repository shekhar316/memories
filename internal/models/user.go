package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    BaseModel
    FirstName       string `gorm:"size:25;not null" json:"first_name"`
    LastName        string `gorm:"size:25" json:"last_name"`
    Email           string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	CountryCode     string `gorm:"size:5" json:"country_code"`
    PhoneNumber     string `gorm:"size:10" json:"phone_number"`
    Username        string `gorm:"uniqueIndex;size:25;not null" json:"username"`
    Password        string `gorm:"size:255" json:"-"`
    GoogleID        string `gorm:"size:255;uniqueIndex" json:"-"`
    LastLoginAt     *time.Time `json:"last_login_at"`
    EmailVerifiedAt *time.Time `json:"email_verified_at"`
    
    // Relationships

	UserSubscription *UserSubscription `gorm:"foreignKey:UserID;references:ID"`
    Albums          []Album            `gorm:"foreignKey:OwnerID;references:ID"`
    Photos          []Photo            `gorm:"foreignKey:OwnerID;references:ID"`
    OwnedShares     []AlbumShare       `gorm:"foreignKey:SharedBy;references:ID"`
    SharedAlbums    []AlbumShare       `gorm:"foreignKey:UserID;references:ID"`
    Transactions    []Transaction      `gorm:"foreignKey:UserID;references:ID"`
}

// TableName overrides the table name
func (User) TableName() string {
    return "users"
}

// Methods
func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}


type ProfilePhoto struct {
    BaseModel
    UserID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
    PhotoID uuid.UUID `gorm:"type:uuid;not null" json:"photo_id"`
    
    // Relationships
    User  User  `gorm:"foreignKey:UserID;references:ID"`
    Photo Photo `gorm:"foreignKey:PhotoID;references:ID"`
}

func (ProfilePhoto) TableName() string {
    return "profile_photos"
}
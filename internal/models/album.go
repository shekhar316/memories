package models

import (
    "github.com/google/uuid"
)

type Album struct {
    BaseModel
    Name            string     `gorm:"size:255;not null" json:"name"`
    Description     string     `gorm:"size:1000" json:"description"`
    OwnerID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"owner_id"`
    CoverPhotoID    *uuid.UUID `gorm:"type:uuid" json:"cover_photo_id"`
    IsPublic        bool       `gorm:"default:false" json:"is_public"`
    IsShared        bool       `gorm:"default:false" json:"is_shared"`
    PhotoCount      int        `gorm:"default:0" json:"photo_count"`
    TotalSizeBytes  int64      `gorm:"default:0" json:"total_size_bytes"`
    
    // Relationships
    Owner       User         `gorm:"foreignKey:OwnerID;references:ID"`
    CoverPhoto  *Photo       `gorm:"foreignKey:CoverPhotoID;references:ID"`
    Photos      []Photo      `gorm:"many2many:album_photos;"`
    SharedUsers []AlbumShare `gorm:"foreignKey:AlbumID;references:ID"`
}

func (Album) TableName() string {
    return "albums"
}


type AlbumShare struct {
    BaseModel
    AlbumID  uuid.UUID `gorm:"type:uuid;not null;index" json:"album_id"`
    UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    Role   int `gorm:"default:1;comment:5=owner,2=collab,1=view;not null" json:"role"`
    SharedBy uuid.UUID `gorm:"type:uuid;not null" json:"shared_by"`
    
    // Relationships
    Album        Album     `gorm:"foreignKey:AlbumID;references:ID"`
    User         User      `gorm:"foreignKey:UserID;references:ID"`
    SharedByUser User      `gorm:"foreignKey:SharedBy;references:ID"`
}

func (AlbumShare) TableName() string {
    return "album_shares"
}


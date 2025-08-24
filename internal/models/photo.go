package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
    "github.com/pgvector/pgvector-go"
)

type Photo struct {
    BaseModel
    OwnerID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"-"`
    OriginalName    string     `gorm:"size:255;not null" json:"original_name"`
    StoragePath     string     `gorm:"size:500;not null" json:"storage_path"`
    ThumbnailPath   string     `gorm:"size:500" json:"thumbnail_path"`
    FileSize        int64      `gorm:"not null" json:"file_size"`
    MimeType        string     `gorm:"size:255" json:"mime_type"`
    Hash            string     `gorm:"size:64;index" json:"hash"`
    
    // EXIF Data
    DateTaken       *time.Time `json:"date_taken"`
    LocationLat     *float64   `gorm:"type:decimal(10,8)" json:"location_lat"`
    LocationLng     *float64   `gorm:"type:decimal(11,8)" json:"location_lng"`
    
    // Data
    Tags            pq.StringArray `gorm:"type:text[];column:tags" json:"tags"`
    AiTags          pq.StringArray `gorm:"type:text[];column:tags" json:"ai_tags"`
    
    ImageEmbedding  pgvector.Vector `gorm:"type:vector(512)" json:"image_embedding,omitempty"`

    // Upload Info
    UploadedAt      time.Time `gorm:"not null" json:"uploaded_at"`
    
    // Relationships
    Owner   User    `gorm:"foreignKey:OwnerID;references:ID"`
    Albums  []Album `gorm:"many2many:album_photos;"`
}

func (Photo) TableName() string {
    return "photos"
}

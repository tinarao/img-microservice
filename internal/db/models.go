package db

import (
	"time"
)

type User struct {
	Id            uint      `gorm:"primarykey" json:"id"`
	Name          string    `json:"name"`
	Provider      string    `json:"provider"`
	Email         string    `gorm:"uniqueIndex" json:"email"`
	NickName      *string   `json:"nickName"`
	AvatarURL     *string   `json:"avatarUrl"`
	Locale        string    `gorm:"default:ru" json:"locale"`
	IsActive      bool      `gorm:"default:true" json:"isActive"`
	LastLogin     time.Time `json:"lastLogin"`
	PublicApiKey  string    `json:"publicApiKey"`
	PrivateApiKey string    `json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type OperationType string

const (
	OperationTypeResize   OperationType = "resize"
	OperationTypeCrop     OperationType = "crop"
	OperationTypeRotate   OperationType = "rotate"
	OperationTypeFilter   OperationType = "filter"
	OperationTypeCompress OperationType = "compress"
)

type OperationStatus string

const (
	StatusPending    OperationStatus = "pending"
	StatusProcessing OperationStatus = "processing"
	StatusCompleted  OperationStatus = "completed"
	StatusFailed     OperationStatus = "failed"
)

type Operation struct {
	ID         uint            `gorm:"primarykey" json:"id"`
	UserID     uint            `json:"userId"`
	User       User            `gorm:"foreignKey:UserID" json:"user"`
	Type       OperationType   `json:"type"`
	GoalWidth  *uint           `json:"goalWidth"`
	InputPath  string          `json:"inputPath"`
	OutputPath string          `json:"outputPath"`
	Status     OperationStatus `gorm:"default:pending" json:"status"` // pending, processing, completed, failed
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

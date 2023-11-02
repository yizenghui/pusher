package db

import (
	"time"
)

// User 用户
type User struct {
	ID      uint   `gorm:"primarykey"`
	Openid  string `gorm:"index"`
	Created int
	Posted  int
	Related int
	Fromid  int
	Isblock int
}

// File 文件
type File struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint
	ToUserID   uint
	FileName   string
	FileSize   int64
	MediaType  string
	URL        string
	Status     int32
	MediaCate  int32
	Hash       string
	CompleteAt time.Time
	CheckAt    time.Time
}

// FileTrace 文件
type FileTrace struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	FileID    uint
	TraceID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

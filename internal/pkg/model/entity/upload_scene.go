package entity

import (
	"mio/internal/pkg/model"
	"time"
)

type UploadScene struct {
	ID        int               `gorm:"primaryKey;not null;type:serial;comment:上传文件记录表"`
	MaxCount  int               `gorm:"type:int4;not null;comment:每天最多上传多少次文件"`
	MaxSize   int64             `gorm:"type:int8;not null;comment:上传文件大小限制 单位B 1MB=1024KB=1048576B"`
	MustLogin bool              `gorm:"type:bool;not null;comment:用户是否必须登录"`
	OssDir    string            `gorm:"type:varchar(100);not null;comment:对象存储路径"`
	MimeTypes model.ArrayString `gorm:"type:varchar(500);not null;comment:可上传的文件mime类型多个用英文逗号隔开 image/png,image/jpg"`
	Scene     string            `gorm:"type:varchar(50);not null;comment:上传场景标识 必须是英文字母 例如 userAvatar"`
	SceneName string            `gorm:"type:varchar(50);not null;comment:上传场景标识名称 例如 用户头像"`
	CreatedAt time.Time         `gorm:"type:timestamp;not null"`
	UpdatedAt time.Time         `gorm:"type:timestamp;not null"`
}

func (UploadScene) TableName() string {
	return "upload_scene"
}

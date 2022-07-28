package entity

import "time"

type UploadLog struct {
	ID           int64     `gorm:"primaryKey;not null;type:serial8;comment:上传文件表"`
	LogId        string    `gorm:"type:varchar(100);not null;comment:文件编号"`
	OssPath      string    `gorm:"type:varchar(1000);not null;comment:阿里云oss路径"`
	Size         int64     `gorm:"type:int8;not null;comment:文件大小 单位B"`
	Url          string    `gorm:"type:varchar(1000);not null;default:'';comment:文件链接"`
	OperatorId   int64     `gorm:"type:int8;not null;default:0;comment:操作员"`
	OperatorType int8      `gorm:"type:int2;not null;default:1;comment:操作员类型 1用户 2管理员"`
	SceneId      int       `gorm:"type:int4;not null;comment:上传场景编号"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null"`
	UpdatedAt    time.Time `gorm:"type:timestamp;not null"`
}

func (UploadLog) TableName() string {
	return "upload_log"
}

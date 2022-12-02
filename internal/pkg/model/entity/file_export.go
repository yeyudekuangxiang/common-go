package entity

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/enum"
	"mio/internal/pkg/model"
)

type FileExportType interface {
	enum.IEnumIntStatus
}

var (
	FileExportTypePoint    FileExportType = enum.NewEnumIntStatus(1, "积分明细", "积分明细")
	FileExportTypeCoupon   FileExportType = enum.NewEnumIntStatus(2, "优惠券码", "优惠券码")
	FileExportTypeExchange FileExportType = enum.NewEnumIntStatus(3, "兑换券码", "兑换券码")
)

var FileExportTypeList = enum.EnumIntStatusList{
	FileExportTypePoint,
	FileExportTypeCoupon,
	FileExportTypeExchange,
}

type FileExportStatus int

func (f FileExportStatus) Text() string {
	switch f {
	case FileExportStatusWaiting:
		return "未开始"
	case FileExportStatusProgress:
		return "进行中"
	case FileExportStatusSuccess:
		return "导出成功"
	case FileExportStatusFailed:
		return "导出失败"
	}
	return "未知状态"
}

const (
	FileExportStatusWaiting  FileExportStatus = 1 //未开始
	FileExportStatusProgress FileExportStatus = 2 //进行中
	FileExportStatusSuccess  FileExportStatus = 3 //导出成功
	FileExportStatusFailed   FileExportStatus = 4 //导出失败
)

var FileExportStatusList = []FileExportStatus{
	FileExportStatusWaiting,
	FileExportStatusProgress,
	FileExportStatusSuccess,
	FileExportStatusFailed,
}

const (
	OrderByFileExportCreatedAtDesc OrderBy = "order_by_file_export_created_at_desc"
	OrderByFileExportUpdatedAtDesc OrderBy = "order_by_file_export_updated_at_desc"
)

type FileExport struct {
	ID        int64            `json:"id"`
	AdminId   int              `json:"adminId"`
	Url       string           `json:"url"`
	Params    string           `json:"params"`
	Status    FileExportStatus `json:"status"` //1 未开始 2进行中 3导出成功 4导出失败
	Message   string           `json:"message"`
	Type      int64            `json:"type"` // 1 积分明细
	CreatedAt model.Time       `json:"createdAt"`
	UpdatedAt model.Time       `json:"updatedAt"`
}

func (FileExport) TableName() string {
	return "file_export"
}

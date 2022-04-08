package entity

const (
	OrderByTagSortDesc OrderBy = "order_by_tag_sort_desc"
	OrderByTagSortAsc  OrderBy = "order_by_tag_sort_asc"
)

// 标签
type Tag struct {
	Id           int64  `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	Name         string `gorm:"size:32;unique" json:"name" form:"name"`            // 名称
	Description  string `gorm:"size:1024" json:"description" form:"description"`   // 描述
	Img          string `gorm:"size:1024" json:"img" form:"img"`                   // 图标
	Sort         int    `gorm:"index:idx_sort_" json:"sort" form:"sort"`           // 排序编号
	Icon         string `gorm:"size:1024" json:"icon" form:"icon"`                 // icon
	ImgWithCover string `gorm:"size:1024" json:"imgWithCover" form:"imgWithCover"` // 图片加遮罩
}

func (Tag) TableName() string {
	return "tag"
}

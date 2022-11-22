package pagetool

import "math"

func TotalPage(total int64, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func GetPageInfo(page int, pageSize int, total int64) PageInfo {
	totalPage := TotalPage(total, pageSize)
	return PageInfo{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
		HasNext:   page < totalPage,
	}
}

type PageInfo struct {
	Total     int64
	Page      int
	PageSize  int
	TotalPage int
	HasNext   bool
}

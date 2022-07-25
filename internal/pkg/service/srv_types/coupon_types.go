package srv_types

import "mio/internal/pkg/util/timeutils"

type GetPageCouponRecordDTO struct {
	OpenId string
	Offset int
	Limit  int
}
type BaseCouponRecordDTO struct {
	ID         int64          `json:"id"`
	CoverImage string         `json:"coverImage"`
	Title      string         `json:"title"`
	UpdateDate timeutils.Date `json:"updateDate"`
}

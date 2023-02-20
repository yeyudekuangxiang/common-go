package srv_types

import "gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"

type GetPageCouponRecordDTO struct {
	OpenId string
	Offset int
	Limit  int
}
type BaseCouponRecordDTO struct {
	ID         int64         `json:"id"`
	CoverImage string        `json:"coverImage"`
	Title      string        `json:"title"`
	UpdateDate timetool.Date `json:"updateDate"`
}

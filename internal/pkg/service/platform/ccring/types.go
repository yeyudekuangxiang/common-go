package ccring

type ccRingReqParams struct {
	MemberId            string  `json:"memberId"`
	DegreeOfCharge      float64 `json:"degreeOfCharge"`
	ProductCategoryName string  `json:"productCategoryName"`
	Name                string  `json:"name"`
	Qua                 string  `json:"qua"`
}

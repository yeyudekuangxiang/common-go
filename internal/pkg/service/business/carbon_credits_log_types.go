package business

import (
	"mio/internal/pkg/model/entity/business"
	business2 "mio/internal/pkg/repository/business"
)

type CarbonCreditsLogSortedListResponse struct {
	Id    int
	Total string
	Title string
	Icon  string
	Type  business.CarbonType
}

type CarbonCreditLogListHistoryResponse struct {
	Total  string
	Month  string
	Detail []business2.CarbonCreditsLogListHistory
	Title  string
}

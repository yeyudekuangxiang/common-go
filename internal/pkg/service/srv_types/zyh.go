package srv_types

import (
	"mio/internal/pkg/model/entity"
)

type GetZyhGetInfoByDTO struct {
	Openid string
	VolId  string
}

type GetZyhLogAddDTO struct {
	Openid         string                      `json:"openid"`
	PointType      entity.PointTransactionType `json:"point_type"`
	Value          int64                       `json:"value"`
	ResultCode     string                      `json:"result_code"`
	AdditionalInfo string                      `json:"additional_info"`
	TransactionId  string                      `json:"transaction_id"`
}

type GetZyhOpenDTO struct {
	VolId  string
	Mobile string
}

type GetZyhLogDTO struct {
	PointType  string `json:"point_type"`
	PointValue int64  `json:"point_value"`
	ResultCode string `json:"result_code"`
	CreateTime string `json:"create_time"`
}

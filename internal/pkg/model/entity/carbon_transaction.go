package entity

import (
	"mio/internal/pkg/model"
)

type CarbonTransactionType string

// Text 获取积分类型的中文名称(给用户看的)
func (p CarbonTransactionType) Text() string {
	switch p {
	case CARBON_STEP:
		return "步行"
	}
	return "未知积分"
}

// RealText 获取积分类型的中文名称(给管理员看的)
func (p CarbonTransactionType) RealText() string {
	switch p {
	case CARBON_STEP:
		return "兑吧话费"
	case CARBON_COFFEE_CUP:
		return "兑吧流量"
	}
	return p.Text()
}

const (
	CARBON_STEP       CarbonTransactionType = "STEP"       //步行
	CARBON_COFFEE_CUP CarbonTransactionType = "COFFEE_CUP" //自带咖啡杯
)

var CarbonTransactionTypeList = []CarbonTransactionType{
	CARBON_STEP,
	CARBON_COFFEE_CUP,
}

var CarbonCollectValueMap = map[CarbonTransactionType]int{
	CARBON_STEP:       150,
	CARBON_COFFEE_CUP: 150,
}

var CarbonCollectLimitMap = map[CarbonTransactionType]int{
	CARBON_STEP:       4,
	CARBON_COFFEE_CUP: 2,
}

/*
CREATE TABLE "public"."carbon_transaction" (
  "id" int8 NOT NULL,
  "openid" text COLLATE "pg_catalog"."default" NOT NULL,
  "user_id" int8 NOT NULL DEFAULT 0,
  "transaction_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "type" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "city" int8  NOT NULL DEFAULT 0,
  "value" numeric(10,2) NOT NULL,
  "info" text COLLATE "pg_catalog"."default",
  "admin_id" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  CONSTRAINT "carbon_transaction_pk" PRIMARY KEY ("id"),
  CONSTRAINT "carbon_transaction_transaction_id_key" UNIQUE ("transaction_id"),
  CONSTRAINT "carbon_transaction_user_id_key" UNIQUE ("user_id")
);
*/

type CarbonTransaction struct {
	ID            int64                 `json:"id"`
	OpenId        string                `gorm:"column:openid" json:"openId"`
	UserId        int64                 `json:"userId"`
	TransactionId string                `json:"transactionId"`
	Type          CarbonTransactionType `json:"type"`
	City          int                   `json:"city"`
	Value         float64               `json:"value"`
	Info          string                `json:"info"`
	AdminId       int                   `json:"adminId"`
	CreatedAt     model.Time            `json:"createdAt"`
	UpdatedAt     model.Time            `json:"updatedAt"`
}

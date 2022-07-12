package duiba

import "time"

type VirtualGoodLog struct {
	ID            int64     `json:"id" gorm:"primaryKey;type:serial8;not null;comment:兑吧虚拟商品充值记录"`
	AppKey        string    `json:"appKey" gorm:"type:varchar(255);not null;uniqueIndex:duiba_virtual_good_unique,priority:3;comment:应用appkey"` //接口appKey，应用的唯一标识
	Timestamp     int64     `json:"timestamp" gorm:"type:int8;not null;comment:时间戳"`                                                            //1970-01-01开始的时间戳，毫秒为单位。
	Uid           string    `json:"uid" gorm:"type:varchar(255);not null;comment:用户唯一性标识"`                                                      //用户唯一性标识
	Sign          string    `json:"sign" gorm:"type:varchar(255);not null;comment:签名"`
	OrderNum      string    `json:"orderNum" gorm:"type:varchar(255);not null;uniqueIndex:duiba_virtual_good_unique,priority:1;comment:兑吧中奖订单号"`
	DevelopBizId  string    `json:"developBizId" gorm:"type:varchar(255);not null;default:'';comment:开发者订单号"`
	Params        string    `json:"params" gorm:"type:varchar(255);not null;uniqueIndex:duiba_virtual_good_unique,priority:2;comment:虚拟商品标识符"`
	Description   string    `json:"description" gorm:"type:varchar(255);not null;default:'';comment:文案描述"`
	Account       string    `json:"account" gorm:"type:varchar(255);not null;default:'';comment:用户兑换虚拟商品时输入的账号"`
	SupplierBizId string    `json:"supplierBizId" gorm:"type:varchar(255);not null;comment:订单流水号，开发者返回给兑吧的凭据"`
	CreatedAt     time.Time `json:"createdAt" gorm:"type:timestamptz;not null;comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"type:timestamptz;not null;comment:更新时间"`
}

func (VirtualGoodLog) TableName() string {
	return "duiba_virtual_good_log"
}

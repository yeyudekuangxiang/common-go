package activity

import "mio/model"

type BocApplyRecord struct {
	Id                 int64      `json:"id"`
	UserId             int64      `json:"userId"`
	ApplyStatus        int8       `json:"applyStatus"`        //卡片申请状态 1申请中 2申请成功 3申请失败
	IsGiveApplyBonus   int        `json:"IsGiveApplyBonus"`   //申请卡片奖励发放状态 1未发放 2已发放
	GiveApplyBonusTime model.Time `json:"giveApplyBonusTime"` //申请卡片奖励发放时间
	ShareUserId        int64      `json:"shareUserId"`
	IsBindWechat       int        `json:"isBindWechat"`      //卡片是否绑定微信 1未绑定 2已绑定
	IsGiveBindBonus    int        `json:"isGiveBindBonus"`   //是否发放绑定微信奖励 1未发放 2已发放
	GiveBindBonusTime  model.Time `json:"giveBindBonusTime"` //发放绑定微信奖励时间
	CreatedAt          model.Time `json:"createAt"`
	UpdatedAt          model.Time `json:"updatedAt"`
}

func (BocApplyRecord) TableName() string {
	return "boc_apply_record"
}

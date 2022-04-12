package activity

import (
	"mio/internal/pkg/model"
	"time"
)

type BocRecord struct {
	Id                      int64      `json:"id"`
	UserId                  int64      `json:"userId"`
	ShareUserId             int64      `json:"shareUserId"`             //邀请者id
	ApplyStatus             int8       `json:"applyStatus"`             //卡片申请状态 1参与中 2审核中 3申请成功 4申请失败
	ApplyBonusStatus        int        `json:"ApplyBonusStatus"`        //申请卡片奖励发放状态 1未申请 2审核中 3已发放
	ApplyBonusTime          model.Time `json:"applyBonusTime"`          //申请卡片奖励发放时间
	BindWechatStatus        int        `json:"bindWechatStatus"`        //卡片是否绑定微信 1未绑定 2已绑定
	BindWechatBonusStatus   int        `json:"bindWechatBonusStatus"`   //申请卡片奖励发放状态 1未申请 审核中 3已发放
	BindWechatBonusTime     model.Time `json:"bindWechatBonusTime"`     //发放绑定微信奖励时间
	AnswerStatus            int        `json:"answerStatus"`            //答题状态 1未作答 2回答正确 3回答错误
	AnswerBonusStatus       int        `json:"answerBonusStatus"`       //答题积分发放状态 1未发放 2已发放
	AnswerBonusTime         model.Time `json:"answerBonusTime"`         //答题积分发放时间
	ShareNum                int        `json:"shareNum"`                //我邀请到的人数
	ShareUserBocBonusStatus int        `json:"shareUserBocBonusStatus"` //未领取 2已领取
	ShareUserBocBonusTime   model.Time `json:"shareUserBocBonusTime"`
	Source                  string     `json:"source"` //用户来源 mio-dialog mio-banner mio-oa(297489) mio-poster(297549)  boc-sms(297108) boc-oa(297490) boc-app(297492)
	CreatedAt               model.Time `json:"createAt"`
	UpdatedAt               model.Time `json:"updatedAt"`
}

func (BocRecord) TableName() string {
	return "boc_record"
}
func NewBocRecord() BocRecord {
	return BocRecord{
		ApplyStatus:           1,
		ApplyBonusStatus:      1,
		BindWechatStatus:      1,
		BindWechatBonusStatus: 1,
		AnswerStatus:          1,
		AnswerBonusStatus:     1,
		CreatedAt:             model.Time{Time: time.Now()},
		UpdatedAt:             model.Time{Time: time.Now()},
	}
}
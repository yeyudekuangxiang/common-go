package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	duibaApi "mio/pkg/duiba/api/model"
	"time"
)

type TopicDetail struct {
	entity.Topic
	IsLike        bool             `json:"isLike"`
	UpdatedAtDate string           `json:"updatedAtDate"` //03-01
	User          entity.ShortUser `json:"user"`
}

type CreatePointTransactionParam struct {
	OpenId       string                      `binding:"required"`
	Type         entity.PointTransactionType `binding:"required"`
	Value        int64
	AdminId      int
	BizId        string
	Note         string
	AdditionInfo string
}
type CreateUserParam struct {
	OpenId      string            `json:"openId"`
	AvatarUrl   string            `json:"avatarUrl"`
	Gender      string            `json:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Nickname    string            `json:"nickname"`
	Birthday    model.Date        `json:"birthday"`
	PhoneNumber string            `json:"phoneNumber"`
	Source      entity.UserSource `json:"source" binding:"oneof=mio mobile"`
	UnionId     string            `json:"unionId"`
	ChannelId   int64             `json:"int64"`
}

type unidianTypeId struct {
	Test       string
	FiveYuan   string
	TwentyYuan string
}

var UnidianTypeId = unidianTypeId{
	Test:       "10013", // 测试
	FiveYuan:   "10689", // 5元话费
	TwentyYuan: "10725", // 20元话费
}

type SubmitOrderParam struct {
	Order           SubmitOrder
	Items           []SubmitOrderItem
	PartnershipType entity.PartnershipType
}
type SubmitOrder struct {
	AddressId string
	UserId    int64
	OrderType entity.OrderType
}
type SubmitOrderItem struct {
	ItemId string
	Count  int
}

type submitOrderParam struct {
	Order           submitOrder
	Items           []submitOrderItem
	PartnershipType entity.PartnershipType
}
type submitOrder struct {
	AddressId string
	UserId    int64
	TotalCost int
	OrderType entity.OrderType
}
type submitOrderItem struct {
	ItemId string
	Count  int
	Cost   int
}
type SubmitOrderForGreenParam struct {
	AddressId string
	UserId    int64
	ItemId    string
}
type SubmitOrderForActivityParam struct {
	AddressId string
	UserId    int64
	ItemId    string
	Activity  string
}
type CalculateProductResult struct {
	TotalCost int
	ItemList  []submitOrderItem
}
type ExchangeCallbackResult struct {
	BizId   string
	Credits int64
}

type AutoLoginParam struct {
	UserId   int64
	Path     string
	DCustom  string
	Transfer string
	SignKeys string
}
type AutoLoginOpenIdParam struct {
	UserId   int64
	Vip      int
	OpenId   string
	Path     string
	DCustom  string
	Transfer string
	SignKeys string
}
type BindPhoneByIVParam struct {
	UserId        int64
	EncryptedData string
	IV            string
}
type FindStepHistoryBy struct {
	OpenId  string
	Day     model.Time
	UserId  int64
	OrderBy entity.OrderByList
}
type GetPointTransactionPageListBy struct {
	UserId    int64
	AdminId   int
	Nickname  string
	OpenId    string
	Phone     string
	StartTime model.Time
	EndTime   model.Time
	Type      entity.PointTransactionType
	Types     []entity.PointTransactionType
	Offset    int
	Limit     int
}
type ExportPointTransactionListBy struct {
	UserId    int64
	Nickname  string
	OpenId    string
	Phone     string
	StartTime model.Time
	EndTime   model.Time
	Type      entity.PointTransactionType
	Types     []entity.PointTransactionType
}
type PointRecord struct {
	ID             int64                       `json:"id"`
	BalanceOfPoint int64                       `json:"balanceOfPoint"`
	Type           entity.PointTransactionType `json:"type"`
	TypeText       string                      `json:"typeText"`
	Value          int64                       `json:"value"`
	CreateTime     model.Time                  `json:"createTime"`
	AdditionalInfo string                      `json:"additionalInfo"`
	User           entity.User                 `json:"user"`
	Note           string                      `json:"note"` //操作备注
	Admin          entity.SystemAdmin          `json:"admin"`
}
type PointTransactionTypeInfo struct {
	Type     entity.PointTransactionType `json:"type"`
	TypeText string                      `json:"typeText"`
}
type FileExportRecord struct {
	entity.FileExport
	StatusText string             `json:"statusText"`
	TypeText   string             `json:"typeText"`
	Admin      entity.SystemAdmin `json:"admin"`
}
type AddFileExportParam struct {
	AdminId int                   `json:"adminId"`
	Params  string                `json:"params"`
	Type    entity.FileExportType `json:"type"`
}

type UpdateFileExportParam struct {
	Url     string                  `json:"url"`
	Status  entity.FileExportStatus `json:"status"` //1 未开始 2进行中 3导出成功 4导出失败
	Message string                  `json:"message"`
}
type FileExportStatusAndType struct {
	StatusList []FileExportStatus `json:"statusList"`
	TypeList   []FileExportType   `json:"typeList"`
}
type FileExportStatus struct {
	Status     entity.FileExportStatus `json:"status"`
	StatusText string                  `json:"statusText"`
}

type FileExportType struct {
	Type     entity.FileExportType `json:"type"`
	TypeText string                `json:"typeText"`
}

type AdminAdjustUserPointParam struct {
	OpenId string                      `binding:"required"`
	Phone  string                      `binding:"required"`
	Type   entity.PointTransactionType `binding:"oneof=SYSTEM_ADD SYSTEM_REDUCE"` //只能时 SYSTEM_ADD 和 SYSTEM_REDUCE 两种类型
	Value  int64                       `binding:"gt=0"`                           //调整积分数量必须大于0 系统会根据类型自动判断加或者减少
	Note   string                      `binding:"required"`
}
type GetPointAdjustRecordPageListParam struct {
	OpenId    string
	Phone     string
	Type      entity.PointTransactionType
	UserId    int64
	Nickname  string
	StartTime time.Time
	EndTime   time.Time
	AdminId   int
	Offset    int
	Limit     int
}
type PointAdjustRecord struct {
	ID         int                         `json:"id"`
	User       entity.User                 `json:"user"`
	Admin      entity.SystemAdmin          `json:"admin"`
	Type       entity.PointTransactionType `json:"type"`
	Note       string                      `json:"note"`
	Value      int                         `json:"value"`
	CreateTime model.Time                  `json:"createTime"`
}
type UpdateStepHistoryByEncryptedParam struct {
	OpenId        string
	EncryptedData string
	IV            string
}
type updateStepHistoryItem struct {
	Count         int
	RecordedEpoch int64
}
type GetStepHistoryListBy struct {
	OpenId            string
	StartRecordedTime model.Time // >=
	EndRecordedTime   model.Time //<=
	RecordEpochs      []int64
	OrderBy           entity.OrderByList
}
type GetStepHistoryPageListBy struct {
	OpenId            string
	StartRecordedTime model.Time // >=
	EndRecordedTime   model.Time //<=
	RecordEpochs      []int64
	OrderBy           entity.OrderByList
	Limit             int
	Offset            int
}
type CreateOrUpdateStepHistoryParam struct {
	OpenId        string
	Count         int
	RecordedTime  model.Time
	RecordedEpoch int64
}
type WeeklyHistoryInfo struct {
	AveragePerWeeklyCo2 float64             `json:"averagePerWeeklyCo2"` //平均每周减少co2
	LifeSavedCo2        float64             `json:"lifeSavedCo2"`        //累计步行减少co2
	LifeSteps           int64               `json:"lifeSteps"`           //累计步行数量
	SevenDaysCo2        float64             `json:"sevenDaysCo2"`        //最近7天减少的co2
	StepList            []WeeklyHistoryStep `json:"stepList"`            //最近7天步行历史数据
}
type WeeklyHistoryStep struct {
	Count     int        `json:"count"`
	Time      model.Time `json:"time"`
	Timestamp int64      `json:"timestamp"`
}
type FilterPointRecordOpenIds struct {
	OpenId   string
	UserId   int64
	Phone    string
	Nickname string
}

type InviteInfo struct {
	OpenId    string     `json:"openId"`
	Nickname  string     `json:"nickname"`
	AvatarUrl string     `json:"avatarUrl"`
	Time      model.Time `json:"time"`
}
type GetPartnershipPromotionListBy struct {
	Partnership  entity.PartnershipType
	PartnerShips []entity.PartnershipType
	Trigger      entity.PartnershipPromotionTrigger
}
type FindCouponBy struct {
	CouponTypeId string
}
type FindCouponTypeBy struct {
	CouponTypeId string
}
type RedeemCouponParam struct {
	OpenId    string
	CouponId  string
	OrderType entity.OrderType
}
type RedeemCouponWithTransactionParam struct {
	OpenId               string
	CouponId             string
	OrderType            entity.OrderType
	PointTransactionType entity.PointTransactionType
}
type FindProductItemBy struct {
	ItemId string
}
type SubmitOrderFromRedeemCodeParam struct {
	UserId           int64
	DefaultAddressId string
	ProductItemId    string
	Count            int
	Partnership      entity.PartnershipType
	OrderType        entity.OrderType
}
type RedeemCouponWithTransactionResult struct {
	Point   int    `json:"point"`
	OrderId string `json:"orderId"`
}
type GenerateCouponBatchParam struct {
	CouponTypeId       string
	BatchSize          int
	GenerateRedeemCode bool
}

type ProcessPromotionInformationInfo struct {
	PromotionId string     `json:"promotionId"`
	CouponId    string     `json:"couponId"`
	Time        model.Date `json:"time"`
	OpenId      string     `json:"openid"`
	OrderId     string     `json:"orderId"`
}
type FindDuiBaPointAddLogBy struct {
	OrderNum string
}
type CreateDuiBaPointAddLog struct {
	Uid           string
	Credits       int64
	Type          duibaApi.PointAddType
	OrderNum      string
	SubOrderNum   string
	Timestamp     int64
	Description   string
	Ip            string
	Sign          string
	AppKey        string
	TransactionId string
}
type UpdateDuiBaPointAddLog struct {
	ID            string
	Uid           string
	Credits       int64
	Type          duibaApi.PointAddType
	OrderNum      string
	SubOrderNum   string
	Timestamp     string
	Description   string
	Ip            string
	Sign          string
	AppKey        string
	TransactionId string
}
type UpdateUserInfoParam struct {
	UserId      int64
	Nickname    string
	Avatar      string
	Gender      entity.UserGender
	Birthday    time.Time
	PhoneNumber string
}

type UpdateUserRiskParam struct {
	UserId int64
	Risk   int
}

type FindCertificateBy struct {
	ProductItemId string
	CertificateId string
}
type GenerateBadgeParam struct {
	OpenId        string
	CertificateId string
	ProductItemId string
	OrderId       string
	Partnership   entity.PartnershipType
}
type OCRResult struct {
	WordsResult []struct {
		Words string `json:"words"`
	} `json:"words_result"`
	WordsResultNum int   `json:"words_result_num"`
	LogID          int64 `json:"log_id"`
}
type UserAccountInfo struct {
	Balance int64 `json:"balance"` //积分余额
	CertNum int64 `json:"certNum"` //证书数量
}

type CommentCount struct {
	Date    time.Time
	TotalID int64
	Total   int64
}

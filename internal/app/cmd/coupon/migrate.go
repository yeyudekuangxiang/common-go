/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package coupon

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
	"strconv"
	"time"
)

// validatorCmd represents the validator command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validator called")
		initialize.Initialize("./config-dev.ini")
		//exec("./coupon.xlsx")
		//lvmiao666()
	},
}

func lvmiao666() {
	type D struct {
		CouponID string
		ID       int64
	}
	list := make([]D, 0)
	err := app.DB.
		Raw("select coupon.coupon_id,\"user\".id from coupon INNER JOIN \"user\" on coupon.openid = \"user\".openid  where coupon_type_id = 'lvmiaoSH666' ").
		Scan(&list).Error
	if err != nil {
		log.Panicln(err)
	}
	list2 := make([]ExchangeCodeRecord, 0)
	for _, item := range list {
		id, err := defaultSonyflakeClient.NextID()
		if err != nil {
			log.Panicln(err)
		}
		list2 = append(list2, ExchangeCodeRecord{
			ExchangeRecordId:   int64(id),
			ExchangeCodeId:     0,
			ExchangeCode:       item.CouponID,
			ExchangeCodeTypeId: 432855507087130904,
			ExchangeCodeType:   "point",
			UserId:             item.ID,
			ExchangeSetting:    "300",
			ExchangeTitle:      "关注公众号得积分",
			ExchangeImage:      "https://resources.miotech.com/static/mp2c/exchange/icon/1667531583963.png",
			CreatedAt:          time.Now(),
			ExchangeBizId:      "",
		})
	}
	err = app.DB.CreateInBatches(list2, 200).Error
	if err != nil {
		log.Panicln(err)
	}
}
func init() {
	CouponCmd.AddCommand(migrateCmd)
}

var defaultSonyflakeClient = sonyflake.NewSonyflake(sonyflake.Settings{})

type OldCouponCardType struct {
	OldCouponCardTypeId string
	CouponCardType
}
type OldCouponCard struct {
	OldCouponCardId     string
	OldCouponCardTypeId string
	CouponCard
}
type OldExchangeCodeType struct {
	OldExchangeCodeTypeId string
	ExchangeCodeType
}
type OldExchangeCode struct {
	OldExchangeCodeId     string
	OldExchangeCodeTypeId string
	ExchangeCode
}

func exec(filename string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		panic(err)
	}
	rows1, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		panic(err)
	}
	rows2, err := f.GetRows(f.GetSheetName(1))
	if err != nil {
		panic(err)
	}
	/*for i, row := range rows2 {
		if i == 0 {
			continue
		}
		userId := getUser(getRow(row,8))
		f.SetCellStr(f.GetSheetName(1), "P"+strconv.Itoa(i+1), strconv.FormatInt(userId.ID, 10))
	}
	f.Save()*/
	rows3, err := f.GetRows(f.GetSheetName(2))
	if err != nil {
		panic(err)
	}
	rows4, err := f.GetRows(f.GetSheetName(3))
	if err != nil {
		panic(err)
	}
	oldCouponTypes, err := loadCouponCodeType(rows1[1:])
	if err != nil {
		panic(err)
	}
	oldCoupons, err := loadCouponCode(oldCouponTypes, rows2[1:])
	if err != nil {
		panic(err)
	}
	oldExchangeTypes, err := loadExchangeCodeType(rows3[1:])
	if err != nil {
		panic(err)
	}
	oldExchanges, records, err := loadExchangeCode(oldExchangeTypes, rows4[1:])
	if err != nil {
		panic(err)
	}

	err = app.DB.Transaction(func(tx *gorm.DB) error {
		list1 := make([]CouponCardType, 0)
		for _, ct := range oldCouponTypes {
			list1 = append(list1, ct.CouponCardType)
		}
		err = tx.CreateInBatches(list1, 100).Error
		if err != nil {
			return err
		}

		list2 := make([]CouponCard, 0)
		for _, ct := range oldCoupons {
			list2 = append(list2, ct.CouponCard)
		}
		err = tx.CreateInBatches(list2, 100).Error
		if err != nil {
			return err
		}

		list3 := make([]ExchangeCodeType, 0)
		for _, ct := range oldExchangeTypes {
			list3 = append(list3, ct.ExchangeCodeType)
		}
		err = tx.CreateInBatches(list3, 100).Error
		if err != nil {
			return err
		}

		list4 := make([]ExchangeCode, 0)
		for _, ct := range oldExchanges {
			list4 = append(list4, ct.ExchangeCode)
		}
		err = tx.CreateInBatches(list4, 100).Error
		if err != nil {
			return err
		}

		err = tx.CreateInBatches(records, 100).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
func parse(str string) sql.NullTime {
	if str == "" {
		return sql.NullTime{}
	}
	c, err := time.Parse("2006-01-02 15:04:05", str[:19])
	if err != nil {
		panic(err)
	}
	return sql.NullTime{Valid: true, Time: c}
}
func loadCouponCodeType(rows [][]string) ([]OldCouponCardType, error) {
	list := make([]OldCouponCardType, 0)
	for _, row := range rows {
		id, err := defaultSonyflakeClient.NextID()
		if err != nil {
			return nil, err
		}

		list = append(list, OldCouponCardType{
			OldCouponCardTypeId: getRow(row, 2),
			CouponCardType: CouponCardType{
				CouponCardTypeId: int64(id),
				CouponCardType:   getRow(row, 1),
				Title:            getRow(row, 0),
				Stock:            0,
				SendCount:        0,
				CardIcon:         "",
				CreatedAt:        parse(getRow(row, 4)).Time,
				UpdatedAt:        parse(getRow(row, 5)).Time,
			},
		})
	}
	return list, nil
}
func loadCouponCode(oldTypes []OldCouponCardType, rows [][]string) ([]OldCouponCard, error) {
	couponTypeIdMap := make(map[string]OldCouponCardType)
	for _, item := range oldTypes {
		couponTypeIdMap[item.OldCouponCardTypeId] = item
	}
	list := make([]OldCouponCard, 0)
	for _, row := range rows {
		id, err := defaultSonyflakeClient.NextID()
		if err != nil {
			return nil, err
		}
		oldType, ok := couponTypeIdMap[getRow(row, 2)]
		if !ok {
			return nil, errors.New(fmt.Sprintf("未查询到%s对应的新id", getRow(row, 2)))
		}
		userId, err := strconv.ParseInt(getRow(row, 15), 10, 64)
		if err != nil {
			return nil, err
		}
		list = append(list, OldCouponCard{
			OldCouponCardId:     getRow(row, 0),
			OldCouponCardTypeId: getRow(row, 2),
			CouponCard: CouponCard{
				CouponCardId:         int64(id),
				CouponCardTypeId:     oldType.CouponCardTypeId,
				Title:                getRow(row, 4),
				CardIcon:             "",
				CouponCardType:       oldType.CouponCardType.CouponCardType,
				StartTime:            parse(getRow(row, 5)),
				EndTime:              parse(getRow(row, 6)),
				Batch:                "",
				CouponCardCode:       getRow(row, 7),
				CouponCardQrcodeText: "",
				BindStatus:           2,
				UsedStatus:           2,
				UserId:               userId,
				CreatedAt:            parse(getRow(row, 9)).Time,
				UpdatedAt:            parse(getRow(row, 10)).Time,
				BizId:                "",
				UsedTime:             parse(getRow(row, 11)),
				BindTime:             parse(getRow(row, 12)),
			},
		})
	}
	return list, nil
}

func loadExchangeCodeType(rows [][]string) ([]OldExchangeCodeType, error) {
	list := make([]OldExchangeCodeType, 0)
	for _, row := range rows {
		id, err := defaultSonyflakeClient.NextID()
		if err != nil {
			return nil, err
		}

		list = append(list, OldExchangeCodeType{
			OldExchangeCodeTypeId: getRow(row, 1),
			ExchangeCodeType: ExchangeCodeType{
				ExchangeCodeTypeId: int64(id),
				ExchangeType:       getRow(row, 3),
				ExchangeSetting:    getRow(row, 9),
				Title:              getRow(row, 0),
				Desc:               getRow(row, 11),
				SystemAdminId:      0,
				CreatedAt:          parse(getRow(row, 4)).Time,
				UpdatedAt:          parse(getRow(row, 5)).Time,
				ExchangeGoodTitle:  getRow(row, 6),
				ExchangeGoodImage:  "",
				Stock:              0,
				ExchangeCount:      0,
				Note:               "",
			},
		})
	}
	return list, nil
}
func getRow(row []string, index int) string {
	if len(row) > index {
		return row[index]
	}
	return ""
}
func loadExchangeCode(oldTypes []OldExchangeCodeType, rows [][]string) ([]OldExchangeCode, []ExchangeCodeRecord, error) {
	exchangeMap := make(map[string]OldExchangeCodeType)
	for _, item := range oldTypes {
		exchangeMap[item.OldExchangeCodeTypeId] = item
	}

	list := make([]OldExchangeCode, 0)
	records := make([]ExchangeCodeRecord, 0)
	for _, row := range rows {
		id, err := defaultSonyflakeClient.NextID()
		if err != nil {
			return nil, nil, err
		}

		exchangeType, ok := exchangeMap[getRow(row, 1)]
		if !ok {
			return nil, nil, errors.New(fmt.Sprintf("未查询到%s对应的新id", getRow(row, 1)))
		}

		usedStatus := 1
		if getRow(row, 5) == "t" {
			usedStatus = 2
		}
		userId, err := strconv.ParseInt(getRow(row, 7), 10, 64)
		if err != nil {
			return nil, nil, err
		}

		list = append(list, OldExchangeCode{
			OldExchangeCodeId:     getRow(row, 0),
			OldExchangeCodeTypeId: getRow(row, 1),
			ExchangeCode: ExchangeCode{
				ExchangeCodeId:     int64(id),
				ExchangeCodeTypeId: exchangeType.ExchangeCodeTypeId,
				ExchangeCode:       getRow(row, 2),
				StartTime:          parse(getRow(row, 3)),
				EndTime:            parse(getRow(row, 4)),
				UsedStatus:         int64(usedStatus),
				CreatedAt:          parse(getRow(row, 6)).Time,
				UserId:             userId,
				ExchangeTime:       parse(getRow(row, 8)),
			},
		})

		if usedStatus == 2 {
			id2, err := defaultSonyflakeClient.NextID()
			if err != nil {
				return nil, nil, err
			}
			records = append(records, ExchangeCodeRecord{
				ExchangeRecordId:   int64(id2),
				ExchangeCodeId:     int64(id),
				ExchangeCode:       getRow(row, 2),
				ExchangeCodeTypeId: exchangeType.ExchangeCodeTypeId,
				ExchangeCodeType:   exchangeType.ExchangeType,
				UserId:             userId,
				ExchangeSetting:    exchangeType.ExchangeSetting,
				ExchangeTitle:      exchangeType.ExchangeGoodTitle,
				ExchangeImage:      exchangeType.ExchangeGoodImage,
				CreatedAt:          parse(getRow(row, 6)).Time,
				ExchangeBizId:      "",
			})
		}
	}
	return list, records, nil
}

type CouponCardType struct {
	CouponCardTypeId int64     `gorm:"primaryKey;column:coupon_card_type_id"` // 券id
	CouponCardType   string    `gorm:"column:coupon_card_type"`               // 券类型 direct(直充) third(三方) lvmiao(绿喵)
	Title            string    `gorm:"column:title"`                          // 券名称
	Stock            int64     `gorm:"column:stock"`                          // 券库存
	SendCount        int64     `gorm:"column:send_count"`                     // 已发送数量
	CardIcon         string    `gorm:"column:card_icon"`                      // 券图标
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}
type CouponCard struct {
	CouponCardId         int64        `gorm:"primaryKey;column:coupon_card_id"`
	CouponCardTypeId     int64        `gorm:"column:coupon_card_type_id"`     // 券类型id
	Title                string       `gorm:"column:title"`                   // 券标题
	CardIcon             string       `gorm:"column:card_icon"`               // 券图标
	CouponCardType       string       `gorm:"column:coupon_card_type"`        // 券类型 direct(直充) third(三方) lvmiao(绿喵)
	StartTime            sql.NullTime `gorm:"column:start_time"`              // 券有效期开始时间
	EndTime              sql.NullTime `gorm:"column:end_time"`                // 券有效期结束时间
	Batch                string       `gorm:"column:batch"`                   // 优惠券导入批次 20060102150405000
	CouponCardCode       string       `gorm:"column:coupon_card_code"`        // 优惠券券码
	CouponCardQrcodeText string       `gorm:"column:coupon_card_qrcode_text"` // 优惠券核销二维码字符串
	BindStatus           int64        `gorm:"column:bind_status"`             // 绑定用户状态 1未发给用户 2已发送给用户
	UsedStatus           int64        `gorm:"column:used_status"`             // 1未使用 2已使用
	UserId               int64        `gorm:"column:user_id"`                 // 用户id
	CreatedAt            time.Time    `gorm:"column:created_at"`              // 创建时间
	UpdatedAt            time.Time    `gorm:"column:updated_at"`              // 更新时间
	BizId                string       `gorm:"column:biz_id"`                  // 业务id
	UsedTime             sql.NullTime `gorm:"column:used_time"`               // 使用时间
	BindTime             sql.NullTime `gorm:"column:bind_time"`               // 绑定时间
}
type ExchangeCode struct {
	ExchangeCodeId     int64        `gorm:"primaryKey;column:exchange_code_id"` // 兑换码id
	ExchangeCodeTypeId int64        `gorm:"column:exchange_code_type_id"`       // 兑换码id
	ExchangeCode       string       `gorm:"column:exchange_code"`               // 兑换码
	StartTime          sql.NullTime `gorm:"column:start_time"`                  // 开始时间
	EndTime            sql.NullTime `gorm:"column:end_time"`                    // 结束时间
	UsedStatus         int64        `gorm:"column:used_status"`                 // 使用状态 1未使用 2已使用
	CreatedAt          time.Time    `gorm:"column:created_at"`                  // 创建时间
	UserId             int64        `gorm:"column:user_id"`                     // 兑换者id
	ExchangeTime       sql.NullTime `gorm:"column:exchange_time"`               // 兑换时间
}
type ExchangeCodeRecord struct {
	ExchangeRecordId   int64     `gorm:"primaryKey;column:exchange_record_id"`
	ExchangeCodeId     int64     `gorm:"column:exchange_code_id"`      // 兑换码id
	ExchangeCode       string    `gorm:"column:exchange_code"`         // 兑换码
	ExchangeCodeTypeId int64     `gorm:"column:exchange_code_type_id"` // 兑换码类型id
	ExchangeCodeType   string    `gorm:"column:exchange_code_type"`    // 兑换码类型
	UserId             int64     `gorm:"column:user_id"`               // 用户id
	ExchangeSetting    string    `gorm:"column:exchange_setting"`      // 兑换物id
	ExchangeTitle      string    `gorm:"column:exchange_title"`        // 兑换物标题
	ExchangeImage      string    `gorm:"column:exchange_image"`        // 兑换物图片
	CreatedAt          time.Time `gorm:"column:created_at"`            // 创建时间(兑换时间)
	ExchangeBizId      string    `gorm:"column:exchange_biz_id"`       // 业务id
}
type ExchangeCodeType struct {
	ExchangeCodeTypeId int64     `gorm:"primaryKey;column:exchange_code_type_id"` // 兑换码类型id
	ExchangeType       string    `gorm:"column:exchange_type"`                    // 兑换码类型 product、point、coupon、certificate
	ExchangeSetting    string    `gorm:"column:exchange_setting"`                 // 兑换物品配置信息
	Title              string    `gorm:"column:title"`                            // 兑换码名称
	Desc               string    `gorm:"column:desc"`                             // 描述
	SystemAdminId      int64     `gorm:"column:system_admin_id"`                  // 管理员id
	CreatedAt          time.Time `gorm:"column:created_at"`                       // 创建时间
	UpdatedAt          time.Time `gorm:"column:updated_at"`                       // 更新时间
	ExchangeGoodTitle  string    `gorm:"column:exchange_good_title"`              // 兑换物品的名称
	ExchangeGoodImage  string    `gorm:"column:exchange_good_image"`              // 兑换物品的图片
	Stock              int64     `gorm:"column:stock"`                            // 库存
	ExchangeCount      int64     `gorm:"column:exchange_count"`                   // 兑换数量
	Note               string    `gorm:"column:note"`                             // 管理员备注
}

var usermap = make(map[string]entity.User)

func getUser(openid string) entity.User {
	if u, ok := usermap[openid]; ok {
		return u
	}
	user := entity.User{}
	err := app.DB.Where("openid = ? and source = 'mio'", openid).Take(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	usermap[openid] = user
	return user
}

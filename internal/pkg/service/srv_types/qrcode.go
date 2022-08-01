package srv_types

import (
	"mio/internal/pkg/model/entity"
)

type CreateQrCodeDTO struct {
	ImagePath    string              //oss 路径
	QrCodeScene  entity.QrCodeScene  //使用场景
	QrCodeSource entity.QrCodeSource //来源
	Key          string              //key 和 QrCodeScene 组成唯一索引
	Content      string              //二维码内容
	Ext          string              //额外参数
	OpenId       string              //用于记录生成用户的id
	Description  string              //二维码描述
}

package business

import (
	"encoding/json"
	ebusiness "mio/internal/pkg/model/entity/business"
)

var DefaultPointSettingService = PointSettingService{}

type PointSettingService struct {
}

func (srv PointSettingService) ParseSaveWaterElectricityRate(setting ebusiness.PointSetting) (*ebusiness.SaveWaterElectricityExchangeRate, error) {
	rate := ebusiness.SaveWaterElectricityExchangeRate{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointSettingService) ParsePublicTransportRate(setting ebusiness.PointSetting) (*ebusiness.PublicTransportExchangeRate, error) {
	rate := ebusiness.PublicTransportExchangeRate{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointSettingService) ParseOnlineMeetingRate(setting ebusiness.PointSetting) (*ebusiness.OnlineMeetingExchangeRate, error) {
	rate := ebusiness.OnlineMeetingExchangeRate{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointSettingService) ParsePointExchangeRate(setting ebusiness.PointSetting) (*ebusiness.PointExchangeRate, error) {
	rate := ebusiness.PointExchangeRate{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointSettingService) EncodePointExchangeRate(rate ebusiness.PointSettingRate) ebusiness.PointSetting {
	data, err := json.Marshal(rate)
	if err != nil {
		panic(err)
	}
	return ebusiness.PointSetting(data)
}

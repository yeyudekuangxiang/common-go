package business

import (
	"encoding/json"
	ebusiness "mio/internal/pkg/model/entity/business"
)

var DefaultPointRateSettingService = PointRateSettingService{}

type PointRateSettingService struct {
}

func (srv PointRateSettingService) ParseSaveWaterElectricityRate(setting ebusiness.PointRateSetting) (*ebusiness.PointRateSaveWaterElectricity, error) {
	rate := ebusiness.PointRateSaveWaterElectricity{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointRateSettingService) ParsePublicTransportRate(setting ebusiness.PointRateSetting) (*ebusiness.PointRatePublicTransport, error) {
	rate := ebusiness.PointRatePublicTransport{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointRateSettingService) ParseOnlineMeetingRate(setting ebusiness.PointRateSetting) (*ebusiness.PointRateOnlineMeeting, error) {
	rate := ebusiness.PointRateOnlineMeeting{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointRateSettingService) ParsePointExchangeRate(setting ebusiness.PointRateSetting) (*ebusiness.PointRate, error) {
	rate := ebusiness.PointRate{}
	err := json.Unmarshal([]byte(setting), &rate)
	return &rate, err
}
func (srv PointRateSettingService) EncodePointExchangeRate(rate ebusiness.PointRate) ebusiness.PointRateSetting {
	data, err := json.Marshal(rate)
	if err != nil {
		panic(err)
	}
	return ebusiness.PointRateSetting(data)
}

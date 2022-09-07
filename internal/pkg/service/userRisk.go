package service

import (
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
)

func NewUserRiskService(ctx *context.MioContext) UserRiskService {
	return UserRiskService{
		ctx:      ctx,
		r:        repository.NewUserRepository(),
		rCity:    repository.NewCityRepository(ctx),
		rChannel: repository.UserChannelRepository{DB: app.DB},
		rPoint:   repository.NewPointRepository(ctx),
	}
}

type UserRiskService struct {
	ctx      *context.MioContext
	r        repository.UserRepository
	rCity    repository.CityRepository
	rChannel repository.UserChannelRepository
	rPoint   *repository.PointRepository
}

func (u UserRiskService) BatchUpdateUserRisk(param UpdateRiskParam) error {
	err := u.r.BatchUpdateUserRisk(repository.UpdateUserRisk{
		UserIdSlice: param.UserIdSlice,
		OpenIdSlice: param.OpenIdSlice,
		PhoneSlice:  param.PhoneSlice,
		Risk:        param.Risk,
	})
	return err
}

func (u UserRiskService) GetUserPageListBy(by repository.GetUserPageListBy) ([]entity.User, int64) {
	return u.r.GetUserPageListBy(by)
}

func (u UserRiskService) GetUserRiskPageListBy(by repository.GetUserPageListBy) ([]api_types.UserVO, int64) {
	list, total := u.r.GetUserPageListBy(by)
	var cidSlice []int64
	var openidSlice, citySlice []string
	cityMap := make(map[string]entity.City)
	pointMap := make(map[string]entity.Point)
	channelMap := make(map[int64]entity.UserChannel)

	for _, user := range list {
		cidSlice = append(cidSlice, user.ChannelId)
		openidSlice = append(openidSlice, user.OpenId)
		citySlice = append(citySlice, user.CityCode)
	}

	//获取渠道信息
	if len(cidSlice) != 0 {
		channelList, _ := u.rChannel.GetUserChannelPageList(repository.GetUserChannelPageListBy{CidSlice: cidSlice})
		for _, channel := range channelList {
			channelMap[channel.Cid] = channel
		}
	}

	//获取积分信息
	if len(openidSlice) != 0 {
		pointList := u.rPoint.FindListPoint(repository.FindListPoint{OpenIds: openidSlice})
		for _, point := range pointList {
			pointMap[point.OpenId] = point
		}
	}

	//城市信息
	if len(citySlice) != 0 {
		cityList, _ := u.rCity.GetList(repotypes.GetCityListDO{CityCodeSlice: citySlice})
		for _, city := range cityList {
			cityMap[city.CityCode] = city
		}
	}

	userVoList := make([]api_types.UserVO, 0)
	for _, l := range list {

		//初始化
		var balance int64
		var cityName, channelName string

		//积分值
		point, ok := pointMap[l.OpenId]
		if ok {
			balance = point.Balance
		}

		//城市名
		city, ok2 := cityMap[l.CityCode]
		if ok2 {
			cityName = city.Name
		}

		//渠道名
		channel, ok3 := channelMap[l.ChannelId]
		if ok3 {
			channelName = channel.Name
		}

		//整理
		userVoList = append(userVoList, api_types.UserVO{
			OpenId:      l.OpenId,
			Nickname:    l.Nickname,
			AvatarUrl:   l.AvatarUrl,
			PhoneNumber: l.PhoneNumber,
			Risk:        l.Risk,
			Point:       balance,
			CityName:    cityName,
			ChannelName: channelName,
			RegTime:     l.Time.Format("2006-01-02"),
		})
	}
	return userVoList, total
}

func (u UserRiskService) GetUserRiskStatisticst() []repository.RiskStatistics {
	list := u.r.GetRiskStatistics()
	return list
}

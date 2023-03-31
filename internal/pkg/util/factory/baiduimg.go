package factory

import (
	"context"
	"gitlab.miotech.com/miotech-application/backend/common-go/baidu"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
)

func NewBaiDuImageFromTokenCenterRpc(centerId string, centerRpc tokencenterclient.TokenCenter) *baidu.ImageClient {
	return baidu.NewImageClient(&bdTokenCenter{
		tokenCenterRpc: centerRpc,
		centerId:       centerId,
	})
}
func NewBaiDuReviewFromTokenCenterRpc(centerId string, centerRpc tokencenterclient.TokenCenter) *baidu.ReviewClient {
	return baidu.NewReviewClient(&bdTokenCenter{
		tokenCenterRpc: centerRpc,
		centerId:       centerId,
	})
}

type bdTokenCenter struct {
	tokenCenterRpc tokencenterclient.TokenCenter
	centerId       string
}

func (b bdTokenCenter) GetToken() (string, error) {
	resp, err := b.tokenCenterRpc.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: b.centerId,
		OldToken: "",
		Refresh:  false,
	})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

func (b bdTokenCenter) RefreshToken() (string, error) {
	resp, err := b.tokenCenterRpc.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: b.centerId,
		OldToken: "",
		Refresh:  true,
	})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

package starcharge

import (
	"net/http"
	"testing"
)

func TestBikeCard(t *testing.T) {
	c := Client{
		htpClient:  http.Client{},
		Domain:     "https://openapi.hellobike.com/bike/activity",
		Version:    "",
		AESSecret:  "",
		SigSecret:  "",
		Token:      "",
		OperatorID: "",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, _ := c.Request(SendStarChargeParam{
		Data:     "",
		QueryUrl: "",
	})

	println(resp.Ret)

}

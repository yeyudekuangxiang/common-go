package api

import (
	bytes2 "bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mio/tests/mp2c"
	"net/http/httptest"
	"testing"
)

func TestDuiBaOrderCallback(t *testing.T) {
	router := mp2c.SetupServer()
	recorder := httptest.NewRecorder()

	form := map[string]string{
		"account":          "13000000000",
		"appKey":           "ngiGp48EcRUC9TjpXEYxdSSJhim",
		"consumerPayPrice": "199.99",
		"createTime":       "1651889176176",
		"developBizId":     "develop-biz-id",
		"errorMsg":         "",
		"expressPrice":     "9.99",
		"finishTime":       "1651889176176",
		"orderItemList":    "[{\"title\":\"商品标题0\",\"isSelf\":\"0\",\"smallImage\":\"smallImage\",\"merchantCode\":\"merchantCode-ads-0\",\"perCredit\":\"100\",\"perPrice\":\"9.99\",\"quantity\":\"1\",\"code\":\"code0\",\"password\":\"password0\",\"cardBeginTime\":\"1651889176176\",\"cardEndTime\":\"1652249176176\",\"deliveryCompanyNo\":\"eexpress-no-0\",\"deliveryCompanyName\":\"eexpress-name-0\",\"duibaSupplyPrice\":\"12.22\"},{\"title\":\"商品标题1\",\"isSelf\":\"1\",\"smallImage\":\"smallImage\",\"merchantCode\":\"merchantCode-ads-1\",\"perCredit\":\"100\",\"perPrice\":\"9.99\",\"quantity\":\"2\",\"code\":\"code1\",\"password\":\"password1\",\"cardBeginTime\":\"1651889176176\",\"cardEndTime\":\"1652249176176\",\"deliveryCompanyNo\":\"eexpress-no-1\",\"deliveryCompanyName\":\"eexpress-name-1\",\"duibaSupplyPrice\":\"12.22\"},{\"title\":\"商品标题2\",\"isSelf\":\"0\",\"smallImage\":\"smallImage\",\"merchantCode\":\"merchantCode-ads-2\",\"perCredit\":\"100\",\"perPrice\":\"9.99\",\"quantity\":\"3\",\"code\":\"code2\",\"password\":\"password2\",\"cardBeginTime\":\"1651889176176\",\"cardEndTime\":\"1652249176176\",\"deliveryCompanyNo\":\"eexpress-no-2\",\"deliveryCompanyName\":\"eexpress-name-2\",\"duibaSupplyPrice\":\"12.22\"}]",
		"orderNum":         "order-num",
		"orderStatus":      "afterSend",
		"receiveAddrInfo":  "{\"province\":\"上海市\",\"city\":\"上海市\",\"area\":\"浦东新区\",\"street\":\"上海中心大厦\",\"address\":\"1402\",\"mobile\":\"13000000000\",\"name\":\"绿喵\"}",
		"sign":             "caf427617f267d582bbe4dc84424f3e2",
		"source":           "普兑",
		"timestamp":        "1651889176176",
		"totalCredits":     "10000",
		"type":             "object",
		"uid":              "greencat",
	}
	body, _ := json.Marshal(form)

	request := httptest.NewRequest("POST", "/api/mp2c/duiba/order/callback", bytes2.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code)

	assert.Equal(t, "ok", recorder.Body.String())
}

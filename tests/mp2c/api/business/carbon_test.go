package business

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mio/tests/mp2c"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetCarbonCreditLogInfoList(t *testing.T) {
	router := mp2c.SetupServer()
	recorder := httptest.NewRecorder()

	form := url.Values{
		"date": {"2022-04"},
	}

	request := httptest.NewRequest("GET", "/api/mp2c/business/carbon/record/list?"+form.Encode(), nil)
	mp2c.AddBusinessToken(request)

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	t.Logf("TestGetPointRecordList %s", recorder.Body.String())

	resp := mp2c.Error{}
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.Code)
}

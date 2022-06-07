package business

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mio/tests/mp2c"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetCarbonUserRankList(t *testing.T) {
	router := mp2c.SetupServer()
	recorder := httptest.NewRecorder()

	form := url.Values{
		"dateType": {"day"},
		"page":     {"1"},
		"pageSize": {"10"},
	}

	request := httptest.NewRequest("GET", "/api/mp2c/business/carbon/rank/user/list?"+form.Encode(), nil)
	mp2c.AddBusinessToken(request)

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	t.Logf("TestGetCarbonUserRankList %s", recorder.Body.String())

	resp := mp2c.Error{}
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.Code)
}
func TestGetCarbonDepartmentRankList(t *testing.T) {
	router := mp2c.SetupServer()
	recorder := httptest.NewRecorder()

	form := url.Values{
		"dateType": {"day"},
		"page":     {"1"},
		"pageSize": {"10"},
	}

	request := httptest.NewRequest("GET", "/api/mp2c/business/carbon/rank/department/list?"+form.Encode(), nil)
	mp2c.AddBusinessToken(request)

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	t.Logf("TestGetCarbonDepartmentRankList %s", recorder.Body.String())

	resp := mp2c.Error{}
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.Code)
}
func TestRankChangeLikeStatus(t *testing.T) {
	router := mp2c.SetupServer()
	recorder := httptest.NewRecorder()

	form := url.Values{
		"dateType": {"day"},
		"page":     {"1"},
		"pageSize": {"10"},
	}

	request := httptest.NewRequest("POST", "/api/mp2c/business/carbon/rank/like/status/change", strings.NewReader(form.Encode()))
	mp2c.AddBusinessToken(request)

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code)
	t.Logf("TestRankChangeLikeStatus %s", recorder.Body.String())

	resp := mp2c.Error{}
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.Code)
}

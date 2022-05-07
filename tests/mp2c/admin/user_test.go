package admin

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mio/tests/mp2c"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	router := mp2c.SetupServer()
	param := url.Values{
		"id": {"1"},
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/admin/user", nil)
	request.Form = param
	mp2c.AddAdminToken(request)
	router.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	var res mp2c.Response
	_ = json.NewDecoder(recorder.Body).Decode(&res)
	bytes, _ := json.Marshal(res)
	t.Logf("%+v", string(bytes))
	assert.Equal(t, 200, res.Code, res.Message)
}

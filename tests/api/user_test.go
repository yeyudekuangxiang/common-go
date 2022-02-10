package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"mio/tests"
	"net/http/httptest"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	tests.SetupMock()
	router := tests.SetupServer()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/user", nil)
	tests.AddUserToken(request)

	router.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	var res tests.Response
	_ = json.NewDecoder(recorder.Body).Decode(&res)
	bytes, _ := json.Marshal(res)
	t.Logf("%+v", string(bytes))
	assert.Equal(t, 200, res.Code, res.Message)
}

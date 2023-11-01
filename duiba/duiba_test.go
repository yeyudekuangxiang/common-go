package duiba

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"log"
	"os"
	"testing"
)

func init() {
	httptool.SetLogger(httptool.FuncLogger(func(data httptool.LogData, err error) {
		log.Printf("%+v %+v %+v", data.Url, string(data.RequestBody), string(data.ResponseBody))
	}))
}
func TestClient_AddActivityTimes(t *testing.T) {
	appKey := os.Getenv("appKey")
	appSecret := os.Getenv("appSecret")
	if appKey == "" || appSecret == "" {
		return
	}
	client := NewClient(appKey, appSecret)
	log.Printf("%+v\n", client)
	resp, err := client.AddActivityTimes(AddActivityTimesParam{
		ActivityId: "243689033869996",
		Times:      1,
		ValidType:  1,
		Uid:        "oy_BA5DGmQBqMeCj_9Eozj8dXhoA",
		BizId:      "1231321kjjjjbjkbjk",
	})
	require.Equal(t, nil, err)
	assert.Equal(t, true, resp.Success, resp.Message)
}

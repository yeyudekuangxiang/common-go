package tjmetro

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSign(t *testing.T) {
	c := Client{
		Domain:  "https://app.trtpazyz.com",
		AppId:   "264735a59163453d9772f92e1f703123",
		Version: "1.0",
	}
	resp, _, _ := c.TicketAllot(TicketAllotParam{
		AllotId:     "12312",
		EtUserId:    "12312",
		EtUserPhone: "123123",
		AllotNum:    1,
	})
	assert.Equal(t, "0001", resp.ResultCode)
}

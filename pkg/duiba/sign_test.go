package duiba

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSign(t *testing.T) {
	params := map[string]string{
		"appKey":     "testappKey",
		"appSecret":  "testappSecret",
		"uid=test":   "test",
		"credits":    "100",
		"timestamp ": "1520559858580",
	}
	assert.Equal(t, "49b12bc5579a2a2a4652a68cd53c1e5e", sign(params))
}

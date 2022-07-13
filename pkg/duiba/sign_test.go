package duiba

import (
	"github.com/stretchr/testify/assert"
	"mio/pkg/duiba/util"
	"testing"
)

func TestMd5(t *testing.T) {
	assert.Equal(t, "c4ca4238a0b923820dcc509a6f75849b", util.Md5("1"))
	assert.Equal(t, "900150983cd24fb0d6963f7d28e17f72", util.Md5("abc"))
	assert.Equal(t, "2c21806ce0d8f171850e27d32735605e", util.Md5("qshoqwd[pnkxqod1234567890.()-=!@#$%^&*()_+-=;':?/.,"))
}
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

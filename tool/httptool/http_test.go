package httptool

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestGet(t *testing.T) {
	_, err := Get("https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&rsv_idx=1&tn=baidu&wd=nihao&fenlei=256&rsv_pq=0xbbc238bb000021db&rsv_t=ff8cWxNk5LAZjIYf5SY2dDHrQMdLh2UOyB1oPKWHH%2FJzgxhFNN7iiSzNfOuc&rqlang=en&rsv_dl=tb&rsv_enter=1&rsv_sug3=5&rsv_sug1=3&rsv_sug7=100&rsv_sug2=0&rsv_btype=i&prefixsug=nihao&rsp=5&inputT=801&rsv_sug4=801")
	assert.Equal(t, nil, err)
}
func TestPostJson(t *testing.T) {
	_, err := PostJson("https://www.baidu.com", map[string]string{
		"userId": "1",
		"name":   "haha",
	})
	assert.Equal(t, nil, err)
}
func TestPostMapForm(t *testing.T) {
	_, err := PostMapFrom("https://www.baidu.com", map[string]string{
		"userId": "1",
		"name":   "haha",
	})
	assert.Equal(t, nil, err)
}
func TestPostForm(t *testing.T) {
	_, err := PostFrom("https://www.baidu.com", url.Values{
		"userId": {"1"},
		"name":   {"haha"},
	})
	assert.Equal(t, nil, err)
}

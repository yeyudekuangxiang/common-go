package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	m := map[string]string{
		"name": "test",
		"sex":  "1",
		"age":  "18",
	}
	assert.Equal(t, "age=18&name=test&sex=1", BuildQuery(m))
}
func TestMapTo(t *testing.T) {
	s := struct {
		Name string `json:"name"`
		Sex  int    `json:"sex"`
		Age  int    `json:"age"`
	}{
		Name: "test",
		Sex:  1,
		Age:  18,
	}

	m := make(map[string]interface{})
	err := MapTo(s, &m)
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", m["name"])
	assert.Equal(t, 1, m["sex"])
	assert.Equal(t, 18, m["age"])
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "c4ca4238a0b923820dcc509a6f75849b", Md5("1"))
	assert.Equal(t, "900150983cd24fb0d6963f7d28e17f72", Md5("abc"))
	assert.Equal(t, "2c21806ce0d8f171850e27d32735605e", Md5("qshoqwd[pnkxqod1234567890.()-=!@#$%^&*()_+-=;':?/.,"))
}
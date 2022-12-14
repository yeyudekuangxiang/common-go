package number

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustParseLongInt64(t *testing.T) {
	assert.Equal(t, int64(12345678987654321), MustParseLongInt64("12345678987654321").Int64())
}
func TestLongInt64_MarshalJSON(t *testing.T) {
	d, err := LongInt64(12345678987654321).MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, "\"12345678987654321\"", string(d))
}
func TestLongInt64_Scan(t *testing.T) {
	t2 := new(LongInt64)
	assert.Equal(t, nil, t2.Scan("12345678987654321"))
	assert.Equal(t, int64(12345678987654321), t2.Int64())
	assert.NotEqual(t, nil, t2.Scan("\"12345678987654321"))
	assert.NotEqual(t, nil, t2.Scan(false))
}
func TestLongInt64_Value(t *testing.T) {
	t2 := LongInt64(12345678987654321)
	v, err := t2.Value()
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(12345678987654321), v)
	LongInt64GormTypeString = true
	v, err = t2.Value()
	assert.Equal(t, nil, err)
	assert.Equal(t, "12345678987654321", v)
}

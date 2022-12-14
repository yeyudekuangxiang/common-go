package number

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustParseLongUInt64(t *testing.T) {
	assert.Equal(t, uint64(12345678987654321), MustParseLongUInt64("12345678987654321").UInt64())
}
func TestLongUInt64_MarshalJSON(t *testing.T) {
	d, err := LongUInt64(12345678987654321).MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, "\"12345678987654321\"", string(d))
}
func TestLongUInt64_Scan(t *testing.T) {
	t2 := new(LongUInt64)
	assert.Equal(t, nil, t2.Scan("12345678987654321"))
	assert.Equal(t, uint64(12345678987654321), t2.UInt64())
	assert.NotEqual(t, nil, t2.Scan("\"12345678987654321"))
	assert.NotEqual(t, nil, t2.Scan(false))
}
func TestLongUInt64_Value(t *testing.T) {
	t2 := LongUInt64(12345678987654321)
	v, err := t2.Value()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(12345678987654321), v)
	LongUInt64GormTypeString = true
	v, err = t2.Value()
	assert.Equal(t, nil, err)
	assert.Equal(t, "12345678987654321", v)
}

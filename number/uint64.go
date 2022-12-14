package number

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

// LongUInt64JSONTypeString json处理时是否以string处理
var LongUInt64JSONTypeString = true

// LongUInt64GormTypeString gorm 处理时是否以string处理
var LongUInt64GormTypeString = false

type LongUInt64 uint64

// UInt64 returns an int64 of the LongUInt64
func (f LongUInt64) UInt64() uint64 {
	return uint64(f)
}

// String returns a string of the LongUInt64
func (f LongUInt64) String() string {
	return strconv.FormatUint(uint64(f), 10)
}

// Base2 returns a string base2 of the LongUInt64
func (f LongUInt64) Base2() string {
	return strconv.FormatUint(uint64(f), 2)
}

// ParseLongUInt64Base2 converts a Base2 string into a LongUInt64
func ParseLongUInt64Base2(id string) (LongUInt64, error) {
	i, err := strconv.ParseUint(id, 2, 64)
	return LongUInt64(i), err
}

// ParseLongUInt64 converts a string into a LongUInt64
func ParseLongUInt64(id string) (LongUInt64, error) {
	i, err := strconv.ParseUint(id, 10, 64)
	return LongUInt64(i), err
}

// MustParseLongUInt64 must converts a string into a LongUInt64
func MustParseLongUInt64(id string) LongUInt64 {
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		panic(err)
	}
	return LongUInt64(i)
}

// MarshalJSON returns a json byte array string of the LongUInt64 ID.
func (f LongUInt64) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	if LongUInt64JSONTypeString {
		buff = append(buff, '"')
	}
	buff = strconv.AppendUint(buff, uint64(f), 10)
	if LongUInt64JSONTypeString {
		buff = append(buff, '"')
	}
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a LongUInt64 ID into an ID type.
func (f *LongUInt64) UnmarshalJSON(b []byte) error {
	var i uint64
	var err error
	if LongUInt64JSONTypeString {
		if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
			return errors.Errorf("json syntax error %q", string(b))
		}
		i, err = strconv.ParseUint(string(b[1:len(b)-1]), 10, 64)
		if err != nil {
			return err
		}
	} else {
		if len(b) < 1 {
			return errors.Errorf("json syntax error %q", string(b))
		}
		i, err = strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return err
		}
	}
	*f = LongUInt64(i)
	return nil
}

func (f *LongUInt64) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		id, err := strconv.ParseUint(value.(string), 10, 64)
		if err != nil {
			return err
		}
		*f = LongUInt64(id)
		return nil
	case uint64:
		*f = LongUInt64(value.(int64))
		return nil
	}
	return errors.New(fmt.Sprintf("unsupport SonyflakeID type %v", value))
}

func (f LongUInt64) Value() (driver.Value, error) {
	if LongUInt64GormTypeString {
		return f.String(), nil
	}
	return f.UInt64(), nil
}

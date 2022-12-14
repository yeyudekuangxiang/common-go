package number

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

// LongInt64JSONTypeString json处理时是否以string处理
var LongInt64JSONTypeString = true

// LongInt64GormTypeString gorm 处理时是否以string处理
var LongInt64GormTypeString = false

type LongInt64 int64

// Int64 returns an int64 of the LongInt64
func (f LongInt64) Int64() int64 {
	return int64(f)
}

// String returns a string of the LongInt64
func (f LongInt64) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// Base2 returns a string base2 of the LongInt64
func (f LongInt64) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// ParseLongInt64Base2 converts a Base2 string into a LongInt64
func ParseLongInt64Base2(id string) (LongInt64, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return LongInt64(i), err
}

// ParseLongInt64 converts a string into a LongInt64
func ParseLongInt64(id string) (LongInt64, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return LongInt64(i), err
}

// MustParseLongInt64 must converts a string into a LongInt64
func MustParseLongInt64(id string) LongInt64 {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}
	return LongInt64(i)
}

// MarshalJSON returns a json byte array string of the LongInt64 ID.
func (f LongInt64) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	if LongInt64JSONTypeString {
		buff = append(buff, '"')
	}
	buff = strconv.AppendInt(buff, int64(f), 10)
	if LongInt64JSONTypeString {
		buff = append(buff, '"')
	}
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a LongInt64 ID into an ID type.
func (f *LongInt64) UnmarshalJSON(b []byte) error {
	var i int64
	var err error
	if LongInt64JSONTypeString {
		if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
			return errors.Errorf("json syntax error %q", string(b))
		}
		i, err = strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
		if err != nil {
			return err
		}
	} else {
		if len(b) < 1 {
			return errors.Errorf("json syntax error %q", string(b))
		}
		i, err = strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
	}
	*f = LongInt64(i)
	return nil
}

func (f *LongInt64) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		id, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return err
		}
		*f = LongInt64(id)
		return nil
	case int64:
		*f = LongInt64(value.(int64))
		return nil
	default:
		return errors.New(fmt.Sprintf("unsupport LongInt64 type %v", value))
	}
}

func (f LongInt64) Value() (driver.Value, error) {
	if LongInt64GormTypeString {
		return f.String(), nil
	}
	return f.Int64(), nil
}

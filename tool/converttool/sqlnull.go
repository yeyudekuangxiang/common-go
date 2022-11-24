package converttool

import (
	"database/sql"
	"time"
)

func NullBool2Pointer(val sql.NullBool) *bool {
	if !val.Valid {
		return nil
	}
	return PointerBool(val.Bool)
}
func NullInt642Pointer(val sql.NullInt64) *int64 {
	if !val.Valid {
		return nil
	}
	return PointerInt64(val.Int64)
}
func NullInt322Pointer(val sql.NullInt32) *int32 {
	if !val.Valid {
		return nil
	}
	return PointerInt32(val.Int32)
}
func NullString2Pointer(val sql.NullString) *string {
	if !val.Valid {
		return nil
	}
	return PointerString(val.String)
}
func NullTime2Pointer(val sql.NullTime) *time.Time {
	if !val.Valid {
		return nil
	}
	return PointerTime(val.Time)
}
func NullFloat642Pointer(val sql.NullFloat64) *float64 {
	if !val.Valid {
		return nil
	}
	return PointerFloat64(val.Float64)
}

func Pointer2NullBool(val *bool) sql.NullBool {
	if val == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Valid: true, Bool: *val}
}
func Pointer2NullInt64(val *int64) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Valid: true, Int64: *val}
}
func Pointer2NullInt32(val *int32) sql.NullInt32 {
	if val == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Valid: true, Int32: *val}
}
func Pointer2NullString(val *string) sql.NullString {
	if val == nil {
		return sql.NullString{}
	}
	return sql.NullString{Valid: true, String: *val}
}
func Pointer2NullTime(val *time.Time) sql.NullTime {
	if val == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true, Time: *val}
}
func Pointer2NullFloat64(val *float64) sql.NullFloat64 {
	if val == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Valid: true, Float64: *val}
}

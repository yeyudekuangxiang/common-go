package converttool

import "time"

func PointerBool(val bool) *bool {
	return &val
}
func PointerInt64(val int64) *int64 {
	return &val
}
func PointerInt32(val int32) *int32 {
	return &val
}
func PointerInt(val int) *int {
	return &val
}
func PointerString(val string) *string {
	return &val
}
func PointerUint64(val uint64) *uint64 {
	return &val
}
func PointerUint(val uint) *uint {
	return &val
}
func PointerFloat64(val float64) *float64 {
	return &val
}
func PointerFloat32(val float32) *float32 {
	return &val
}
func PointerTime(t time.Time) *time.Time {
	return &t
}

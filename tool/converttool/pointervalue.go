package converttool

func PointerBool(val bool) *bool {
	return &val
}
func PointerInt64(val int64) *int64 {
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

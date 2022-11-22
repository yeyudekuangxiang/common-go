package converttool

import "strconv"

func MustInt(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}
func AbsInt(i int64) uint64 {
	if i > 0 {
		return uint64(i)
	}
	return uint64(-i)
}

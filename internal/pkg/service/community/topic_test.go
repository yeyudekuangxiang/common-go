package community

import (
	"fmt"
	"testing"
)

func TestParaseTime(t *testing.T) {
	var tp int
	var tp2 *int
	fmt.Println(tp)
	fmt.Println(tp2)
	tp2 = &tp
	fmt.Println(tp2)
	fmt.Println(*tp2)
}

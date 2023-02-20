package open

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t1 := time.UnixMilli(1676626163000)
	t2 := time.UnixMilli(1676340683000)

	fmt.Println(t1)
	fmt.Println(t2)

	diff := t1.Sub(t2).Hours()

	fmt.Println(diff)

	parseInt, err := strconv.ParseInt("1676859083000", 10, 64)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(parseInt)

	t3 := time.UnixMilli(parseInt)

	fmt.Println(t3)

	diff2 := t3.Sub(t1)

	fmt.Println(diff2)
	fmt.Println(diff2.Hours())

	if time.UnixMilli(parseInt).Sub(t1).Hours() > 72.0 {
		fmt.Println("true")
	}

}

package open

import (
	"fmt"
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
}

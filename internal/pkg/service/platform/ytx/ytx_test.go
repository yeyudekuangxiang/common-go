package ytx

import (
	"fmt"
	mioctx "mio/internal/pkg/core/context"
	"testing"
)

func TestGoFunc(t *testing.T) {
	var options []Options
	options = append(options, WithPoolCode("RP202110251300002"))
	options = append(options, WithSecret("a123456"))
	jhxService := NewYtxService(mioctx.NewMioContext(), options...)
	go func() {
		code, err := jhxService.SendCoupon()
		if err != nil {
			return
		}
		fmt.Println(code)
	}()
}

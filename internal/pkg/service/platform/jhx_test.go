package platform

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimal(t *testing.T) {
	amount := "3.14112312"
	fromString, err := decimal.NewFromString(amount)
	if err != nil {
		return
	}
	point := fromString.Mul(decimal.NewFromInt(10)).Round(2).String()
	fmt.Printf("point:%s\n", point)
}

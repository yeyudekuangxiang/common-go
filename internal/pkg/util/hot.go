package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

//(log10())
func hot(Qviews, Qanswers, Qscore, Ascores float64, dataAsk, dateActive time.Time) {
	Qage := time.Now().Sub(dataAsk).Seconds()
	Qage, _ = decimal.NewFromFloat(Qage / 3600).Round(1).Float64()

	Qupdate := time.Now().Sub(dateActive).Seconds()
	Qupdate, _ = decimal.NewFromFloat(Qupdate / 3600).Round(1).Float64()

	dividend := (math.Log10(Qviews) * 4) + ((Qanswers * Qscore) / 5) + Ascores
	divisor := math.Pow((Qage+1)-(Qage-Qupdate)/2, 1.5)
	res := dividend / divisor
	fmt.Println(res)
}

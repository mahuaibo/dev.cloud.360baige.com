package money

import (
	"fmt"
	//"math/big"
)

func CentToDollar(cent int64) string {
	dollar := fmt.Sprintf("%.2f", float64(cent)/100)
	return dollar
}

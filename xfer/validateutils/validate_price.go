package validateutils

import (
	"github.com/shopspring/decimal"
)

const MaxAmount = 1e6

func ValidPriceString(str string) (bool, string) {

	// don't try to replace the required validation
	if str == `` {
		return true, ``
	}

	d, err := decimal.NewFromString(str)
	if err != nil {
		return false, "That isn't a valid amount"
	}

	return ValidPrice(&d)
}

func ValidPrice(d *decimal.Decimal) (bool, string) {
	if d == nil {
		return false, "Invalid amount"
	}

	f, _ := d.Float64()
	if f < 0 {
		return false, "The amount cannot be negative"
	}
	if f > MaxAmount {
		return false, "The amount is too big"
	}

	return true, ""
}

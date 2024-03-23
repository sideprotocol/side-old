package types

import (
	"errors"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

// DecimalFractional constant for decimal calculations
var DecimalFractional = decimal.NewFromBigInt(big.NewInt(1_000_000_000_000_000_000), 0)

// SubSign returns mod subtraction and boolean indicating if the result is negative
func SubSign(a, b decimal.Decimal) (decimal.Decimal, bool) {
	if a.GreaterThanOrEqual(b) {
		return a.Sub(b), false
	}
	return b.Sub(a), true
}

// Sqrt computes the square root of a decimal number using Newton's method.
func Sqrt(value decimal.Decimal, precision decimal.Decimal) (decimal.Decimal, error) {
	if value.LessThan(decimal.Zero) {
		return decimal.Decimal{}, errors.New("square root of negative number")
	}

	x := value.DivRound(decimal.NewFromInt(2), precision.Exponent())
	lastX := decimal.Zero

	for x.Sub(lastX).Abs().GreaterThan(precision) {
		lastX = x
		x = decimal.Avg(x, value.Div(x))
	}

	return x, nil
}

// CalculatePow computes base^(exp) using an approximation algorithm
func ApproximatePow(baseS, expS string, precisionS string) (decimal.Decimal, error) {
	// Convert string to Decimal
	base, err := decimal.NewFromString(baseS)
	if err != nil {
		return decimal.Decimal{}, err
	}

	exp, err := decimal.NewFromString(expS)
	if err != nil {
		return decimal.Decimal{}, err
	}

	precision, err := decimal.NewFromString(precisionS)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if base.IsZero() && !exp.IsZero() {
		return base, nil
	}

	if base.GreaterThan(decimal.NewFromInt(2)) {
		return decimal.Decimal{}, errors.New("calculatePow: base must be less than 2")
	}

	integer := exp.Div(DecimalFractional).Floor()
	fractional := exp.Mod(DecimalFractional)
	integerPow := base.Pow(integer)

	if fractional.IsZero() {
		return integerPow, nil
	}

	fractionalPow, err := PowApprox(base, fractional, precision)
	if err != nil {
		return decimal.Decimal{}, err
	}

	result := integerPow.Mul(fractionalPow)
	return result, nil
}

// PowApprox approximates power for fractional exponents
func PowApprox(base, exp, precision decimal.Decimal) (decimal.Decimal, error) {
	if exp.Equals(decimal.NewFromInt(1).Div(decimal.NewFromInt(2))) {
		sqrtBase, err := Sqrt(base, precision)
		if err != nil {
			return decimal.Decimal{}, err
		}
		return sqrtBase, nil
	}

	x, xNeg := SubSign(base, decimal.NewFromInt(1))
	term := decimal.NewFromInt(1)
	sum := decimal.NewFromInt(1)
	negative := false

	a := exp
	bigK := decimal.Zero
	i := decimal.NewFromInt(1)

	for term.GreaterThanOrEqual(precision) {
		c, cNeg := SubSign(a, bigK)
		bigK = i

		newTerm := term.Mul(c).Mul(x).Div(bigK)
		term = newTerm

		if term.IsZero() {
			break
		}

		if xNeg {
			negative = !negative
		}

		if cNeg {
			negative = !negative
		}

		if negative {
			sum = sum.Sub(term)
		} else {
			sum = sum.Add(term)
		}

		i = i.Add(decimal.NewFromInt(1))
	}
	fixedPoint := DecimalPlacesFromPrecision(precision)
	return sum.RoundDown(fixedPoint), nil
}

// DecimalPlacesFromPrecision calculates the number of decimal places based on the precision decimal.
// DecimalPlacesFromPrecision calculates the number of decimal places based on the precision decimal.
func DecimalPlacesFromPrecision(precision decimal.Decimal) int32 {
	// Convert precision to string using no fixed decimal places to avoid trailing zeros.
	str := precision.String()
	dotIndex := strings.IndexRune(str, '.')
	if dotIndex != -1 {
		// Remove trailing zeros to get the correct count of significant decimal places.
		str = strings.TrimRight(str, "0")
		// Count the number of characters after the decimal point to determine the scale.
		return int32(len(str) - dotIndex - 1)
	}
	return 0 // Return 0 if there is no fractional part.
}

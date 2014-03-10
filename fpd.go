// Package implements a fixed-point decimal
package fpd

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

var Precision = 16 // precision during division when it doesn't divide exactly

// Decimal represents a fixed-point decimal.
type Decimal struct {
	value *big.Int
	scale int
}

// New returns a new fixed-point decimal
func New(value int64, scale int) Decimal {
	return Decimal{big.NewInt(value), scale}
}

func NewFromString(value string) (Decimal, error) {
	var intString string
	var scale int
	parts := strings.Split(value, ".")
	if len(parts) == 1 {
		// There is no decimal point, we can just parse the original string as
		// an int
		intString = value
		scale = 0
	} else if len(parts) == 2 {
		intString = parts[0] + parts[1]
		scale = -len(parts[1])
	} else {
		return Decimal{}, fmt.Errorf("can't convert %s to decimal", value)
	}

	dValue := big.NewInt(0)
	_, ok := dValue.SetString(intString, 10)
	if !ok {
		return Decimal{}, fmt.Errorf("can't convert %s to decimal", value)
	}

	return Decimal{dValue, scale}, nil
}

func NewFromFloat(value float64) Decimal {
	intPortion := math.Floor(math.Abs(value))

	var intDigits int
	if intPortion == 0 {
		intDigits = 0
	} else {
		intDigits = int(math.Log10(intPortion)) + 1
	}
	decDigits := 16 - intDigits // 16 is max significant digits in float64
	return NewFromFloatWithScale(value, -decDigits)
}

func NewFromFloatWithScale(value float64, scale int) Decimal {
	scaleMul := math.Pow(10, -float64(scale))
	intValue := int64(value * scaleMul)
	dValue := big.NewInt(intValue)

	return Decimal{dValue, scale}
}

func (d *Decimal) ensureInitialized() {
	if d.value == nil {
		d.value = big.NewInt(0)
	}
}

// Rescale returns a rescaled version of the decimal. Returned
// decimal may be less precise if the given scale is bigger
// than the initial scale of the Decimal
//
// Example:
//
// 	d := New(12345, -4)
//	d2 := d.rescale(-1)
//	d3 := d2.rescale(-4)
//	println(d1)
//	println(d2)
//	println(d3)
//
// Output:
//
//	1.2345
//	1.2
//	1.2000
//
func (d Decimal) rescale(scale int) Decimal {
	d.ensureInitialized()
	diff := int(math.Abs(float64(scale - d.scale)))
	value := big.NewInt(0).Set(d.value)
	ten := big.NewInt(10)

	for diff > 0 {
		if scale > d.scale {
			value = value.Quo(value, ten)
		} else if scale < d.scale {
			value = value.Mul(value, ten)
		}

		diff--
	}

	return Decimal{value, scale}
}

func (d Decimal) Abs() Decimal {
	d.ensureInitialized()
	d2Value := big.NewInt(0).Abs(d.value)
	return Decimal{d2Value, d.scale}
}

// Add adds d to d2 and return d3
func (d Decimal) Add(d2 Decimal) Decimal {
	baseScale := smallestOf(d.scale, d2.scale)
	rd := d.rescale(baseScale)
	rd2 := d2.rescale(baseScale)

	d.ensureInitialized()
	d3Value := big.NewInt(0).Add(rd.value, rd2.value)
	return Decimal{d3Value, baseScale}
}

// Sub subtracts d2 from d and returns d3
func (d Decimal) Sub(d2 Decimal) Decimal {
	baseScale := smallestOf(d.scale, d2.scale)
	rd := d.rescale(baseScale)
	rd2 := d2.rescale(baseScale)

	d3Value := big.NewInt(0).Sub(rd.value, rd2.value)
	return Decimal{d3Value, baseScale}
}

// Mul multiplies d with d2 and returns d3
func (d Decimal) Mul(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	d3Value := big.NewInt(0).Mul(d.value, d2.value)
	return Decimal{d3Value, d.scale + d2.scale}
}

// Mul divides d by d2 and returns d3
func (d Decimal) Div(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	negative := (d.value.Sign() > 0) != (d2.value.Sign() > 0)
	ad := d.Abs()
	ad2 := d2.Abs()

	shift := bigIntLen(ad2.value) - bigIntLen(ad.value) + Precision + 1

	scale := ad.scale - ad2.scale - shift

	coeff := big.NewInt(0)
	rem := big.NewInt(0)
	if shift >= 0 {
		shiftMult := big.NewInt(0).Exp(big.NewInt(10),
			big.NewInt(int64(shift)),
			nil)
		num := big.NewInt(0).Mul(ad.value, shiftMult)
		coeff, rem = coeff.DivMod(num, ad2.value, rem)
	} else {
		coeff, rem = coeff.DivMod(ad.value, ad2.value, rem)
		scale += shift
	}

	if rem.Cmp(big.NewInt(0)) == 0 {
		// result is exact, get as close to ideal exponent as possible
		ideal := ad.scale - ad2.scale
		modRes := big.NewInt(0)
		for scale < ideal &&
			modRes.Mod(coeff, big.NewInt(10)).Cmp(big.NewInt(0)) == 0 {
			coeff.Div(coeff, big.NewInt(10))
			scale += 1
		}
	}

	if negative {
		coeff.Neg(coeff)
	}

	return Decimal{coeff, scale}
}

// Cmp compares x and y and returns -1, 0 or 1
//
// Example
//
//-1 if x <  y
// 0 if x == y
//+1 if x >  y
//
func (d Decimal) Cmp(d2 Decimal) int {
	smallestScale := smallestOf(d.scale, d2.scale)
	rd := d.rescale(smallestScale)
	rd2 := d2.rescale(smallestScale)

	return rd.value.Cmp(rd2.value)
}

func (d Decimal) Scale() int {
	return d.scale
}

// String returns the string representatino of the decimal
// with the fixed point
//
// Example:
//
//     d := New(-12345, -3)
//     println(d.String())
//
// Output:
//
//     -12.345
//
func (d Decimal) String() string {
	if d.scale >= 0 {
		return d.rescale(0).value.String()
	}

	abs := big.NewInt(0).Abs(d.value)
	str := abs.String()

	var a, b string
	if len(str) > -d.scale {
		a = str[:len(str)+d.scale]
		b = str[len(str)+d.scale:]
	} else {
		a = "0"

		num0s := -d.scale - len(str)
		b = strings.Repeat("0", num0s) + str
	}

	if d.value.Sign() < 0 {
		return fmt.Sprintf("-%v.%v", a, b)
	}

	return fmt.Sprintf("%v.%v", a, b)
}

func (d Decimal) unformattedString() string {
	return d.value.String()
}

func smallestOf(x, y int) int {
	if x >= y {
		return y
	}
	return x
}

func bigIntLen(value *big.Int) int {
	// TODO: optimize
	ret := len(value.String())
	if value.Sign() < 0 {
		ret -= 1 // don't count - sign
	}
	return ret
}

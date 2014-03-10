package fpd

import "testing"

var testTable = map[float64]string{
	3.141592653589793:   "3.141592653589793",
	3:                   "3.000000000000000",
	1234567890123456:    "1234567890123456",
	1234567890123456000: "1234567890123456000",
	1234.567890123456:   "1234.567890123456",
	.1234567890123456:   "0.1234567890123456",
	0:                   "0.0000000000000000",
	.1111111111111110:   "0.1111111111111110",
	.1111111111111111:   "0.1111111111111111",
	.1111111111111119:   "0.1111111111111119",
	.0000000000000001:   "0.0000000000000001", // TODO: these should be able
	.0000000000000002:   "0.0000000000000002", // TODO: to have more than 16
	.0000000000000003:   "0.0000000000000003", // TODO: decimal digits
	.0000000000000005:   "0.0000000000000005",
	.0000000000000008:   "0.0000000000000008",
	.1000000000000001:   "0.1000000000000001",
	.1000000000000002:   "0.1000000000000002",
	.1000000000000003:   "0.1000000000000003",
	.1000000000000005:   "0.1000000000000005",
	.1000000000000008:   "0.1000000000000008",
}

func TestNewFromFloat(t *testing.T) {
	// add negatives
	for f, s := range testTable {
		if f > 0 {
			testTable[-f] = "-" + s
		}
	}

	for f, s := range testTable {
		d := NewFromFloat(f)
		if d.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, d.String(),
				d.unformattedString(), d.scale)
		}
	}
}

func TestNewFromString(t *testing.T) {
	// add negatives
	for f, s := range testTable {
		if f > 0 {
			testTable[-f] = "-" + s
		}
	}

	for _, s := range testTable {
		d, err := NewFromString(s)
		if err != nil {
			t.Errorf("error while parsing %s", s)
		} else if d.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, d.String(),
				d.unformattedString(), d.scale)
		}
	}
}

func TestNewFromStringErrs(t *testing.T) {
	tests := []string{
		"",
		"qwert",
		"-",
		".",
		"-.",
		".-",
		"234-.56",
		"234-56",
		"2-",
	}

	for _, s := range tests {
		_, err := NewFromString(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestNewFromFloatWithScale(t *testing.T) {
	type Inp struct {
		float float64
		scale int
	}
	tests := map[Inp]string{
		Inp{123.4, -3}:     "123.400",
		Inp{123.4, -1}:     "123.4",
		Inp{123.412345, 1}: "120",
		Inp{123.412345, 0}: "123",
		Inp{123.412345, -5}: "123.41234",
		Inp{123.412345, -6}: "123.412345",
		Inp{123.412345, -7}: "123.4123450",
	}

	// add negatives
	for p, s := range tests {
		if p.float > 0 {
			tests[Inp{-p.float, p.scale}] = "-" + s
		}
	}

	for input, s := range tests {
		d := NewFromFloatWithScale(input.float, input.scale)
		if d.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, d.String(),
				d.unformattedString(), d.scale)
		}
	}
}

func TestDecimal_rescale(t *testing.T) {
	type Inp struct {
		int     int64
		scale   int
		rescale int
	}
	tests := map[Inp]string{
		Inp{1234, -3, -5}: "1.23400",
		Inp{1234, -3, 0}:  "1",
		Inp{1234, 3, 0}:   "1234000",
		Inp{1234, -4, -4}: "0.1234",
	}

	// add negatives
	for p, s := range tests {
		if p.int > 0 {
			tests[Inp{-p.int, p.scale, p.rescale}] = "-" + s
		}
	}

	for input, s := range tests {
		d := New(input.int, input.scale).rescale(input.rescale)

		if d.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, d.String(),
				d.unformattedString(), d.scale)
		}
	}
}

func TestDecimal_Uninitialized(t *testing.T) {
	a := Decimal{}
	b := Decimal{}

	decs := []Decimal{
		a,
		a.rescale(10),
		a.Abs(),
		a.Add(b),
		a.Sub(b),
		a.Mul(b),
		a.Div(New(1, -1)),
	}

	for _, d := range decs {
		if d.String() != "0" {
			t.Errorf("expected 0, got %s", d.String())
		}
	}

	if a.Cmp(b) != 0 {
		t.Errorf("a != b")
	}
	if a.Scale() != 0 {
		t.Errorf("a.Scale() != 0")
	}
}

func TestDecimal_Add(t *testing.T) {
	type Inp struct {
		a string
		b string
	}

	adds := map[Inp]string {
		Inp{"2", "3"}: "5",
		Inp{"2454495034", "3451204593"}: "5905699627",
		Inp{"24544.95034", ".3451204593"}: "24545.2954604593",
		Inp{".1", ".1"}: "0.2",
		Inp{".1", "-.1"}: "0.0", // TODO: should this just be "0"?
		Inp{"0", "1.001"}: "1.001",
	}

	for inp, res := range adds {
		a, err := NewFromString(inp.a)
		if err != nil {
			t.FailNow()
		}
		b, err := NewFromString(inp.b)
		if err != nil {
			t.FailNow()
		}
		c := a.Add(b)
		if c.String() != res {
			t.Errorf("expected %s, got %s", res, c.String())
		}
	}
}

func TestDecimal_Sub(t *testing.T) {
	type Inp struct {
		a string
		b string
	}

	adds := map[Inp]string {
		Inp{"2", "3"}: "-1",
		Inp{"12", "3"}: "9",
		Inp{"-2", "9"}: "-11",
		Inp{"2454495034", "3451204593"}: "-996709559",
		Inp{"24544.95034", ".3451204593"}: "24544.6052195407",
		Inp{".1", "-.1"}: "0.2",
		Inp{".1", ".1"}: "0.0", // TODO: should this just be "0"?
		Inp{"0", "1.001"}: "-1.001",
		Inp{"1.001", "0"}: "1.001",
		Inp{"2.3", ".3"}: "2.0",
	}

	for inp, res := range adds {
		a, err := NewFromString(inp.a)
		if err != nil {
			t.FailNow()
		}
		b, err := NewFromString(inp.b)
		if err != nil {
			t.FailNow()
		}
		c := a.Sub(b)
		if c.String() != res {
			t.Errorf("expected %s, got %s", res, c.String())
		}
	}
}

func TestDecimal_Mul(t *testing.T) {
	type Inp struct {
		a string
		b string
	}

	mults := map[Inp]string {
		Inp{"2", "3"}: "6",
		Inp{"2454495034", "3451204593"}: "8470964534836491162",
		Inp{"24544.95034", ".3451204593"}: "8470.964534836491162",
		Inp{".1", ".1"}: "0.01",
		Inp{"0", "1.001"}: "0.000", // TODO: should this just be "0"?
	}

	for inp, res := range mults {
		a, err := NewFromString(inp.a)
		if err != nil {
			t.FailNow()
		}
		b, err := NewFromString(inp.b)
		if err != nil {
			t.FailNow()
		}
		c := a.Mul(b)
		if c.String() != res {
			t.Errorf("expected %s, got %s", res, c.String())
		}
	}

	// positive scale
	c := New(1234,5).Mul(New(45,-1))
	if c.String() != "555300000" {
		t.Errorf("Expected %s, got %s", "555300000", c.String())
	}
}

// old tests after this line

func TestDecimal_Scale(t *testing.T) {
	a := New(1234, -3)
	if a.Scale() != -3 {
		t.Errorf("error")
	}
}

func TestDecimal_Abs1(t *testing.T) {
	a := New(-1234, -4)
	b := New(1234, -4)

	c := a.Abs()
	if c.Cmp(b) != 0 {
		t.Errorf("error")
	}
}

func TestDecimal_Abs2(t *testing.T) {
	a := New(-1234, -4)
	b := New(1234, -4)

	c := b.Abs()
	if c.Cmp(a) == 0 {
		t.Errorf("error")
	}
}

func TestDecimal_Div1(t *testing.T) {
	a := New(1398699, -4)
	b := New(1006, -3)

	c := a.Div(b)
	if c.unformattedString() != "1390356" {
		t.Errorf(c.unformattedString())
	}
}

func TestDecimal_Div2(t *testing.T) {
	a := New(2345, -3)
	b := New(2, 0)

	c := a.Div(b)
	if c.unformattedString() != "1172" {
		t.Errorf(c.unformattedString())
	}
}

func TestDecimal_Div3(t *testing.T) {
	a := New(18275499, -6)
	b := New(16275499, -6)

	c := a.Div(b)
	if c.unformattedString() != "1122884" {
		t.Errorf(c.unformattedString())
	}
}
func TestDecimal_Cmp1(t *testing.T) {
	a := New(123, 3)
	b := New(-1234, 2)

	if a.Cmp(b) != 1 {
		t.Errorf("Error")
	}
}

func TestDecimal_Cmp2(t *testing.T) {
	a := New(123, 3)
	b := New(1234, 2)

	if a.Cmp(b) != -1 {
		t.Errorf("Error")
	}
}

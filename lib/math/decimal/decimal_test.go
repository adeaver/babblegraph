package decimal_test

import (
	"babblegraph/lib/math/decimal"
	"testing"
)

func TestAddInt64(t *testing.T) {
	d1 := decimal.FromInt64(5)
	d2 := decimal.FromInt64(2)
	d3 := decimal.FromInt64(0)
	d4 := decimal.FromInt64(-5)

	d5 := d1.Add(d2)
	d6 := d1.Add(d3)
	d7 := d1.Add(d4)

	expectedD5 := decimal.FromInt64(7)
	expectedD6 := decimal.FromInt64(5)
	expectedD7 := decimal.FromInt64(0)
	if !d5.EqualTo(expectedD5) {
		t.Errorf("D5: expected %+v, Got %+v", expectedD5, d5)
	}
	if !d6.EqualTo(expectedD6) {
		t.Errorf("D6: expected %+v, Got %+v", expectedD6, d6)
	}
	if !d7.EqualTo(expectedD7) {
		t.Errorf("D7: expected %+v, Got %+v", expectedD7, d7)
	}
}

func TestMultiplyInt64(t *testing.T) {
	d1 := decimal.FromInt64(5)
	d2 := decimal.FromInt64(2)
	d3 := decimal.FromInt64(0)
	d4 := decimal.FromInt64(-5)

	d5 := d1.Multiply(d2)
	d6 := d1.Multiply(d3)
	d7 := d1.Multiply(d4)

	expectedD5 := decimal.FromInt64(10)
	expectedD6 := decimal.FromInt64(0)
	expectedD7 := decimal.FromInt64(-25)
	if !d5.EqualTo(expectedD5) {
		t.Errorf("D5: expected %+v, Got %+v", expectedD5, d5)
	}
	if !d6.EqualTo(expectedD6) {
		t.Errorf("D6: expected %+v, Got %+v", expectedD6, d6)
	}
	if !d7.EqualTo(expectedD7) {
		t.Errorf("D7: expected %+v, Got %+v", expectedD7, d7)
	}
}

func TestDivide(t *testing.T) {
	d1 := decimal.FromInt64(10)
	d2 := decimal.FromInt64(5)
	d3 := decimal.FromInt64(1)
	d4 := decimal.FromInt64(4)

	d5 := d1.Divide(d2)
	d6 := d1.Divide(d3)
	d7 := d1.Divide(d4)

	expectedD5 := decimal.FromInt64(2)
	expectedD6 := decimal.FromInt64(10)
	expectedD7 := decimal.FromFloat64(2.5)
	if !d5.EqualTo(expectedD5) {
		t.Errorf("D5: expected %+v, Got %+v", expectedD5, d5)
	}
	if !d6.EqualTo(expectedD6) {
		t.Errorf("D6: expected %+v, Got %+v", expectedD6, d6)
	}
	if !d7.EqualTo(expectedD7) {
		t.Errorf("D7: expected %+v, Got %+v", expectedD7, d7)
	}
}

func TestFloat64(t *testing.T) {
	d1 := decimal.FromFloat64(5.5)
	d2 := decimal.FromFloat64(2.5)
	d3 := decimal.FromFloat64(.0001)

	d4 := d1.Multiply(d3)
	d5 := d1.Add(d2)
	d6 := d1.Subtract(d2)
	d7 := d1.Multiply(d2)
	d8 := d1.Divide(d2)

	expectedD4 := decimal.FromFloat64(0.00055)
	expectedD5 := decimal.FromFloat64(8.0)
	expectedD6 := decimal.FromFloat64(3.0)
	expectedD7 := decimal.FromFloat64(13.75)
	expectedD8 := decimal.FromFloat64(2.2)
	if !d4.EqualTo(expectedD4) {
		t.Errorf("D4: expected %+v, Got %+v", expectedD4, d4)
	}
	if !d5.EqualTo(expectedD5) {
		t.Errorf("D5: expected %+v, Got %+v", expectedD5, d5)
	}
	if !d6.EqualTo(expectedD6) {
		t.Errorf("D6: expected %+v, Got %+v", expectedD6, d6)
	}
	if !d7.EqualTo(expectedD7) {
		t.Errorf("D7: expected %+v, Got %+v", expectedD7, d7)
	}
	if !d8.EqualTo(expectedD8) {
		t.Errorf("D8: expected %+v, Got %+v", expectedD8, d8)
	}
}

package decimal

import (
	"encoding/json"
	"math"
)

const percision = 1000000

type Number struct {
	millionths int64
}

func (n Number) Ptr() *Number {
	return &n
}

func FromInt64(i int64) Number {
	return Number{
		millionths: i * percision,
	}
}

func FromFloat64(f float64) Number {
	return Number{
		millionths: int64(float64(percision) * f),
	}
}

func (d Number) Multiply(n Number) Number {
	return Number{
		millionths: (d.millionths * n.millionths) / percision,
	}
}

func (d Number) Divide(n Number) Number {
	divisor := float64(d.millionths) / float64(n.millionths)
	return Number{
		millionths: int64(divisor * float64(percision)),
	}
}

func (d Number) Add(n Number) Number {
	return Number{
		millionths: (d.millionths + n.millionths),
	}
}

func (d Number) Subtract(n Number) Number {
	return Number{
		millionths: (d.millionths - n.millionths),
	}
}

func (d Number) Log10() Number {
	return FromFloat64(math.Log10(float64(d.millionths)) - math.Log10(percision))
}

func (d Number) GreaterThan(n Number) bool {
	return d.millionths > n.millionths
}

func (d Number) GreaterThanOrEqualTo(n Number) bool {
	return d.millionths >= n.millionths
}

func (d Number) LessThan(n Number) bool {
	return d.millionths < n.millionths
}

func (d Number) LessThanOrEqualTo(n Number) bool {
	return d.millionths <= n.millionths
}

func (d Number) EqualTo(n Number) bool {
	return d.millionths == n.millionths
}

func (d Number) ToFloat64() float64 {
	return float64(d.millionths) / float64(percision)
}

func (d Number) ToInt64Rounded() int64 {
	return int64(math.Round(d.ToFloat64()))
}

func (d Number) ToInt64Truncated() int64 {
	return int64(d.ToFloat64())
}

type jsonNumber struct {
	Millionths int64 `json:"millionths"`
}

func (d Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonNumber{Millionths: d.millionths})
}

func (d *Number) UnmarshalJSON(data []byte) error {
	var n jsonNumber
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	d.millionths = n.Millionths
	return nil
}

package decimal

import (
	"fmt"
	"log"
	"math"
)

const percision = 10000

type Number struct {
	tenthousandths int64
}

func FromInt64(i int64) Number {
	return Number{
		tenthousandths: i * percision,
	}
}

func FromFloat64(f float64) Number {
	return Number{
		tenthousandths: int64(float64(percision) * f),
	}
}

func (d Number) Multiply(n Number) Number {
	return Number{
		tenthousandths: (d.tenthousandths * n.tenthousandths) / percision,
	}
}

func (d Number) Divide(n Number) Number {
	divisor := float64(d.tenthousandths) / float64(n.tenthousandths)
	log.Println(fmt.Sprintf("top %d, bottom %d, Val %f", d.tenthousandths, n.tenthousandths, divisor))
	return Number{
		tenthousandths: int64(divisor * float64(percision)),
	}
}

func (d Number) Add(n Number) Number {
	return Number{
		tenthousandths: (d.tenthousandths + n.tenthousandths),
	}
}

func (d Number) Subtract(n Number) Number {
	return Number{
		tenthousandths: (d.tenthousandths - n.tenthousandths),
	}
}

func (d Number) Log10() Number {
	return FromFloat64(math.Log10(float64(d.tenthousandths)) - math.Log10(percision))
}

func (d Number) GreaterThan(n Number) bool {
	return d.tenthousandths > n.tenthousandths
}

func (d Number) GreaterThanOrEqualTo(n Number) bool {
	return d.tenthousandths >= n.tenthousandths
}

func (d Number) LessThan(n Number) bool {
	return d.tenthousandths < n.tenthousandths
}

func (d Number) LessThanOrEqualTo(n Number) bool {
	return d.tenthousandths <= n.tenthousandths
}

func (d Number) EqualTo(n Number) bool {
	return d.tenthousandths == n.tenthousandths
}

func (d Number) ToFloat64() float64 {
	return float64(d.tenthousandths) / float64(percision)
}

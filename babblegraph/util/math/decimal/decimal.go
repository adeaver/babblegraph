package decimal

const percision = 1000000

type Number struct {
	millionths int64
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
	return Number{
		millionths: int64((float64(d.millionths) / float64(n.millionths)) * float64(percision)),
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
	return float64(d.millionths) * float64(percision)
}

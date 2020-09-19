package decimal

const percision = 1000000

type DecimalNumber struct {
	millionths int64
}

func FromInt64(i int64) DecimalNumber {
	return DecimalNumber{
		millionths: i * percision,
	}
}

func FromFloat64(f float64) DecimalNumber {
	return DecimalNumber{
		millionths: int64(float64(percision) * f),
	}
}

func (d DecimalNumber) Multiply(n DecimalNumber) DecimalNumber {
	return DecimalNumber{
		millionths: (d.millionths * n.millionths) / percision,
	}
}

func (d DecimalNumber) Divide(n DecimalNumber) DecimalNumber {
	return DecimalNumber{
		millionths: int64((float64(d.millionths) / float64(n.millionths)) * float64(percision)),
	}
}

func (d DecimalNumber) Add(n DecimalNumber) DecimalNumber {
	return DecimalNumber{
		millionths: (d.millionths + n.millionths),
	}
}

func (d DecimalNumber) Subtract(n DecimalNumber) DecimalNumber {
	return DecimalNumber{
		millionths: (d.millionths - n.millionths),
	}
}

func (d DecimalNumber) GreaterThan(n DecimalNumber) bool {
	return d.millionths > n.millionths
}

func (d DecimalNumber) GreaterThanOrEqualTo(n DecimalNumber) bool {
	return d.millionths >= n.millionths
}

func (d DecimalNumber) LessThan(n DecimalNumber) bool {
	return d.millionths < n.millionths
}

func (d DecimalNumber) LessThanOrEqualTo(n DecimalNumber) bool {
	return d.millionths <= n.millionths
}

func (d DecimalNumber) EqualTo(n DecimalNumber) bool {
	return d.millionths == n.millionths
}

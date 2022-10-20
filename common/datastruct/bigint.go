package datastruct

import "math/big"

type BigInt struct {
	v        *big.Int
	readonly bool
}

func NewBigInt(v int64) *BigInt {
	return &BigInt{
		v:        big.NewInt(v),
		readonly: false,
	}
}

func NewReadonlyBigInt(v int64) *BigInt {
	return &BigInt{
		v:        big.NewInt(v),
		readonly: true,
	}
}

// copy will set readonly to false
func (z *BigInt) Copy() *BigInt {
	return &BigInt{
		v:        big.NewInt(0).Set(z.v),
		readonly: false,
	}
}

// Add sets z to the sum x+y and returns z.
func (z *BigInt) Add(x, y *BigInt) *BigInt {
	checkReadonly(z)
	z.v.Add(x.v, y.v)
	return z
}

// Sub sets z to the difference x-y and returns z.
func (z *BigInt) Sub(x, y *BigInt) *BigInt {
	checkReadonly(z)
	z.v.Sub(x.v, y.v)
	return z
}

// Mul sets z to the product x*y and returns z.
func (z *BigInt) Mul(x, y *BigInt) *BigInt {
	checkReadonly(z)
	z.v.Mul(x.v, y.v)
	return z
}

// Div sets z to the quotient x/y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Div implements Euclidean division (unlike Go); see DivMod for more details.
func (z *BigInt) Div(x, y *BigInt) *BigInt {
	checkReadonly(z)
	z.v.Div(x.v, y.v)
	return z
}

func checkReadonly(v *BigInt) {
	if v.readonly == true {
		panic("invade bigint readonly part")
	}
}

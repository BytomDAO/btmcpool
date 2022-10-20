package datastruct

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	x := NewReadonlyBigInt(100)
	y := x.Copy()
	y.v = big.NewInt(32)

	assert.Equal(t, big.NewInt(100), x.v)
	assert.True(t, x.readonly)

	assert.Equal(t, big.NewInt(32), y.v)
	assert.False(t, y.readonly)
}

func TestInvadeOP(t *testing.T) {
	x := NewBigInt(200)
	y := NewBigInt(100)
	z := NewReadonlyBigInt(0)

	assert.Panics(t, func() {
		z.Add(x, y)
	})
	assert.Panics(t, func() {
		z.Sub(x, y)
	})
	assert.Panics(t, func() {
		z.Mul(x, y)
	})
	assert.Panics(t, func() {
		z.Div(x, y)
	})
}

func TestNormalOP(t *testing.T) {
	x := NewBigInt(200)
	y := NewBigInt(100)
	z := NewBigInt(0)

	z.Add(x, y)
	assert.Equal(t, z.v, big.NewInt(300))
	assert.Equal(t, z.readonly, false)

	z.Sub(x, y)
	assert.Equal(t, z.v, big.NewInt(100))
	assert.Equal(t, z.readonly, false)

	z.Mul(x, y)
	assert.Equal(t, z.v, big.NewInt(20000))
	assert.Equal(t, z.readonly, false)

	z.Div(x, y)
	assert.Equal(t, z.v, big.NewInt(2))
	assert.Equal(t, z.readonly, false)
}

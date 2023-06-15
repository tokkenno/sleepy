package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const (
	lessThan    = -1
	equal       = 0
	greaterThan = 1
)

type UInt128 interface {
	Clone() UInt128
	Equal(o UInt128) bool
	Compare(o UInt128) int

	GetBit(position int) uint8

	ToInt() (int, error)
	ToUInt64() (uint64, uint64)
	ToBytes() []byte
	ToBitString() string
	ToHexString() string

	Add(o UInt128)
	Sub(num2 UInt128)
	Mul(num2 UInt128)
	DivMod(denominator UInt128) (mod UInt128, err error)

	Not()
	And(o UInt128)
	Or(o UInt128)
	Xor(o UInt128)

	LeftShift(numBits uint)
	RightShift(numBits uint)
}

type uint128Imp struct {
	lo uint64
	hi uint64
}

var _ UInt128 = (*uint128Imp)(nil)

func NewUInt128(low uint64, high uint64) UInt128 {
	return &uint128Imp{lo: low, hi: high}
}

func NewUInt128FromInt(val int) UInt128 {
	if val < 0 {
		return NewUInt128(uint64(^uint64(0)-uint64(-val-1)), ^uint64(0))
	} else {
		return NewUInt128(uint64(val), 0)
	}
}

func NewUInt128FromLong(val int64) UInt128 {
	num := new(uint128Imp)
	num.lo = uint64(val)
	return num
}

func NewUInt128FromByteArray(data []byte) (UInt128, error) {
	if len(data) > 16 {
		return nil, errors.New("uint128Imp only can be parsed from 16 byte length array")
	}

	if len(data) < 16 {
		data = append(make([]byte, 16-len(data)), data...)
	}

	num := new(uint128Imp)
	num.hi = binary.BigEndian.Uint64(data[0:8])
	num.lo = binary.BigEndian.Uint64(data[8:16])
	return num, nil
}

func (uint128 *uint128Imp) Clone() UInt128 {
	return &uint128Imp{
		uint128.lo,
		uint128.hi,
	}
}

func (uint128 *uint128Imp) Compare(o UInt128) int {
	lo, hi := o.ToUInt64()
	if uint128.hi < hi {
		return lessThan
	} else if uint128.hi > hi {
		return greaterThan
	} else if uint128.lo < lo {
		return lessThan
	} else if uint128.lo > lo {
		return greaterThan
	}

	return equal
}

func (uint128 *uint128Imp) Equal(o UInt128) bool {
	lo, hi := o.ToUInt64()
	return uint128.hi == hi && uint128.lo == lo
}

func (uint128 *uint128Imp) GetBit(position int) uint8 {
	valT := uint128.hi
	if position > 63 {
		valT = uint128.lo
		position -= 64
	}
	iShift := uint(63 - (position % 64))
	return uint8((valT >> iShift) & 1)
}

// ToUInt64 returns the value of the UInt128 as two uint64 values, the first being the lower 64 bits, the second the higher 64 bits.
func (uint128 *uint128Imp) ToUInt64() (uint64, uint64) {
	return uint128.lo, uint128.hi
}

func And(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.And(val2)
	return newValue
}

// Calculate the binary operation AND between the actual number and other
func (uint128 *uint128Imp) And(o UInt128) {
	lo, hi := o.ToUInt64()
	uint128.hi &= hi
	uint128.lo &= lo
}

// Calculate the binary operation OR between two uint128Imp
func Or(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Or(val2)
	return newValue
}

// Calculate the binary operation OR between the actual number and other
func (uint128 *uint128Imp) Or(o UInt128) {
	lo, hi := o.ToUInt64()
	uint128.hi |= hi
	uint128.lo |= lo
}

// Calculate the binary operation XOR between two uint128Imp
func Xor(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Xor(val2)
	return newValue
}

// Calculate the binary operation XOR between the actual number and other
func (uint128 *uint128Imp) Xor(o UInt128) {
	lo, hi := o.ToUInt64()
	uint128.hi ^= hi
	uint128.lo ^= lo
}

// Calculate the binary operation NOT for an uint128Imp number
func Not(val1 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Not()
	return newValue
}

func (uint128 *uint128Imp) Not() {
	uint128.lo = ^uint128.lo
	uint128.hi = ^uint128.hi
}

func Add(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Add(val2)
	return newValue
}

// Calculate the arithmetic operation ADD between the actual number and other
func (uint128 *uint128Imp) Add(o UInt128) {
	carry := uint128.lo
	lo, hi := o.ToUInt64()
	uint128.lo += lo
	uint128.hi += hi

	if uint128.lo < carry {
		uint128.hi += 1
	}
}

// Calculate the arithmetic operation SUB (subtract) between two uint128Imp
func Sub(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Sub(val2)
	return newValue
}

// Calculate the arithmetic operation SUB (subtract) between the actual number and other
func (uint128 *uint128Imp) Sub(num2 UInt128) {
	var AI, BI, R big.Int
	AI.SetBytes(uint128.ToBytes())
	BI.SetBytes(num2.ToBytes())
	R.Sub(&AI, &BI)

	tValue, _ := NewUInt128FromByteArray(R.Bytes())
	lo, hi := tValue.ToUInt64()
	uint128.hi = hi
	uint128.lo = lo
}

// Calculate the arithmetic operation MUL (product) between two uint128Imp
func Mul(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Mul(val2)
	return newValue
}

// Calculate the arithmetic operation MUL (product) between the actual number and other
func (uint128 *uint128Imp) Mul(num2 UInt128) {
	var AI, BI, R big.Int
	AI.SetBytes(uint128.ToBytes())
	BI.SetBytes(num2.ToBytes())
	R.Mul(&AI, &BI)

	tValue, _ := NewUInt128FromByteArray(R.Bytes())
	lo, hi := tValue.ToUInt64()
	uint128.hi = hi
	uint128.lo = lo
}

// Calculate the arithmetic operation DIVMOD (division & modulus) between two uint128Imp
func DivMod(val1 *uint128Imp, val2 *uint128Imp) (newValue UInt128, mod UInt128, err error) {
	newValue = val1.Clone()
	mod, err = newValue.DivMod(val2)
	return
}

// Calculate the arithmetic operation DIVMOD (division & modulus) between the actual number and other
func (uint128 *uint128Imp) DivMod(denominator UInt128) (mod UInt128, err error) {
	if denominator.Equal(NewUInt128FromInt(0)) {
		mod, err = nil, errors.New("divide by 0")
	} else {
		var AI, BI, R, Q big.Int
		AI.SetBytes(uint128.ToBytes())
		BI.SetBytes(denominator.ToBytes())
		R.DivMod(&AI, &BI, &Q)

		tValue, _ := NewUInt128FromByteArray(R.Bytes())
		lo, hi := tValue.ToUInt64()
		uint128.hi = hi
		uint128.lo = lo
		mod, err = NewUInt128FromByteArray(Q.Bytes())
	}
	return
}

func (uint128 *uint128Imp) LeftShift(numBits uint) {
	if numBits == 0 {
		return
	} else if numBits > 128 {
		uint128.hi = 0
		uint128.lo = 0
	} else if numBits > 64 {
		uint128.hi = uint128.lo << (numBits - 64)
		uint128.lo = 0
	} else {
		tmp := uint128.lo >> (64 - numBits)
		uint128.lo = uint128.lo << numBits
		uint128.hi = (uint128.hi << numBits) | tmp
	}
}

func (uint128 *uint128Imp) RightShift(numBits uint) {
	if numBits == 0 {
		return
	} else if numBits > 128 {
		uint128.hi = 0
		uint128.lo = 0
	} else if numBits > 64 {
		uint128.lo = uint128.hi >> (numBits - 64)
		uint128.hi = 0
	} else {
		tmp := uint128.hi << (64 - numBits)
		uint128.hi = uint128.hi >> numBits
		uint128.lo = (uint128.lo >> numBits) | tmp
	}
}

func (uint128 *uint128Imp) ToInt() (int, error) {
	if uint128.hi != 0 || uint128.lo > math.MaxUint32 {
		return 0, errors.New("the value inside uint128 can't be represented by an int32")
	} else {
		return int(uint32(uint128.lo)), nil
	}
}

func (uint128 *uint128Imp) ToBytes() []byte {
	hi, lo := make([]byte, 8), make([]byte, 8)
	binary.BigEndian.PutUint64(hi, uint128.hi)
	binary.BigEndian.PutUint64(lo, uint128.lo)
	return append(hi, lo...)
}

func (uint128 *uint128Imp) ToHexString() string {
	loStr := strconv.FormatUint(uint128.lo, 16)
	hiStr := strconv.FormatUint(uint128.hi, 16)

	var loFill []rune
	if len(loStr) < 16 {
		loFill = make([]rune, 16-len(loStr))
	}
	for i := range loFill {
		loFill[i] = '\u0030'
	}

	var hiFill []rune
	if len(hiFill) < 16 {
		hiFill = make([]rune, 16-len(hiStr))
	}
	for i := range hiFill {
		hiFill[i] = '\u0030'
	}

	return string(hiFill) + hiStr + string(loFill) + loStr
}

func (uint128 *uint128Imp) ToBitString() string {
	loStr := fmt.Sprintf("%b", uint128.lo)
	hiStr := fmt.Sprintf("%b", uint128.hi)

	var loFill []rune
	if len(loStr) < 64 {
		loFill = make([]rune, 64-len(loStr))
	}
	for i := range loFill {
		loFill[i] = '\u0030'
	}

	var hiFill []rune
	if len(hiFill) < 64 {
		hiFill = make([]rune, 64-len(hiStr))
	}
	for i := range hiFill {
		hiFill[i] = '\u0030'
	}

	return string(hiFill) + hiStr + string(loFill) + loStr
}

func (uint128 *uint128Imp) ToBaseString(base int) (string, error) {
	if uint128.lo == 0 && uint128.hi == 0 {
		return "0", nil
	}

	if base < 2 || base > 10 {
		return "", errors.New("invalid base")
	}

	numStr := string("")

	zero := NewUInt128FromInt(0)
	radix := NewUInt128FromInt(base)
	ii := uint128.Clone()

	for !ii.Equal(zero) {
		tmpRemainder, _ := ii.DivMod(radix)
		remainder, _ := tmpRemainder.ToInt()
		character := "0123456789"[remainder]
		numStr = string(character) + numStr
	}

	return numStr, nil
}

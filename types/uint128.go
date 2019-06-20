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

type UInt128 struct {
	lo uint64
	hi uint64
}

func NewUInt128(low uint64, high uint64) *UInt128 {
	return &UInt128{lo: low, hi: high}
}

func NewUInt128FromInt(val int) *UInt128 {
	if val < 0 {
		return NewUInt128(uint64(^uint64(0)-uint64(-val-1)), ^uint64(0))
	} else {
		return NewUInt128(uint64(val), 0)
	}
}

func NewUInt128FromLong(val int64) *UInt128 {
	num := new(UInt128)
	num.lo = uint64(val)
	return num
}

func NewUInt128FromByteArray(data []byte) (*UInt128, error) {
	if len(data) > 16 {
		return nil, errors.New("UInt128 only can be parsed from 16 byte length array")
	}

	if len(data) < 16 {
		data = append(make([]byte, 16-len(data)), data...)
	}

	num := new(UInt128)
	num.hi = binary.BigEndian.Uint64(data[0:8])
	num.lo = binary.BigEndian.Uint64(data[8:16])
	return num, nil
}

// Create a instance copy
func (uint128 *UInt128) Clone() *UInt128 {
	return &UInt128{
		uint128.lo,
		uint128.hi,
	}
}

// Compare if two UInt128 has equals, greater or less in value
func (uint128 *UInt128) Compare(o *UInt128) int {
	if uint128.hi < o.hi {
		return lessThan
	} else if uint128.hi > o.hi {
		return greaterThan
	} else if uint128.lo < o.lo {
		return lessThan
	} else if uint128.lo > o.lo {
		return greaterThan
	}

	return equal
}

// Check if two UInt128 has same value
func (uint128 *UInt128) Equal(o *UInt128) bool {
	return uint128.hi == o.hi && uint128.lo == o.lo
}

// Get the bit value in the provided position, from left to right
func (uint128 *UInt128) GetBit(position int) uint8 {
	valT := uint128.hi
	if position > 63 {
		valT = uint128.lo
		position -= 64
	}
	iShift := uint(63 - (position % 64))
	return uint8((valT >> iShift) & 1)
}

// Calculate the binary operation AND between two UInt128
func And(val1 *UInt128, val2 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.And(val2)
	return newValue
}

// Calculate the binary operation AND between the actual number and other
func (uint128 *UInt128) And(o *UInt128) {
	uint128.hi &= o.hi
	uint128.lo &= o.lo
}

// Calculate the binary operation OR between two UInt128
func Or(val1 *UInt128, val2 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.Or(val2)
	return newValue
}

// Calculate the binary operation OR between the actual number and other
func (uint128 *UInt128) Or(o *UInt128) {
	uint128.hi |= o.hi
	uint128.lo |= o.lo
}

// Calculate the binary operation XOR between two UInt128
func Xor(val1 *UInt128, val2 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.Xor(val2)
	return newValue
}

// Calculate the binary operation XOR between the actual number and other
func (uint128 *UInt128) Xor(o *UInt128) {
	uint128.hi ^= o.hi
	uint128.lo ^= o.lo
}

// Calculate the binary operation NOT for an UInt128 number
func Not(val1 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.Not()
	return newValue
}

// Calculate the binary operation NOT for the actual number
func (uint128 *UInt128) Not() {
	uint128.lo = ^uint128.lo
	uint128.hi = ^uint128.hi
}

// Calculate the arithmetic operation ADD between two UInt128
func Add(val1 *UInt128, val2 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.Add(val2)
	return newValue
}

// Calculate the arithmetic operation ADD between the actual number and other
func (uint128 *UInt128) Add(o *UInt128) {
	carry := uint128.lo
	uint128.lo += o.lo
	uint128.hi += o.hi

	if uint128.lo < carry {
		uint128.hi += 1
	}
}

// Calculate the arithmetic operation SUB (subtract) between two UInt128
func Sub(val1 *UInt128, val2 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.Sub(val2)
	return newValue
}

// Calculate the arithmetic operation SUB (subtract) between the actual number and other
func (uint128 *UInt128) Sub(num2 *UInt128) {
	var AI, BI, R big.Int
	AI.SetBytes(uint128.ToBytes())
	BI.SetBytes(num2.ToBytes())
	R.Sub(&AI, &BI)

	tValue, _ := NewUInt128FromByteArray(R.Bytes())
	uint128.hi = tValue.hi
	uint128.lo = tValue.lo
}

// Calculate the arithmetic operation MUL (product) between two UInt128
func Mul(val1 *UInt128, val2 *UInt128) *UInt128 {
	newValue := val1.Clone()
	newValue.Mul(val2)
	return newValue
}

// Calculate the arithmetic operation MUL (product) between the actual number and other
func (uint128 *UInt128) Mul(num2 *UInt128) {
	var AI, BI, R big.Int
	AI.SetBytes(uint128.ToBytes())
	BI.SetBytes(num2.ToBytes())
	R.Mul(&AI, &BI)

	tValue, _ := NewUInt128FromByteArray(R.Bytes())
	uint128.hi = tValue.hi
	uint128.lo = tValue.lo
}

// Calculate the arithmetic operation DIVMOD (division & modulus) between two UInt128
func DivMod(val1 *UInt128, val2 *UInt128) (newValue *UInt128, mod *UInt128, err error) {
	newValue = val1.Clone()
	mod, err = newValue.DivMod(val2)
	return
}

// Calculate the arithmetic operation DIVMOD (division & modulus) between the actual number and other
func (uint128 *UInt128) DivMod(denominator *UInt128) (mod *UInt128, err error) {
	if denominator.Equal(NewUInt128FromInt(0)) {
		mod, err = nil, errors.New("divide by 0")
	} else {
		var AI, BI, R, Q big.Int
		AI.SetBytes(uint128.ToBytes())
		BI.SetBytes(denominator.ToBytes())
		R.DivMod(&AI, &BI, &Q)

		tValue, _ := NewUInt128FromByteArray(R.Bytes())
		uint128.hi = tValue.hi
		uint128.lo = tValue.lo
		mod, err = NewUInt128FromByteArray(Q.Bytes())
	}
	return
}

// Move [numBits] number of bits to the left
func (uint128 *UInt128) LeftShift(numBits uint) {
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

// Move [numBits] number of bits to the right
func (uint128 *UInt128) RightShift(numBits uint) {
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

// Get the int representation, if the value take inside it
func (uint128 *UInt128) ToInt() (int, error) {
	if uint128.hi != 0 || uint128.lo > math.MaxUint32 {
		return 0, errors.New("the value inside uint128 can't be represented by an int32")
	} else {
		return int(uint32(uint128.lo)), nil
	}
}

// Get the high and low uint64 that forms the number
func (uint128 *UInt128) ToUInt64() (low uint64, high uint64) {
	low = uint128.lo
	high = uint128.hi
	return
}

// Get the 16 bytes of the number
func (uint128 *UInt128) ToBytes() []byte {
	hi, lo := make([]byte, 8), make([]byte, 8)
	binary.BigEndian.PutUint64(hi, uint128.hi)
	binary.BigEndian.PutUint64(lo, uint128.lo)
	return append(hi, lo...)
}

// Convert to the hexadecimal string representation
func (uint128 *UInt128) ToHexString() string {
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

// Convert to the bit string representation
func (uint128 *UInt128) ToBitString() string {
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

// Convert to string in an arbitrary base between 2 and 10
func (uint128 *UInt128) ToBaseString(base int) (string, error) {
	if uint128.lo == 0 && uint128.hi == 0 {
		return "0", nil;
	}

	if base < 2 || base > 10 {
		return "", errors.New("invalid base");
	}

	numStr := string("")

	zero := NewUInt128FromInt(0)
	radix := NewUInt128FromInt(base)
	ii := uint128.Clone()

	for !ii.Equal(zero) {
		tmpRemainder, _ := ii.DivMod(radix)
		remainder, _ := tmpRemainder.ToInt()
		character := "0123456789"[remainder]
		numStr = string(character) + numStr;
	}

	return numStr, nil;
}

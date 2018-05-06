package types

import (
	"errors"
	"encoding/binary"
	"strconv"
)

const (
	lessThan = -1
	equal = 0
	greaterThan = 1

	Len = 32
)

type UInt128 struct {
	lo uint64
	hi uint64
}

func NewUInt128FromInt(val int) *UInt128  { return NewUInt128FromLong(int64(val)) }

func NewUInt128FromLong(val int64) *UInt128 {
	num := new(UInt128)
	if val < 0 {
		num.hi = 1
	} else {
		num.hi = 0
	}
	num.lo = uint64(val)
	return num
}

func NewUInt128FromByteArray(data []byte) (*UInt128, error) {
	if len(data) != 16 {
		return nil, errors.New("UInt128 only can be parsed from 16 byte length array")
	}

	num := new(UInt128)
	num.hi = binary.BigEndian.Uint64(data[0:8])
	num.lo = binary.BigEndian.Uint64(data[8:16])
	return num, nil
}

func (this *UInt128) Clone() UInt128 {
	return UInt128{
		this.lo,
		this.hi,
	}
}

func (this *UInt128) Compare(o *UInt128) int {
	if this.hi < o.hi {
		return lessThan
	} else if this.hi > o.hi {
		return greaterThan
	} else if this.lo < o.lo {
		return lessThan
	} else if this.lo > o.lo {
		return greaterThan
	}

	return equal
}

func (this *UInt128) Equal(o UInt128) bool {
	return this.hi == o.hi && this.lo == o.lo
}

// Get the bit value in the provided position
func (this *UInt128) GetBit(position uint8) uint8 {
	valT := this.hi
	if position > 63 {
		valT = this.lo
		position -= 64
	}
	iShift := 63 - (position % 64);
	return uint8((valT >> iShift) & 1);
}

func And(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.And(val2)
	return newValue
}

func (this *UInt128) And(o UInt128) {
	this.hi &= o.hi
	this.lo &= o.lo
}

func Or(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Or(val2)
	return newValue
}

func (this UInt128) Or(o UInt128) {
	this.hi |= o.hi
	this.lo |= o.lo
}

func Xor(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Xor(val2)
	return newValue
}

func (this *UInt128) Xor(o UInt128) {
	this.hi ^= o.hi
	this.lo ^= o.lo
}

func Add(val1 UInt128, val2 UInt128) UInt128 {
	newValue := val1.Clone()
	newValue.Add(val2)
	return newValue
}

func (this *UInt128) Add(o UInt128) {
	carry := this.lo
	this.lo += o.lo
	this.hi += o.hi

	if this.lo < carry {
		this.hi += 1
	}
}

func (this *UInt128) ToHexString() string {
	return strconv.FormatUint(this.hi, 16) + strconv.FormatUint(this.lo, 16)
}
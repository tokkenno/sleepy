package uint128

import (
	"errors"
	"encoding/binary"
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

func FromInt(val int) *UInt128  { return FromLong(int64(val)) }

func FromLong(val int64) *UInt128 {
	num := new(UInt128)
	if val < 0 {
		num.hi = ^0
	} else {
		num.hi = 0
	}
	num.lo = uint64(val)
	return num
}

func FromByteArray(data []byte) (*UInt128, error) {
	if (len(data) != 8) {
		return nil, errors.New("UInt128 only can be parsed from 8 byte length array")
	}

	num := new(UInt128)
	num.hi = binary.BigEndian.Uint64(data[0:3])
	num.lo = binary.BigEndian.Uint64(data[4:7])
	return num, nil
}

func (u *UInt128) Compare(o *UInt128) int {
	if u.hi < o.hi {
		return lessThan
	} else if u.hi > o.hi {
		return greaterThan
	} else if u.lo < o.lo {
		return lessThan
	} else if u.lo > o.lo {
		return greaterThan
	}

	return equal
}

/*func (u *UInt128) Equal(o *UInt128) bool {
	return u.Compare(o) == equal
}*/

func (u *UInt128) And(o *UInt128) {
	u.hi &= o.hi
	u.lo &= o.lo
}

func (u *UInt128) Or(o *UInt128) {
	u.hi |= o.hi
	u.lo |= u.lo
}

func (u *UInt128) Xor(o *UInt128) {
	u.hi ^= o.hi
	u.lo ^= o.lo
}

func (u *UInt128) Add(o *UInt128) {
	carry := u.lo
	u.lo += o.lo
	u.hi += o.hi

	if u.lo < carry {
		u.hi += 1
	}
}
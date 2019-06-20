package types

import (
	"math"
	"testing"
)

func TestNewUInt128FromInt(t *testing.T) {
	testNum := 75634
	uint1 := NewUInt128(uint64(testNum), 0)
	uint2 := NewUInt128FromInt(testNum)

	if !uint1.Equal(uint2) {
		t.Errorf("uint1 and uint2 must be equals")
	}

	uint1 = NewUInt128(uint64(^uint64(0)-uint64(testNum-1)), ^uint64(0))
	uint2 = NewUInt128FromInt(-testNum)

	if !uint1.Equal(uint2) {
		t.Errorf("uint1 and uint2 must be equals")
	}
}

func TestNewUInt128FromLong(t *testing.T) {
	uint1 := NewUInt128(math.MaxUint64, 0)
	uint2 := NewUInt128FromLong(-1)

	if !uint1.Equal(uint2) {
		t.Errorf("uint1 and uint2 must be equals")
	}
}

func TestUInt128_Equal(t *testing.T) {
	uint1 := NewUInt128FromInt(0)
	uint2 := NewUInt128FromLong(0)
	uint3, err := NewUInt128FromByteArray(make([]byte, 16))
	uint4 := NewUInt128FromInt(1)

	if err != nil {
		t.Errorf("Error while create UInt128 from byte array: %s", err.Error())
	} else {
		if !uint1.Equal(uint2) || !uint2.Equal(uint3) {
			t.Errorf("uint1, uint2 and uint3 must be equals")
		}

		if uint1.Equal(uint4) {
			t.Errorf("uint4 must be different than uint1")
		}
	}
}

func TestUInt128_Clone(t *testing.T) {
	uint1 := NewUInt128FromInt(1)
	uint2 := NewUInt128FromInt(1)

	uint3 := uint1.Clone()

	if !uint1.Equal(uint2) || !uint1.Equal(uint3) {
		t.Errorf("On this step all values must be equals")
	}

	uint3.Add(uint2)

	if !(uint1.Equal(uint2) && !uint3.Equal(uint2)) {
		t.Errorf("The value of uint3 must be different than uint2, but uint1 should be equal")
	}
}

func TestUInt128_Compare(t *testing.T) {
	uint1 := NewUInt128FromInt(1)
	uint2 := NewUInt128FromInt(2)
	uint4 := uint2.Clone()

	if uint1.Compare(uint2) != lessThan {
		t.Errorf("uint1 must be less than uint2")
	}
	if uint2.Compare(uint1) != greaterThan {
		t.Errorf("uint2 must be greater than uint1")
	}
	if uint2.Compare(uint4) != equal {
		t.Errorf("uint2 must be equal than uint3")
	}
}

func TestUInt128_Add(t *testing.T) {
	num1 := uint64(300)
	num2 := uint64(900)

	uint1 := NewUInt128(num1, 1)
	uint2 := NewUInt128(1, num2)
	uint3 := NewUInt128(num1+1, num2+1)

	uint1.Add(uint2)

	if !uint3.Equal(uint1) {
		t.Errorf("ADD operation failed")
	}
}

func TestUInt128_Sub(t *testing.T) {
	uint1 := NewUInt128FromInt(900)
	uint2 := NewUInt128FromInt(300)
	uint3 := NewUInt128FromInt(600)

	uint1.Sub(uint2)

	if !uint3.Equal(uint1) {
		t.Errorf("SUBSTRACT operation failed")
	}
}

func TestUInt128_Mul(t *testing.T) {
	uint1 := NewUInt128FromInt(900)
	uint2 := NewUInt128FromInt(300)
	uint3 := NewUInt128FromInt(900 * 300)

	uint1.Mul(uint2)

	if !uint3.Equal(uint1) {
		t.Errorf("SUBSTRACT operation failed")
	}
}

func TestUInt128_DivMod(t *testing.T) {
	uint1 := NewUInt128FromInt(900)
	uint2 := NewUInt128FromInt(298)

	mod, err := uint1.DivMod(uint2)
	if err != nil {
		t.Errorf("Error while divide the uint128 numbers: " + err.Error())
	} else if !uint1.Equal(NewUInt128FromInt(3)) || !mod.Equal(NewUInt128FromInt(6)) {
		t.Errorf("DIV MOD operation failed")
	}
}

func TestUInt128_Or(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint2, err2 := NewUInt128FromByteArray([]byte{0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x80, 0x00, 0x08, 0x00, 0x8c, 0x00})
	uint3, err3 := NewUInt128FromByteArray([]byte{0xff, 0xff, 0x00, 0xff, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0x80, 0x00, 0x08, 0x00, 0xff, 0x00})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		uint1.Or(uint2)

		if !uint3.Equal(uint1) {
			t.Errorf("OR operation failed")
		}
	}
}

func TestUInt128_And(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint2, err2 := NewUInt128FromByteArray([]byte{0xff, 0xff, 0x00, 0x00, 0xff, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x80, 0x00, 0x08, 0x00, 0x8c, 0x00})
	uint3, err3 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x08, 0x00, 0x00, 0x00})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		uint1.And(uint2)

		if !uint3.Equal(uint1) {
			t.Errorf("AND operation failed")
		}
	}
}

func TestUInt128_Xor(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint2, err2 := NewUInt128FromByteArray([]byte{0xff, 0xff, 0x00, 0x00, 0xff, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x80, 0x00, 0x08, 0x00, 0x8c, 0x00})
	uint3, err3 := NewUInt128FromByteArray([]byte{0x00, 0xff, 0x00, 0xff, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		uint1.Xor(uint2)

		if !uint3.Equal(uint1) {
			t.Errorf("XOR operation failed")
		}
	}
}

func TestUInt128_Not(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint2, err2 := NewUInt128FromByteArray([]byte{0x00, 0xff, 0xff, 0x00, 0xff, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x7f, 0xff, 0xf7, 0xff, 0x8c, 0xff})

	if err1 != nil || err2 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		uint1.Not()

		if !uint1.Equal(uint2) {
			t.Errorf("uint1 is distinct than his negated value on uint2")
		}
	}
}

func TestUInt128_LeftShift(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xaa, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint2, err2 := NewUInt128FromByteArray([]byte{0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00, 0x00, 0x00})
	uint3, err3 := NewUInt128FromByteArray([]byte{0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		uint1.LeftShift(16)

		if !uint1.Equal(uint2) {
			t.Errorf("uint1 must be now equal than uint2")
		}

		uint1.LeftShift(56)

		if !uint1.Equal(uint3) {
			t.Errorf("uint1 must be now equal than uint3")
		}
	}
}

func TestUInt128_RightShift(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xaa, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint2, err2 := NewUInt128FromByteArray([]byte{0x00, 0x00, 0xaa, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00})
	uint3, err3 := NewUInt128FromByteArray([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		uint1.RightShift(16)

		if !uint1.Equal(uint2) {
			t.Errorf("uint1 must be now equal than uint2")
		}

		uint1.RightShift(56)

		if !uint1.Equal(uint3) {
			t.Errorf("uint1 must be now equal than uint3")
		}
	}
}

func TestUInt128_GetBit(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})

	if err1 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		// Test the 0x73
		if uint1.GetBit(128-9) != 1 || uint1.GetBit(128-10) != 1 ||
			uint1.GetBit(128-11) != 0 || uint1.GetBit(128-12) != 0 ||
			uint1.GetBit(128-13) != 1 || uint1.GetBit(128-14) != 1 ||
			uint1.GetBit(128-15) != 1 || uint1.GetBit(128-16) != 0 {
			t.Errorf("The bits retrieved are not the expected")
		}
	}
}

func TestUInt128_ToBitString(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint1Str := "11111111000000000000000011111111000000000000000011110000000000001111000000000000100000000000000000001000000000000111001100000000"

	if err1 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		if uint1.ToBitString() != uint1Str {
			t.Errorf("String conversion to bit string failed")
		}
	}
}

func TestUInt128_ToHexString(t *testing.T) {
	uint1, err1 := NewUInt128FromByteArray([]byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf0, 0x00, 0xf0, 0x00, 0x80, 0x00, 0x08, 0x00, 0x73, 0x00})
	uint1Str := "ff0000ff0000f000f000800008007300"

	if err1 != nil {
		t.Errorf("Error while construct the uint128 numbers")
	} else {
		if uint1.ToHexString() != uint1Str {
			t.Errorf("String conversion to hex string failed")
		}
	}
}

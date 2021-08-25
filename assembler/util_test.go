package assembler

import (
	"testing"
)

func Test_bitsMinMax(t *testing.T) {
	tests := []struct{
		bitsLength uint
		unsigned bool
		min int64
		max uint64
	}{
		{bitsLength: 1, unsigned: true, min: int64(0), max: uint64(1)},
		{bitsLength: 2, unsigned: false, min: int64(-2), max: uint64(1)},
		{bitsLength: 2, unsigned: true, min: int64(0), max: uint64(3)},
		{bitsLength: 15, unsigned: false, min: int64(-16384), max: uint64(16383)},
		{bitsLength: 15, unsigned: true, min: int64(0), max: uint64(32767)},
		{bitsLength: 64, unsigned: false, min: int64(-9223372036854775808), max: uint64(9223372036854775807)},
		{bitsLength: 64, unsigned: true, min: int64(0), max: uint64(18446744073709551615)},
	}

	for _, test := range tests {
		min, max := bitsMinMax(test.bitsLength, test.unsigned)
		if min != test.min || max != test.max {
			t.Errorf("bitsMinMax() = %d, %d, want %d, %d", min, max, test.min, test.max)
		}
	}
}

func Test_NewBits(t *testing.T) {
	tests := []struct{
		originalIntString string
		bitsLength uint
		unsigned bool
		bits Bits
		isErr bool
	}{
		{originalIntString: "10", bitsLength: 0, unsigned: true,
			bits: Bits{}, isErr: true},
		{originalIntString: "10", bitsLength: 1, unsigned: true,
			bits: Bits{}, isErr: true},
		{originalIntString: "-1", bitsLength: 1, unsigned: false,
			bits: Bits{}, isErr: true},
		{originalIntString: "1", bitsLength: 1, unsigned: true,
			bits: Bits{originalUInt64: uint64(1), unsigned: true}, isErr: false},
		{originalIntString: "2", bitsLength: 2, unsigned: false,
			bits: Bits{}, isErr: true},
		{originalIntString: "-3", bitsLength: 2, unsigned: false,
			bits: Bits{}, isErr: true},
		{originalIntString: "1", bitsLength: 2, unsigned: false,
			bits: Bits{originalInt64: int64(1), unsigned: false}, isErr: false},
		{originalIntString: "-2", bitsLength: 2, unsigned: false,
			bits: Bits{originalInt64: int64(-2), unsigned: false}, isErr: false},
		{originalIntString: "3", bitsLength: 2, unsigned: true,
			bits: Bits{originalUInt64: uint64(3), unsigned: true}, isErr: false},
		{originalIntString: "32767", bitsLength: 15, unsigned: true,
			bits: Bits{originalUInt64: uint64(32767), unsigned: true}, isErr: false},
		{originalIntString: "-9223372036854775808", bitsLength: 64, unsigned: false,
			bits: Bits{originalInt64: int64(-9223372036854775808), unsigned: false}, isErr: false},
		{originalIntString: "9223372036854775807", bitsLength: 64, unsigned: false,
			bits: Bits{originalInt64: int64(9223372036854775807), unsigned: false}, isErr: false},
		{originalIntString: "18446744073709551615", bitsLength: 64, unsigned: true,
			bits: Bits{originalUInt64: uint64(18446744073709551615), unsigned: true}, isErr: false},
	}

	for _, test := range tests {
		bits, err := NewBits(test.originalIntString, test.bitsLength, test.unsigned)
		if test.isErr && err == nil {
			t.Errorf("NewBits() should return error: %+v", test)
		}

		if !test.isErr && err != nil {
			t.Errorf("NewBits() should not return error: %s", err)
		}

		if bits != nil {
			if *bits != test.bits {
				t.Errorf("NewBits() = %+v, but want%+v", *bits, test.bits)
			}
		}

	}
}

func Test_notBits(t *testing.T) {
	tests := []struct {
		bits string
		wanted string
	}{
		{bits: "0000", wanted: "1111"},
		{bits: "1111", wanted: "0000"},
		{bits: "10101101", wanted: "01010010"},
		{bits: "", wanted: ""},
	}

	for _, test := range tests {
		if got := notBits(test.bits); got != test.wanted {
			t.Errorf("notBits() = %s, want %s", got, test.wanted)
		}
	}
}

func Test_fillZerosToBits(t *testing.T) {
	tests := []struct {
		originalBits string
		length int
		wanted string
		err bool
	}{
		{originalBits: "1010", length: 2, wanted: "", err: true},
		{originalBits: "1010", length: 4, wanted: "1010", err: false},
		{originalBits: "1010", length: 6, wanted: "001010", err: false},
	}

	for _, test := range tests {
		bits, err := fillZerosToBits(test.originalBits, test.length)
		if err != nil && !test.err {
			t.Errorf("fillZerosToBits() should cause error for %s, %d", test.originalBits, test.length)
		}

		if bits != test.wanted {
			t.Errorf("fillZerosToBits() = %s, want %s", bits, test.wanted)
		}
	}
}

func Test_IntTo15BitsString(t *testing.T) {
	tests := []struct{
		integer int
		wanted string
	}{
		{integer: 16383, wanted: "011111111111111"},
		{integer: 14, wanted: "000000000001110"},
		{integer: 0, wanted: "000000000000000"},
		{integer: -16384, wanted: "100000000000000"},
		{integer: -15, wanted: "111111111110001"},
	}

	for _, test := range tests {
		if bits := IntTo15BitsString(test.integer); bits != test.wanted {
			t.Errorf("InteTo15BitsString() = %s, want %s", bits, test.wanted)
		}
	}
}

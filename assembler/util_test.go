package assembler

import (
	"testing"
)

func Test_bitsMinMax(t *testing.T) {
	tests := []struct{
		bitsLength int
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
		intString string
		length    int
		unsigned  bool
		bits     IntToBitsConverter
		isErr    bool
	}{
		{intString: "10", length: 0, unsigned: true,
			bits: IntToBitsConverter{}, isErr: true},
		{intString: "10", length: 1, unsigned: true,
			bits: IntToBitsConverter{}, isErr: true},
		{intString: "-1", length: 1, unsigned: false,
			bits: IntToBitsConverter{}, isErr: true},
		{intString: "1", length: 1, unsigned: true,
			bits: IntToBitsConverter{originalUInt64: uint64(1), length: 1, unsigned: true}, isErr: false},
		{intString: "2", length: 2, unsigned: false,
			bits: IntToBitsConverter{}, isErr: true},
		{intString: "-3", length: 2, unsigned: false,
			bits: IntToBitsConverter{}, isErr: true},
		{intString: "1", length: 2, unsigned: false,
			bits: IntToBitsConverter{originalInt64: int64(1), length: 2, unsigned: false}, isErr: false},
		{intString: "-2", length: 2, unsigned: false,
			bits: IntToBitsConverter{originalInt64: int64(-2), length: 2, unsigned: false}, isErr: false},
		{intString: "3", length: 2, unsigned: true,
			bits: IntToBitsConverter{originalUInt64: uint64(3), length: 2, unsigned: true}, isErr: false},
		{intString: "32767", length: 15, unsigned: true,
			bits: IntToBitsConverter{originalUInt64: uint64(32767), length: 15, unsigned: true}, isErr: false},
		{intString: "-9223372036854775808", length: 64, unsigned: false,
			bits: IntToBitsConverter{originalInt64: int64(-9223372036854775808), length: 64, unsigned: false}, isErr: false},
		{intString: "9223372036854775807", length: 64, unsigned: false,
			bits: IntToBitsConverter{originalInt64: int64(9223372036854775807), length: 64, unsigned: false}, isErr: false},
		{intString: "18446744073709551615", length: 64, unsigned: true,
			bits: IntToBitsConverter{originalUInt64: uint64(18446744073709551615), length: 64, unsigned: true}, isErr: false},
	}

	for _, test := range tests {
		bits, err := NewIntToBitsConverter(test.intString, test.length, test.unsigned)
		if test.isErr && err == nil {
			t.Errorf("NewIntToBitsConverter() should return error: %+v", test)
		}

		if !test.isErr && err != nil {
			t.Errorf("NewIntToBitsConverter() should not return error: %s", err)
		}

		if bits != nil {
			if *bits != test.bits {
				t.Errorf("NewIntToBitsConverter() = %+v, but want%+v", *bits, test.bits)
			}
		}

	}
}

func TestIntToBitsConverter_ToBits(t *testing.T) {
	tests := []struct{
		intString string
		length int
		unsigned bool
		wanted string
	}{
		{intString: "0", length: 1, unsigned: true, wanted: "0"},
		{intString: "1", length: 1, unsigned: true, wanted: "1"},
		{intString: "-1", length: 2, unsigned: false, wanted: "11"},
		{intString: "1", length: 2, unsigned: false, wanted: "01"},
		{intString: "-2", length: 2, unsigned: false, wanted: "10"},
		{intString: "32767", length: 15, unsigned: true, wanted: "111111111111111"},
		{intString: "14", length: 15, unsigned: true, wanted: "000000000001110"},
		{intString: "0", length: 15, unsigned: true, wanted: "000000000000000"},
		{intString: "-16384", length: 15, unsigned: false, wanted: "100000000000000"},
		{intString: "-15", length: 15, unsigned: false, wanted: "111111111110001"},
		{intString: "-9223372036854775808", length: 64, unsigned: false,
			wanted: "1000000000000000000000000000000000000000000000000000000000000000"},

		{intString: "9223372036854775807", length: 64, unsigned: false,
			wanted: "0111111111111111111111111111111111111111111111111111111111111111"},
		{intString: "18446744073709551615", length: 64, unsigned: true,
			wanted: "1111111111111111111111111111111111111111111111111111111111111111"},
	}

	for _, test := range tests {
		converter, err := NewIntToBitsConverter(test.intString, test.length, test.unsigned)
		if err != nil {
			t.Errorf("IntToBitsConverter.ToBits() results in error: %s", err.Error())
		}
		if got := converter.ToBits(); got != test.wanted {
			t.Errorf("IntToBitsConverter.ToBits() = %s, want %s", got, test.wanted)
		}
	}
}

func TestBits_not(t *testing.T) {
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
		b := IntToBitsConverter{}
		if got := b.not(test.bits); got != test.wanted {
			t.Errorf("notBits() = %s, want %s", got, test.wanted)
		}
	}
}

func TestBits_fillZerosToBits(t *testing.T) {
	tests := []struct {
		originalBits string
		length int
		wanted string
		err bool
	}{
		{originalBits: "1010", length: 4, wanted: "1010", err: false},
		{originalBits: "1010", length: 6, wanted: "001010", err: false},
	}

	for _, test := range tests {
		bits := &IntToBitsConverter{length: test.length}
		if got := bits.fillZerosToBits(test.originalBits); got != test.wanted {
			t.Errorf("IntToBitsConverter.fillZerosToBits() = %s, want %s", got, test.wanted)
		}
	}
}

package assembler

import (
	"testing"
)

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

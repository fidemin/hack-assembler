package assembler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func bitsMinMax(bitsLength uint, unsigned bool) (int64, uint64) {
	if bitsLength < 1 || bitsLength > 64 {
		panic("bitsMinMax(): bitsLength should be between 1 and 64")
	}

	if bitsLength == 1 && !unsigned {
		panic("bitsMinMax(): singed 1 bitsLength is nonsense")
	}

	boundary := int64(1)
	for i := 0; i < int(bitsLength)-1; i++ {
		boundary = boundary * 2
	}

	min := -boundary
	if unsigned {
		min = int64(0)
	}

	max := uint64(boundary - 1)
	if unsigned {
		max = uint64(boundary * 2 - 1)
	}

	return min, max
}

type Bits struct {
	originalInt64 int64
	originalUInt64 uint64
	unsigned bool
}

func NewBits(originalIntString string, bitsLength uint, unsigned bool) (*Bits, error) {
	bits := &Bits{}
	bits.unsigned = unsigned
	if unsigned {
		originalUint64, err := strconv.ParseUint(originalIntString, 10, 64)
		if err != nil {
			return nil, err
		}
		bits.originalUInt64 = originalUint64

	} else {
		originalInt64, err := strconv.ParseInt(originalIntString, 10, 64)
		if err != nil {
			return nil, err
		}
		bits.originalInt64 = originalInt64
	}

	if bitsLength < 1 || bitsLength > 64 {
		return nil, errors.New("bitsLength should be between 1 and 64")
	}

	if bitsLength == 1 && !unsigned {

		return nil, errors.New("singed 1 bitsLength is nonsense")
	}

	min, max := bitsMinMax(bitsLength, unsigned)

	if unsigned {
		if bits.originalUInt64 < uint64(min) || bits.originalUInt64 > max {
			return nil, errors.New(fmt.Sprintf("%d is not between %d, %d", bits.originalUInt64, min, max))
		}
	} else {
		if bits.originalInt64 < min || bits.originalInt64 > int64(max) {
			return nil, errors.New(fmt.Sprintf("%d is not between %d, %d", bits.originalInt64, min, max))
		}
	}
	return bits, nil
}

func notBits(bits string) string {
	length := len(bits)
	var bitsArray = make([]string, length, length)

	for i, b := range bits {
		if string(b) == "0" {
			bitsArray[i] = "1"
		} else if string(b) == "1" {
			bitsArray[i] = "0"
		} else {
			panic(fmt.Sprintf("notBits() bits is not bits: %s", bits))
		}
	}
	return strings.Join(bitsArray, "")
}

func fillZerosToBits(originalBits string, length int) (string, error) {
	originalLength := len(originalBits)

	lengthOfZeros := length - originalLength

	if lengthOfZeros < 0 {
		return "", errors.New("length of original bits is larger than required length")
	}

	if lengthOfZeros == 0 {
		return originalBits, nil
	}

	var bitsArray = make([]string, length, length)
	cursor := 0

	for i:= 0; i < lengthOfZeros; i++ {
		bitsArray[cursor] = "0"
		cursor += 1
	}

	for _, b := range originalBits {
		bitsArray[cursor] = string(b)
		cursor += 1
	}

	return strings.Join(bitsArray, ""), nil
}

func IntTo15BitsString(integer int) string {
	// if not 15bit integer => if not (-2^14 <= integer < 2^14 - 1)
	if integer < -16384 || integer > 16383 {
		panic(fmt.Sprintf("integer %d is not in range of 15bit int", integer))
	}

	if integer < 0 {
		//// For negative integer -K, not(K-1) is binary expression of -K
		//// e.g. -5 -> not(5-1) -> not(4) -> Not(0010) -> 1101
		integer = -integer - 1
		originalBits := strconv.FormatInt(int64(integer), 2)
		bits, err := fillZerosToBits(originalBits, 15)
		if err != nil {
			panic(err.Error())
		}
		return notBits(bits)
	}

	// for zero or positive integer
	originalBits := strconv.FormatInt(int64(integer), 2)
	bits, err := fillZerosToBits(originalBits, 15)
	if err != nil {
		panic(err.Error())
	}
	return bits
}

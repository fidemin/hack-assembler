package assembler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func bitsMinMax(bitsLength int, unsigned bool) (int64, uint64) {
	if bitsLength < 1 || bitsLength > 64 {
		panic("bitsMinMax(): length should be between 1 and 64")
	}

	if bitsLength == 1 && !unsigned {
		panic("bitsMinMax(): singed 1 length is nonsense")
	}

	boundary := int64(1)
	for i := 0; i < bitsLength-1; i++ {
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

// IntToBitsConverter can be used to convert int string to bits string
// integer with size from 1 bit to 64 bit size can be converted!
type IntToBitsConverter struct {
	originalInt64 int64
	originalUInt64 uint64
	length         int
	unsigned       bool
}

func NewIntToBitsConverter(originalIntString string, bitsLength int, unsigned bool) (*IntToBitsConverter, error) {
	bits := &IntToBitsConverter{}
	bits.unsigned = unsigned
	bits.length = bitsLength
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
		return nil, errors.New("length should be between 1 and 64")
	}

	if bitsLength == 1 && !unsigned {
		return nil, errors.New("singed 1 length is nonsense")
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

func (b *IntToBitsConverter) ToBits() string {
	if !b.unsigned {
		if b.originalInt64 < 0 {
			// For negative integer -K, not(K-1) is binary expression of -K
			// e.g. -5 -> not(5-1) -> not(4) -> Not(0010) -> 1101
			integer := -b.originalInt64 - 1
			bitsString := strconv.FormatInt(integer, 2)
			bitsString = b.fillZerosToBits(bitsString)
			return b.not(bitsString)
		} else {
			bitsString := strconv.FormatInt(b.originalInt64, 2)
			return b.fillZerosToBits(bitsString)
		}
	} else {
		bitsString := strconv.FormatUint(b.originalUInt64, 2)
		return b.fillZerosToBits(bitsString)
	}
}

func (b *IntToBitsConverter) fillZerosToBits(originalBits string) string {
	originalLength := len(originalBits)

	lengthOfZeros := b.length - originalLength

	if lengthOfZeros < 0 {
		panic("length of originalBits should be smaller or equal than length")
	}

	if lengthOfZeros == 0 {
		return originalBits
	}

	var bitsArray = make([]string, b.length, b.length)
	cursor := 0

	for i:= 0; i < lengthOfZeros; i++ {
		bitsArray[cursor] = "0"
		cursor += 1
	}

	for _, b := range originalBits {
		bitsArray[cursor] = string(b)
		cursor += 1
	}

	return strings.Join(bitsArray, "")
}

func (b *IntToBitsConverter) not(bits string) string {
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

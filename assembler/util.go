package assembler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

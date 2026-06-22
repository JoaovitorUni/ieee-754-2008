package ieee

import (
	"fmt"
	"math"
	"math/bits"
)

type IEEENumber struct {
	IEEENumber uint32
	Signal     uint  // 1 bit
	Expoent    uint8 // 8 bits
	Mantissa   uint  // 23 bits
}

func NewIEEENumber(data any) (*IEEENumber, error) {
	switch v := data.(type) {
	case int:
		if v == 0 {
			return createIEEENumber(0), nil
		}
		return fromInt(v)
	case float32:
		return fromFloat32(v)
	default:
		return nil, fmt.Errorf("Tipo %t não é aceito. Tipos aceitos: string, int e float32.", data)
	}
}

func createIEEENumber(n uint32) *IEEENumber {
	return &IEEENumber{
		IEEENumber: n,
		Signal:     uint((n & signalMask) >> 31),
		Expoent:    uint8((n & exponentMask) >> 23),
		Mantissa:   uint((n & mantissaMask)),
	}
}

func fromInt(value int) (*IEEENumber, error) {
	var iNumber uint32 = 0
	if value < 0 {
		iNumber |= signalMask
		value *= -1
	}
	uNumber := uint(value)

	msbPos := bits.Len(uNumber)
	exponent := int(bias)
	if msbPos != 0 {
		exponent += msbPos - 1
	}
	exponentMask := uint32(exponent) << 23

	distance := 23 - msbPos
	var mantissa uint32
	if distance < 0 {
		distance *= -1
		mantissa = uint32(uNumber >> distance)
	} else {
		mantissa = uint32(uNumber << distance)
	}
	mantissa <<= 1 // Remove bit implícito
	mantissa &= mantissaMask

	iNumber |= exponentMask
	iNumber |= mantissa
	return createIEEENumber(iNumber), nil
}
func fromFloat32(value float32) (*IEEENumber, error) {
	iNumber := math.Float32bits(value)
	return createIEEENumber(iNumber), nil
}

func (i *IEEENumber) ToInt(t Rounding) int {
	notSignalMask := ^signalMask // NOT em signalMask
	if (i.IEEENumber & notSignalMask) == 0 {
		return 0
	}
	uNumber := i.MantissaWithHideOne()
	fractionalBitsCount := int(23 - i.ExpoentToInt())

	var number int

	if fractionalBitsCount <= 0 {
		fractionalBitsCount *= -1
		number = int(uNumber << fractionalBitsCount)
	} else {
		if t == Truncate {
			number = i.truncate(uint(uNumber), fractionalBitsCount)
		}
		if t == Nearest {
			number = i.roundNearest(uint(uNumber), fractionalBitsCount)
		}
	}

	if i.Signal > 0 {
		number *= -1
	}
	return number
}
func (i *IEEENumber) ToFloat32() float32 {
	return math.Float32frombits(i.IEEENumber)
}

func (i *IEEENumber) truncate(uNumber uint, fractionalBitsCount int) int {
	return int(uNumber >> fractionalBitsCount)
}
func (i *IEEENumber) roundNearest(uNumber uint, fractionalBitsCount int) int {
	tNumber := i.truncate(uNumber, fractionalBitsCount)

	firstFractionalBitsCount := uint(fractionalBitsCount - 1)

	lsb := tNumber & 1                                         // Least Significant Bit (LSB): Último bit da parte inteira
	guardBit := int((uNumber >> firstFractionalBitsCount) & 1) // Guard Bit (G): Primeiro fracionário
	var stickyBit int = 0                                      // Sticky Bit (S): Se 1 após o Guard Bit temos outros bits

	if fractionalBitsCount > 1 {
		stickyMask := uint32((1 << firstFractionalBitsCount) - 1)
		if (uint32(uNumber) & stickyMask) != 0 {
			stickyBit = 1
		}
	}

	// O arredondamento do IEEE 754 é chamado de "Ties to Even"
	// Isso significa que o arredondamento para cima no caso de empate exato, primeiro fracionário é 5,
	// só é feito se o número for impar. Se não for empate exato verificamos o Stick Bit.
	if guardBit == 1 && (stickyBit == 1 || lsb == 1) {
		tNumber++
	}
	return tNumber
}

func (i *IEEENumber) ExpoentToInt() int8 {
	return int8((i.Expoent)) - bias
}
func (i *IEEENumber) MantissaWithHideOne() uint {
	return uint((i.IEEENumber & mantissaMask) | hideOneMask)
}

func (i *IEEENumber) Debug() {
	fmt.Printf("Sinal: %b\n", i.IEEENumber>>31)
	fmt.Printf("Expoente: %08b (%d) Expoente Real: %d\n", i.Expoent, (i.Expoent), i.ExpoentToInt())
	fmt.Printf("Mantissa: %023b (%d)\n", i.Mantissa, i.Mantissa)
}
func (i *IEEENumber) String() string {
	strValue := fmt.Sprintf("%032b", i.IEEENumber)
	return fmt.Sprintf("IEE 754-2008: \033[34m%c\033[0m\033[31m%s\033[0m\033[33m%s\033[0m", strValue[0], strValue[1:9], strValue[9:])
}

package ieee

import (
	"fmt"
	"math/bits"
)

func Sum(iNumber1 *IEEENumber, iNumber2 *IEEENumber) (*IEEENumber, error) {
	baseExpoent, m1, m2 := nomalizeIEEENumbers(iNumber1, iNumber2)

	if iNumber1.Signal == iNumber2.Signal {
		return addMantissas(baseExpoent, m1, m2, iNumber1.Signal)
	}
	return subtractMantissas(baseExpoent, m1, m2, iNumber1.Signal)
}

func Sub(iNumber1 *IEEENumber, iNumber2 *IEEENumber) (*IEEENumber, error) {
	baseExpoent, m1, m2 := nomalizeIEEENumbers(iNumber1, iNumber2)

	if iNumber1.Signal != iNumber2.Signal {
		return addMantissas(baseExpoent, m1, m2, iNumber1.Signal)
	}
	return subtractMantissas(baseExpoent, m1, m2, iNumber1.Signal)
}

func Mult(iNumber1 *IEEENumber, iNumber2 *IEEENumber) (*IEEENumber, error) {
	s := iNumber1.Signal ^ iNumber2.Signal
	e := int(iNumber1.Expoent) + int(iNumber2.Expoent) - 127
	m64 := uint64(iNumber1.MantissaWithHideOne()) * uint64(iNumber2.MantissaWithHideOne())

	if m64&(1<<47) != 0 {
		m64 >>= 1
		e++
	}

	m := uint(m64>>23) & uint(mantissaMask)
	return mathResult(e, m, s)
}

func Div(iNumber1 *IEEENumber, iNumber2 *IEEENumber) (*IEEENumber, error) {
	s := iNumber1.Signal ^ iNumber2.Signal
	e := int(iNumber1.Expoent) - int(iNumber2.Expoent) + 127
	m1 := uint64(iNumber1.MantissaWithHideOne()) << 23
	m2 := uint64(iNumber2.MantissaWithHideOne())
	m64 := m1 / m2

	if m64&(1<<23) == 0 {
		m64 <<= 1
		e--
	}

	m := uint(m64) & uint(mantissaMask)
	return mathResult(e, m, s)
}

func addMantissas(baseExpoent uint8, m1 uint, m2 uint, resultSignal uint) (*IEEENumber, error) {
	e := int(baseExpoent)
	newMantissa := m1 + m2

	if newMantissa&(1<<24) != 0 {
		newMantissa >>= 1
		e++
	}

	mantissa := newMantissa & uint(mantissaMask)
	return mathResult(e, mantissa, resultSignal)
}

func subtractMantissas(baseExpoent uint8, m1 uint, m2 uint, defaultSignal uint) (*IEEENumber, error) {
	if m1 == m2 {
		return createIEEENumber(0), nil
	}

	var newMantissa uint
	var resultSignal uint = defaultSignal

	if m1 > m2 {
		newMantissa = m1 - m2
	} else {
		newMantissa = m2 - m1
		resultSignal = defaultSignal ^ 1 // Ajustado para inverter o bit de sinal de forma limpa
	}

	e := int(baseExpoent)
	msbPos := bits.Len(uint(newMantissa))
	if msbPos < 24 {
		shift := 24 - msbPos
		newMantissa <<= uint(shift)
		e -= shift
	}

	mantissa := newMantissa & uint(mantissaMask)
	return mathResult(e, mantissa, resultSignal)
}

func nomalizeIEEENumbers(n1 *IEEENumber, n2 *IEEENumber) (uint8, uint, uint) {
	exp1 := n1.ExpoentToInt()
	exp2 := n2.ExpoentToInt()

	if exp1 > exp2 {
		eDiff := uint(exp1 - exp2)
		m1 := uint(n1.MantissaWithHideOne())
		m2 := uint(n2.MantissaWithHideOne()) >> eDiff
		return n1.Expoent, m1, m2
	}

	eDiff := uint(exp2 - exp1)
	m2 := uint(n2.MantissaWithHideOne())
	m1 := uint(n1.MantissaWithHideOne()) >> eDiff
	return n2.Expoent, m1, m2
}

func mathResult(e int, m uint, s uint) (*IEEENumber, error) {
	if e > 254 { // Overflow
		return createIEEENumber(createINumber(255, 0, s)), fmt.Errorf("%s", Overflow)
	}
	if e < 1 { // Underflow
		return createIEEENumber(0), fmt.Errorf("%s", Underflow)
	}
	return createIEEENumber(createINumber(uint8(e), m, s)), nil
}

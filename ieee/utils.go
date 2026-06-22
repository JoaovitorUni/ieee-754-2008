package ieee

type Rounding string

const (
	Truncate Rounding = "T"
	Nearest  Rounding = "N"
)

type MathError string

const (
	Overflow  MathError = "OVERFLOW"
	Underflow MathError = "UNDERFLOW"
)

const (
	signalMask   uint32 = 1 << 31
	exponentMask uint32 = 0xff << 23
	hideOneMask  uint32 = 0x800000
	mantissaMask uint32 = 0x7fffff
	bias         int8   = 127
)

func createINumber(expoent uint8, mantissa uint, signal uint) uint32 {
	var iNumber uint32 = 0
	iNumber |= uint32(signal) << 31
	iNumber |= uint32(expoent) << 23
	iNumber |= uint32(mantissa)
	return iNumber
}

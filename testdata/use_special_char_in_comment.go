package testdata

// SpecialCharType is SpecialCharType.
// +enameg
type SpecialCharType int

// SpecialCharTypes
const (
	SpecialCharTypeBackSlash   SpecialCharType = 0 // A\B
	SpecialCharTypeDoubleQuote SpecialCharType = 1 // "B"
)

//go:generate enameg $GOFILE

package testdata

// HogeType is HogeType.
// +enameg
type HogeType int

// HogeTypes
const (
	HogeTypeA HogeType = 0 // A
	HogeTypeB HogeType = 1 // B
	HogeTypeC HogeType = 2 // C
)

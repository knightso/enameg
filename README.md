enameg
======

enameg is a enum names generator.

# Install

```sh
$ go get github.com/knightso/enameg/cmd/enameg
```

# How to use

## Usage

```sh
$ ename [-flag] [directory]
$ ename [-flag] files... # Must be a single package

-flag:
  -output="": file name; default srcdir/<filename>_ename.go
```

## Annotation

Thie generator generates `Name` method for the type which have an annotation comment `// +enameg`.

from:

```go
// +enameg
type HogeType int

const (
	HogeTypeA HogeType = 0 // A
	HogeTypeB HogeType = 1 // B
	HogeTypeC HogeType = 2 // C
)
```

to:

```go
// Name returns the HogeType Name.
func (src HogeType) Name() string {
	switch src {
	case HogeTypeA:
		return "A"
	case HogeTypeB:
		return "B"
	case HogeTypeC:
		return "C"
	default:
		return fmt.Sprintf("%v", src)
	}
}
```

# With `go generate`

```go
//go:generate enameg $GOFILE

// +enameg
type HogeType int

const (
	HogeTypeA HogeType = 0 // A
	HogeTypeB HogeType = 1 // B
	HogeTypeC HogeType = 2 // C
)
```


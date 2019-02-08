package enameg_test

import (
	"strings"
	"testing"

	"github.com/knightso/enameg"
)

func TestSimple(t *testing.T) {
	testCases := []struct {
		Name     string
		Path     string
		Expected string // ignore if empty
	}{
		{"simple", "./testdata/simple.go", expectedForSimple},
		{"special char", "./testdata/use_special_char_in_comment.go", expectedForSpecialChar},
		{"lack comment", "./testdata/lack_comment.go", expectedForLackComment},
		{"without constants", "./testdata/without_constants.go", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic: %#v", r)
				}
			}()

			_, g := enameg.Generate([]string{tc.Path})

			expected := strings.TrimSpace(tc.Expected)
			g = strings.TrimSpace(g)
			if expected != "" && g != expected {
				t.Errorf("generated = \n%s, wants = \n%s", g, expected)
			}
		})
	}
}

const expectedForSimple = `
package testdata

import "fmt"

// Name returns the SimpleType Name.
func (src SimpleType) Name() string {
	switch src {
	case SimpleTypeA:
		return "A"
	case SimpleTypeB:
		return "B"
	case SimpleTypeC:
		return "C"
	default:
		return fmt.Sprintf("%v", src)
	}
}
`

const expectedForSpecialChar = `
package testdata

import "fmt"

// Name returns the SpecialCharType Name.
func (src SpecialCharType) Name() string {
	switch src {
	case SpecialCharTypeBackSlash:
		return "A\\B"
	case SpecialCharTypeDoubleQuote:
		return "\"B\""
	case SpecialCharTypePercent:
		return "C%%D"
	default:
		return fmt.Sprintf("%v", src)
	}
}
`

const expectedForLackComment = `
package testdata

import "fmt"

// Name returns the LackCommentType Name.
func (src LackCommentType) Name() string {
	switch src {
	case LackCommentTypeA:
		return "A"
	case LackCommentTypeC:
		return "C"
	default:
		return fmt.Sprintf("%v", src)
	}
}
`

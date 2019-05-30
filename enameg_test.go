package enameg_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/knightso/enameg"
)

func TestSimple(t *testing.T) {
	testCases := []struct {
		Name                  string
		Path                  string
		optionNoFmt           bool
		optionEnableEmptyName bool
		Expected              string // ignore if empty
	}{
		{"simple", "./testdata/simple.go", false, false, expectedForSimple},
		{"simple", "./testdata/simple.go", true, false, expectedForSimple},
		{"simple", "./testdata/simple.go", false, true, expectedForSimple},
		{"special char", "./testdata/use_special_char_in_comment.go", false, false, expectedForSpecialChar},
		{"special char", "./testdata/use_special_char_in_comment.go", true, false, expectedForSpecialChar},
		{"special char", "./testdata/use_special_char_in_comment.go", false, true, expectedForSpecialChar},
		{"lack comment", "./testdata/lack_comment.go", false, false, expectedForLackComment},
		{"lack comment", "./testdata/lack_comment.go", true, false, expectedForLackComment},
		{"lack comment", "./testdata/lack_comment.go", false, true, expectedForLackCommentForEmptyName},
		{"without constants", "./testdata/without_constants.go", false, false, ""},
		{"without constants", "./testdata/without_constants.go", true, false, ""},
		{"without constants", "./testdata/without_constants.go", false, true, ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s,--nofmt=%t,--empty-nocomment=%t", tc.Name, tc.optionNoFmt, tc.optionEnableEmptyName), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic: %#v", r)
				}
			}()

			_, g := enameg.Generate([]string{tc.Path}, tc.optionNoFmt, tc.optionEnableEmptyName)

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

import (
	"fmt"
)

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

import (
	"fmt"
)

// Name returns the SpecialCharType Name.
func (src SpecialCharType) Name() string {
	switch src {
	case SpecialCharTypeBackSlash:
		return "A\\B"
	case SpecialCharTypeDoubleQuote:
		return "\"B\""
	default:
		return fmt.Sprintf("%v", src)
	}
}
`

const expectedForLackComment = `
package testdata

import (
	"fmt"
)

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

const expectedForLackCommentForEmptyName = `
package testdata

import (
	"fmt"
)

// Name returns the LackCommentType Name.
func (src LackCommentType) Name() string {
	switch src {
	case LackCommentTypeA:
		return "A"
	case LackCommentTypeB:
		return ""
	case LackCommentTypeC:
		return "C"
	default:
		return fmt.Sprintf("%v", src)
	}
}
`

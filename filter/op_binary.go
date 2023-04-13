package filter

import (
	"fmt"
	fp "github.com/scim2/filter-parser/v2"
	"strings"
)

// cmpBinary returns a compare function that compares a given value to the reference string based on the given attribute
// expression and binary attribute. The filter operators gt, lt, ge and le are not supported on binary attributes.
//
// Expects a binary attribute. Will panic on unknown filter operator.
// Known operators: eq, ne, co, sw, ew, gt, lt, ge and le.
func cmpBinary(e *fp.AttributeExpression, ref string) (func(interface{}) error, error) {
	switch op := e.Operator; op {
	case fp.EQ:
		return cmpStr(ref, true, func(v, ref string) error {
			if v != ref {
				return fmt.Errorf("%s is not equal to %s", v, ref)
			}
			return nil
		})
	case fp.NE:
		return cmpStr(ref, true, func(v, ref string) error {
			if v == ref {
				return fmt.Errorf("%s is equal to %s", v, ref)
			}
			return nil
		})
	case fp.CO:
		return cmpStr(ref, true, func(v, ref string) error {
			if !strings.Contains(v, ref) {
				return fmt.Errorf("%s does not contain %s", v, ref)
			}
			return nil
		})
	case fp.SW:
		return cmpStr(ref, true, func(v, ref string) error {
			if !strings.HasPrefix(v, ref) {
				return fmt.Errorf("%s does not start with %s", v, ref)
			}
			return nil
		})
	case fp.EW:
		return cmpStr(ref, true, func(v, ref string) error {
			if !strings.HasSuffix(v, ref) {
				return fmt.Errorf("%s does not end with %s", v, ref)
			}
			return nil
		})
	case fp.GT, fp.LT, fp.GE, fp.LE:
		return nil, fmt.Errorf("can not use op %q on binary values", op)
	default:
		panic(fmt.Sprintf("unknown operator in expression: %s", e))
	}
}

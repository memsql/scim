package filter

import (
	"fmt"
	fp "github.com/scim2/filter-parser/v2"
	"strings"
)

func cmpInt(ref int, cmp func(v, ref int) error) func(interface{}) error {
	return func(i interface{}) error {
		v, ok := i.(int)
		if !ok {
			panic(fmt.Sprintf("given value is not an integer: %v", i))
		}
		return cmp(v, ref)
	}
}

func cmpIntStr(ref int, cmp func(v, ref string) error) (func(interface{}) error, error) {
	return func(i interface{}) error {
		if _, ok := i.(int); !ok {
			panic(fmt.Sprintf("given value is not an integer: %v", i))
		}
		return cmp(fmt.Sprintf("%d", i), fmt.Sprintf("%d", ref))
	}, nil
}

// cmpInteger returns a compare function that compares a given value to the reference int based on the given attribute
// expression and integer attribute.
//
// Expects a integer attribute. Will panic on unknown filter operator.
// Known operators: eq, ne, co, sw, ew, gt, lt, ge and le.
func cmpInteger(e *fp.AttributeExpression, ref int) (func(interface{}) error, error) {
	switch op := e.Operator; op {
	case fp.EQ:
		return cmpInt(ref, func(v, ref int) error {
			if v != ref {
				return fmt.Errorf("%d is not equal to %d", v, ref)
			}
			return nil
		}), nil
	case fp.NE:
		return cmpInt(ref, func(v, ref int) error {
			if v == ref {
				return fmt.Errorf("%d is equal to %d", v, ref)
			}
			return nil
		}), nil
	case fp.CO:
		return cmpIntStr(ref, func(v, ref string) error {
			if !strings.Contains(v, ref) {
				return fmt.Errorf("%s does not contain %s", v, ref)
			}
			return nil
		})
	case fp.SW:
		return cmpIntStr(ref, func(v, ref string) error {
			if !strings.HasPrefix(v, ref) {
				return fmt.Errorf("%s does not start with %s", v, ref)
			}
			return nil
		})
	case fp.EW:
		return cmpIntStr(ref, func(v, ref string) error {
			if !strings.HasSuffix(v, ref) {
				return fmt.Errorf("%s does not end with %s", v, ref)
			}
			return nil
		})
	case fp.GT:
		return cmpInt(ref, func(v, ref int) error {
			if v <= ref {
				return fmt.Errorf("%d is not greater than %d", v, ref)
			}
			return nil
		}), nil
	case fp.LT:
		return cmpInt(ref, func(v, ref int) error {
			if v >= ref {
				return fmt.Errorf("%d is not less than %d", v, ref)
			}
			return nil
		}), nil
	case fp.GE:
		return cmpInt(ref, func(v, ref int) error {
			if v < ref {
				return fmt.Errorf("%d is not greater or equal to %d", v, ref)
			}
			return nil
		}), nil
	case fp.LE:
		return cmpInt(ref, func(v, ref int) error {
			if v > ref {
				return fmt.Errorf("%d is not less or equal to %d", v, ref)
			}
			return nil
		}), nil
	default:
		panic(fmt.Sprintf("unknown operator in expression: %s", e))
	}
}

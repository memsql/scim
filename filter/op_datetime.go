package filter

import (
	"fmt"
	datetime "github.com/di-wu/xsd-datetime"
	fp "github.com/scim2/filter-parser/v2"
	"strings"
	"time"
)

// cmpDateTime returns a compare function that compares a given value to the reference string/time based on the given
// attribute expression and dateTime attribute.
//
// Expects a dateTime attribute. Will panic on unknown filter operator.
// Known operators: eq, ne, co, sw, ew, gt, lt, ge and le.
func cmpDateTime(e *fp.AttributeExpression, date string, ref time.Time) (func(interface{}) error, error) {
	switch op := e.Operator; op {
	case fp.EQ:
		return cmpTime(ref, func(v, ref time.Time) error {
			if !v.Equal(ref) {
				return fmt.Errorf("%s is not equal to %s", v.Format(time.RFC3339), ref.Format(time.RFC3339))
			}
			return nil
		}), nil
	case fp.NE:
		return cmpTime(ref, func(v, ref time.Time) error {
			if v.Equal(ref) {
				return fmt.Errorf("%s is equal to %s", v.Format(time.RFC3339), ref.Format(time.RFC3339))
			}
			return nil
		}), nil
	case fp.CO:
		return cmpStr(date, false, func(v, ref string) error {
			if !strings.Contains(v, ref) {
				return fmt.Errorf("%s does not contain %s", v, ref)
			}
			return nil
		})
	case fp.SW:
		return cmpStr(date, false, func(v, ref string) error {
			if !strings.HasPrefix(v, ref) {
				return fmt.Errorf("%s does not start with %s", v, ref)
			}
			return nil
		})
	case fp.EW:
		return cmpStr(date, false, func(v, ref string) error {
			if !strings.HasSuffix(v, ref) {
				return fmt.Errorf("%s does not end with %s", v, ref)
			}
			return nil
		})
	case fp.GT:
		return cmpTime(ref, func(v, ref time.Time) error {
			if !v.After(ref) {
				return fmt.Errorf("%s is not greater than %s", v.Format(time.RFC3339), ref.Format(time.RFC3339))
			}
			return nil
		}), nil
	case fp.LT:
		return cmpTime(ref, func(v, ref time.Time) error {
			if !v.Before(ref) {
				return fmt.Errorf("%s is not less than %s", v.Format(time.RFC3339), ref.Format(time.RFC3339))
			}
			return nil
		}), nil
	case fp.GE:
		return cmpTime(ref, func(v, ref time.Time) error {
			if !v.After(ref) && !v.Equal(ref) {
				return fmt.Errorf("%s is not greater or equal to %s", v.Format(time.RFC3339), ref.Format(time.RFC3339))
			}
			return nil
		}), nil
	case fp.LE:
		return cmpTime(ref, func(v, ref time.Time) error {
			if !v.Before(ref) && !v.Equal(ref) {
				return fmt.Errorf("%s is not less or equal to %s", v.Format(time.RFC3339), ref.Format(time.RFC3339))
			}
			return nil
		}), nil
	default:
		panic(fmt.Sprintf("unknown operator in expression: %s", e))
	}
}

func cmpTime(ref time.Time, cmp func(v, ref time.Time) error) func(interface{}) error {
	return func(i interface{}) error {
		date, ok := i.(string)
		if !ok {
			panic(fmt.Sprintf("given value is not a string: %v", i))
		}
		value, err := datetime.Parse(date)
		if err != nil {
			panic(fmt.Sprintf("given value is not a date time (%v): %s", i, err))
		}
		return cmp(value, ref)
	}
}

package filter_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim/filter"

	"github.com/elimity-com/scim/schema"
	fp "github.com/scim2/filter-parser/v2"
)

func TestValidatorBoolean(t *testing.T) {
	var (
		exp = func(op fp.CompareOperator) string {
			return fmt.Sprintf("bool %s true", op)
		}
		ref = schema.Schema{
			Attributes: []schema.CoreAttribute{
				schema.SimpleCoreAttribute(schema.SimpleBooleanParams(schema.BooleanParams{
					Name: "bool",
				})),
			},
		}
		attr = map[string]interface{}{
			"bool": true,
		}
	)

	for _, test := range []struct {
		op    fp.CompareOperator
		valid bool // Whether the filter is valid.
	}{
		{fp.EQ, true},
		{fp.NE, false},
		{fp.CO, true},
		{fp.SW, true},
		{fp.EW, true},
		{fp.GT, false},
		{fp.LT, false},
		{fp.GE, false},
		{fp.LE, false},
	} {
		t.Run(string(test.op), func(t *testing.T) {
			f := exp(test.op)
			validator, err := filter.NewValidator(f, ref)
			if err != nil {
				t.Fatal(err)
			}
			if err := validator.PassesFilter(attr); (err == nil) != test.valid {
				t.Errorf("%s %v | actual %v, expected %v", f, attr, err, test.valid)
			}
		})
	}
}

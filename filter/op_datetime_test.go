package filter_test

import (
	"fmt"
	"github.com/elimity-com/scim/filter"
	"github.com/elimity-com/scim/schema"
	fp "github.com/scim2/filter-parser/v2"
	"testing"
)

func TestValidatorDateTime(t *testing.T) {
	var (
		exp = func(op fp.CompareOperator) string {
			return fmt.Sprintf("time %s \"2021-01-01T12:00:00Z\"", op)
		}
		ref = schema.Schema{
			Attributes: []schema.CoreAttribute{
				schema.SimpleCoreAttribute(schema.SimpleDateTimeParams(schema.DateTimeParams{
					Name: "time",
				})),
			},
		}
		attrs = [3]map[string]interface{}{
			{"time": "2021-01-01T08:00:00Z"}, // before
			{"time": "2021-01-01T12:00:00Z"}, // equal
			{"time": "2021-01-01T16:00:00Z"}, // after
		}
	)

	for _, test := range []struct {
		op    fp.CompareOperator
		valid [3]bool
	}{
		{fp.EQ, [3]bool{false, true, false}},
		{fp.NE, [3]bool{true, false, true}},
		{fp.CO, [3]bool{false, true, false}},
		{fp.SW, [3]bool{false, true, false}},
		{fp.EW, [3]bool{false, true, false}},
		{fp.GT, [3]bool{false, false, true}},
		{fp.LT, [3]bool{true, false, false}},
		{fp.GE, [3]bool{false, true, true}},
		{fp.LE, [3]bool{true, true, false}},
	} {
		t.Run(string(test.op), func(t *testing.T) {
			f := exp(test.op)
			validator, err := filter.NewValidator(f, ref)
			if err != nil {
				t.Fatal(err)
			}
			for i, attr := range attrs {
				if err := validator.PassesFilter(attr); (err == nil) != test.valid[i] {
					t.Errorf("(%d) %s %v | actual %v, expected %v", i, f, attr, err, test.valid[i])
				}
			}
		})
	}
}

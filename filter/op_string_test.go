package filter_test

import (
	"fmt"
	"github.com/elimity-com/scim/filter"
	"github.com/elimity-com/scim/schema"
	fp "github.com/scim2/filter-parser/v2"
	"testing"
)

func TestValidatorString(t *testing.T) {
	var (
		exp = func(op fp.CompareOperator) string {
			return fmt.Sprintf("str %s \"x\"", op)
		}
		attrs = [3]map[string]interface{}{
			{"str": "x"},
			{"str": "X"},
			{"str": "y"},
		}
	)

	for _, test := range []struct {
		op      fp.CompareOperator
		valid   [3]bool
		validCE [3]bool
	}{
		{fp.EQ, [3]bool{true, true, false}, [3]bool{true, false, false}},
		{fp.NE, [3]bool{false, false, true}, [3]bool{false, true, true}},
		{fp.CO, [3]bool{true, true, false}, [3]bool{true, false, false}},
		{fp.SW, [3]bool{true, true, false}, [3]bool{true, false, false}},
		{fp.EW, [3]bool{true, true, false}, [3]bool{true, false, false}},
		{fp.GT, [3]bool{false, false, true}, [3]bool{false, false, true}},
		{fp.LT, [3]bool{false, false, false}, [3]bool{false, true, false}},
		{fp.GE, [3]bool{true, true, true}, [3]bool{true, false, true}},
		{fp.LE, [3]bool{true, true, false}, [3]bool{true, true, false}},
	} {
		t.Run(string(test.op), func(t *testing.T) {
			f := exp(test.op)
			for i, attr := range attrs {
				validator, err := filter.NewValidator(f, schema.Schema{
					Attributes: []schema.CoreAttribute{
						schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
							Name: "str",
						})),
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if err := validator.PassesFilter(attr); (err == nil) != test.valid[i] {
					t.Errorf("(0.%d) %s %s | actual %v, expected %v", i, f, attr, err, test.valid[i])
				}
				validatorCE, err := filter.NewValidator(f, schema.Schema{
					Attributes: []schema.CoreAttribute{
						schema.SimpleCoreAttribute(schema.SimpleReferenceParams(schema.ReferenceParams{
							Name: "str",
						})),
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if err := validatorCE.PassesFilter(attr); (err == nil) != test.validCE[i] {
					t.Errorf("(1.%d) %s %s | actual %v, expected %v", i, f, attr, err, test.validCE[i])
				}
			}
		})
	}
}

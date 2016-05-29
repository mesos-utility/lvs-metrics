package g

import (
	"strings"
	"testing"
)

func TestCalculateMetricRatio(t *testing.T) {
	type CalMetricRadioTest struct {
		result string
		opstr  string
	}

	var tests = []CalMetricRadioTest{
		{"0.0", "0.0"},
		{"0", "0.0,0.0"},
		{"0.00", "0.0,1.0"},
		{"100.00", "1.0,0.0"},
		{"50.00", "1.0, 1.0"},
		{"0", "0.0,a.0"},
	}

	for _, test := range tests {
		array := strings.SplitN(test.opstr, ",", -1)
		ret := CalculateMetricRatio(array...)

		if test.result != ret {
			t.Errorf("CalculateMetricRatio(%v) = %v, want %v", array, ret, test.result)
		}
	}
}

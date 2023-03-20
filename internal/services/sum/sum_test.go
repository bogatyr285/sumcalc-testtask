package sum

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestSum(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		exp  int
	}{
		{"simple array of integers", `[1,2,3,4]`, 10},
		{"simple object of integers", `{"a":6,"b":4}`, 10},
		{"nested array of integers", `[[[2]]]`, 2},
		{"nested object of integers", `{"a":{"b":4},"c":-2}`, 2},
		{"empty array", `[]`, 0},
		{"empty object", `{}`, 0},
		{"array of mixed values", `{"a":[-1,1,"dark"]}`, 0},
		{"object of mixed values", `[-1,{"a":1, "b":"light"}]`, 0},
		{"empty string", ``, 0},
		{"bool", `false`, 0},
		{"mix of nums & arrs", `[1, [2, 3], 4]`, 10},
		{"mix of objs & arrs", `{"a": [1, 2], "b": [3, 4]}`, 10},
		// string should be considered
		{"mix of numbers, strings, objects and arrays", `[1, "2", {"a": 3}, [4, 5]]`, 13},
		{"deep nested", `[1, {"nested": {"multiple":{"levels":1334}}}, 2]`, 1337},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var arg interface{}
			err := json.NewDecoder(strings.NewReader(tc.in)).Decode(&arg)
			if err != nil {
				return
			}

			res := NewSumCalculator().SumNumbers(arg)
			if res != tc.exp {
				t.Errorf("Expected: %d, Got: %d", tc.exp, res)
			}
		})
	}
}

// benchmarks for testing/comparison purposes of another solutions & algorithms
func BenchmarkSumSmallInput(b *testing.B) {
	arg := []interface{}{
		map[string]interface{}{
			"a": map[string]interface{}{
				"b": 4,
			},
			"c": -2,
		},
	}

	for i := 0; i < b.N; i++ {
		NewSumCalculator().SumNumbers(arg)
	}
}

func BenchmarkSumHugeInput(b *testing.B) {
	arg := GenerateHugeJSON()

	for i := 0; i < b.N; i++ {
		NewSumCalculator().SumNumbers(arg)
	}
}

// GenerateHugeJSON generates a huge JSON object
func GenerateHugeJSON() interface{} {
	data := make(map[string]interface{})
	for i := 0; i < 10000; i++ {
		data[fmt.Sprintf("key_%d", i)] = i
	}
	return data
}

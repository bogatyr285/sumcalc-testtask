package sum

type SumCalculator struct {
}

func NewSumCalculator() *SumCalculator {
	return &SumCalculator{}
}

// SumNumbers is a function which Accepts arbitrary JSON document as payload, which can contain a variety of things:
// arrays [1,2,3,4], arbitrary objects {"a":1, "b":2, "c":3}, numbers, and strings.
// The code should find all of the numbers throughout the request and add them together.
func (s *SumCalculator) SumNumbers(v interface{}) int {
	sum := 0
	s.recursiveSum(v, &sum)

	return sum
}

func (m *SumCalculator) recursiveSum(src interface{}, sum *int) {
	switch part := src.(type) {
	// corner case, we found number
	case float64:
		*sum += int(part)
	// deal with arrays
	case []interface{}:
		for _, v := range part {
			m.recursiveSum(v, sum)
		}
	// with objects
	case map[string]interface{}:
		for _, v := range part {
			m.recursiveSum(v, sum)
		}
	}
}

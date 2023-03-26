package parser

import "fmt"

func compareNumbers[V int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64](l, r V, comparator Operator) (isMatch bool, err error) {
	switch comparator {
	case Any:
		isMatch = true
	case Equal:
		isMatch = l == r
	case NotEqual:
		isMatch = l != r
	case GreaterThan:
		isMatch = l > r
	case LessThan:
		isMatch = l < r
	case GreaterThanOrEqual:
		isMatch = l >= r
	case LessThanOrEqual:
		isMatch = l <= r
	default:
		return false, fmt.Errorf("unknown comparator %d", comparator)
	}
	return isMatch, nil
}

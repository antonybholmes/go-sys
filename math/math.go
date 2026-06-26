package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func AbsDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func AbsDiffint(x, y int) int {
	if x < y {
		return y - x
	}

	return x - y
}

// func IntMin(x, y int) int {
// 	if x < y {
// 		return x
// 	}
// 	return y
// }

// func IntMax(x, y int) int {
// 	if x > y {
// 		return x
// 	}
// 	return y
// }

// func intMin(x, y int) int {
// 	if x < y {
// 		return x
// 	}
// 	return y
// }

// func intMax(x, y int) int {
// 	if x > y {
// 		return x
// 	}
// 	return y
// }

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func LnFactorial(n int) float64 {
	// for property 0! = 1 since exp(0) == 1
	if n == 0 {
		return 0
	}

	var ret float64 = 0

	for i := range n {
		ret += math.Log(float64(i + 1))
	}

	return ret
}

func Factorial(n int) int64 {
	if n == 0 {
		return 1
	}

	return int64(math.Round(math.Exp(LnFactorial(n))))
}

// Returns the max value of an array of floats
func MaxFloat64Array(values *[]float64) float64 {

	if len(*values) == 0 {
		return 0 // Or handle the empty array case appropriately
	}

	max := 0.0

	for _, num := range *values {
		max = math.Max(num, max)
	}

	return max

}

// Returns the max value in an array of ints.
func MaxintArray(values *[]int) int {

	if len(*values) == 0 {
		return 0 // Or handle the empty array case appropriately
	}

	max := int(0)

	for _, num := range *values {
		max = Max(num, max)
	}

	return max

}

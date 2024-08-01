package sys

import (
	"cmp"
	"slices"
)

// func main() {
// 	fmt.Println("Hello, playground")

// 	example := []float64{1, 25, 3, 5, 4}
// 	fmt.Println("pre :", example)
// 	ia := ArgsortNew(example)
// 	fmt.Println("post:", example)
// 	fmt.Println("ia  :", ia)
// 	sortedExample := []float64{}
// 	for i := range ia {
// 		sortedExample = append(sortedExample, example[ia[i]])
// 	}
// 	fmt.Println("sortedExample:", sortedExample)
// }

// simplified and updated to use slices from https://github.com/mkmik/argsort/tree/main

// argsort, like in Numpy, it returns an array of indexes into an array. Note
// that the gonum version of argsort reorders the original array and returns
// indexes to reconstruct the original order.
type argsort[T cmp.Ordered] struct {
	v   T   // value to be sorted
	idx int // keep track of the original index
}

// ArgsortNew allocates and returns an array of indexes into the source float
// array.
func Argsort[T cmp.Ordered](src []T) []int {
	argsortList := make([]argsort[T], len(src))
	for i := range src {
		argsortList[i] = argsort[T]{v: src[i], idx: i}
	}

	slices.SortStableFunc(argsortList, func(a, b argsort[T]) int {
		return cmp.Compare(a.v, b.v)
	})

	indices := make([]int, len(src))
	for i := range src {
		indices[i] = argsortList[i].idx
	}

	return indices
}

func Reorder[T any](slice []T, order []int) []T {
	reordered := make([]T, len(slice))

	for i, index := range order {
		if index < 0 || index >= len(slice) {
			panic("Invalid index in order array")
		}
		reordered[i] = slice[index]
	}

	return reordered
}

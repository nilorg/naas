package tools

import (
	mapset "github.com/deckarep/golang-set"
)

// SliceIsEqual 是相等
func SliceIsEqual(input, source []string) bool {
	inputSet := mapset.NewSet()
	for _, i := range input {
		inputSet.Add(i)
	}
	sourceSet := mapset.NewSet()
	for _, s := range source {
		sourceSet.Add(s)
	}
	return inputSet.Equal(sourceSet)
}

// SliceIsSubset 是包含
func SliceIsSubset(input, source []string) bool {
	inputSet := mapset.NewSet()
	for _, i := range input {
		inputSet.Add(i)
	}
	sourceSet := mapset.NewSet()
	for _, s := range source {
		sourceSet.Add(s)
	}
	return inputSet.IsSubset(sourceSet)
}

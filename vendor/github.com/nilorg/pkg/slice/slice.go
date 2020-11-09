package slice

import (
	mapset "github.com/deckarep/golang-set"
)

// IsEqual 是相等
func IsEqual(input, source []string) bool {
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

// IsSubset 是包含
func IsSubset(input, source []string) bool {
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

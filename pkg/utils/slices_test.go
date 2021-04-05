package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64SliceEquals_Equals(t *testing.T) {
	s1 := []int64{0, 1, 2, 3}
	s2 := []int64{0, 1, 2, 3}
	expectedValue := true
	actualValue := Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInt64SliceEquals_NotEquals(t *testing.T) {
	s1 := []int64{0, 1, 2, 3}
	s2 := []int64{0, 1, 3, 2}
	expectedValue := false
	actualValue := Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = []int64{0}
	expectedValue = false
	actualValue = Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = []int64{}
	expectedValue = false
	actualValue = Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInt64SliceEquals_NilOrEmpty(t *testing.T) {
	var s1 []int64 = nil
	s2 := []int64{0, 1, 2, 3}
	expectedValue := false
	actualValue := Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = nil
	expectedValue = false
	actualValue = Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = nil
	s2 = nil
	expectedValue = true
	actualValue = Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{}
	s2 = []int64{}
	expectedValue = true
	actualValue = Int64SliceEquals(s1, s2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInt64SliceMinus(t *testing.T) {
	s1 := []int64{0, 1, 2, 3}
	s2 := []int64{0, 1, 2, 3}
	expectedValue := []int64{}
	actualValue := Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = []int64{0, 2}
	expectedValue = []int64{1, 3}
	actualValue = Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = []int64{}
	expectedValue = []int64{0, 1, 2, 3}
	actualValue = Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{}
	s2 = []int64{0, 1, 2, 3}
	expectedValue = []int64{}
	actualValue = Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = []int64{1, 2, 3, 4}
	expectedValue = []int64{0}
	actualValue = Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInt64SliceMinus_NilOrEmpty(t *testing.T) {
	var s1 []int64 = nil
	s2 := []int64{0, 1, 2, 3}
	var expectedValue []int64 = nil
	actualValue := Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = []int64{0, 1, 2, 3}
	s2 = nil
	expectedValue = []int64{0, 1, 2, 3}
	actualValue = Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)

	s1 = nil
	s2 = nil
	expectedValue = nil
	actualValue = Int64SliceMinus(s1, s2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestToUniqueInt64Slice(t *testing.T) {
	arr := []int64{0, 1, 2, 3, 2, 4, 0}
	expectedValue := []int64{0, 1, 2, 3, 4}
	actualValue := ToUniqueInt64Slice(arr)
	assert.Equal(t, expectedValue, actualValue)
}

func TestToUniqueInt64Slice_NilOrEmpty(t *testing.T) {
	var arr []int64 = nil
	expectedValue := []int64{}
	actualValue := ToUniqueInt64Slice(arr)
	assert.Equal(t, expectedValue, actualValue)

	arr = []int64{}
	expectedValue = []int64{}
	actualValue = ToUniqueInt64Slice(arr)
	assert.Equal(t, expectedValue, actualValue)
}

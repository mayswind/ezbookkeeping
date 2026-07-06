package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeMaps(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"c": 3, "d": 4}
	expectedValue := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	actualValue := MergeMaps(m1, m2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestMergeMaps_Override(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	expectedValue := map[string]int{"a": 1, "b": 3, "c": 4}
	actualValue := MergeMaps(m1, m2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestMergeMaps_NilOrEmpty(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	var m2 map[string]int = nil
	expectedValue := map[string]int{"a": 1, "b": 2}
	actualValue := MergeMaps(m1, m2)
	assert.Equal(t, expectedValue, actualValue)

	m1 = nil
	m2 = map[string]int{"c": 3, "d": 4}
	expectedValue = map[string]int{"c": 3, "d": 4}
	actualValue = MergeMaps(m1, m2)
	assert.Equal(t, expectedValue, actualValue)

	m1 = nil
	m2 = nil
	expectedValue = map[string]int{}
	actualValue = MergeMaps(m1, m2)
	assert.Equal(t, expectedValue, actualValue)

	m1 = map[string]int{}
	m2 = map[string]int{}
	expectedValue = map[string]int{}
	actualValue = MergeMaps(m1, m2)
	assert.Equal(t, expectedValue, actualValue)
}

func TestMergeMaps_MultipleMaps(t *testing.T) {
	m1 := map[string]int{"a": 1}
	m2 := map[string]int{"b": 2}
	m3 := map[string]int{"c": 3}
	expectedValue := map[string]int{"a": 1, "b": 2, "c": 3}
	actualValue := MergeMaps(m1, m2, m3)
	assert.Equal(t, expectedValue, actualValue)

	m1 = map[string]int{"a": 1}
	m2 = map[string]int{"a": 2}
	m3 = map[string]int{"a": 3}
	expectedValue = map[string]int{"a": 3}
	actualValue = MergeMaps(m1, m2, m3)
	assert.Equal(t, expectedValue, actualValue)
}

package core

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestNormalTokenAdditionalData_MarshalBinary(t *testing.T) {
	data := &NormalTokenAdditionalData{Scope: NormalTokenScopeAll}
	var marshaler encoding.BinaryMarshaler = data
	expectedValue := []byte{1}
	actualValue, err := marshaler.MarshalBinary()

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_MarshalBinary_OtherScopes(t *testing.T) {
	data := &NormalTokenAdditionalData{}
	actualValue, err := data.MarshalBinary()
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{0}, actualValue)

	data.Scope = NormalTokenScope(255)
	actualValue, err = data.MarshalBinary()
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{255}, actualValue)
}

func TestNormalTokenAdditionalData_UnmarshalBinary(t *testing.T) {
	data := &NormalTokenAdditionalData{}
	var unmarshaler encoding.BinaryUnmarshaler = data
	err := unmarshaler.UnmarshalBinary([]byte{1})

	assert.Equal(t, nil, err)
	assert.Equal(t, NormalTokenScopeAll, data.Scope)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_InvalidLength(t *testing.T) {
	data := &NormalTokenAdditionalData{}
	err := data.UnmarshalBinary(nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, NormalTokenScope(0), data.Scope)

	err = data.UnmarshalBinary([]byte{})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, NormalTokenScope(0), data.Scope)

	err = data.UnmarshalBinary([]byte{1, 0})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, NormalTokenScope(0), data.Scope)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_OtherScopes(t *testing.T) {
	data := &NormalTokenAdditionalData{}
	err := data.UnmarshalBinary([]byte{0})
	assert.Equal(t, nil, err)
	assert.Equal(t, NormalTokenScope(0), data.Scope)

	err = data.UnmarshalBinary([]byte{255})
	assert.Equal(t, nil, err)
	assert.Equal(t, NormalTokenScope(255), data.Scope)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_InvalidDataKeepsScope(t *testing.T) {
	data := &NormalTokenAdditionalData{Scope: NormalTokenScopeAll}
	err := data.UnmarshalBinary(nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, NormalTokenScopeAll, data.Scope)
}

func TestNormalTokenAdditionalData_BinaryRoundTrip(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{Scope: NormalTokenScopeAll}
	encodedData, err := expectedValue.MarshalBinary()
	assert.Equal(t, nil, err)

	actualValue := &NormalTokenAdditionalData{}
	err = actualValue.UnmarshalBinary(encodedData)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue.Scope = NormalTokenScope(255)
	encodedData, err = expectedValue.MarshalBinary()
	assert.Equal(t, nil, err)

	err = actualValue.UnmarshalBinary(encodedData)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

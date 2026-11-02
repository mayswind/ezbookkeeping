package core

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestNormalTokenAdditionalData_MarshalBinary(t *testing.T) {
	data := &NormalTokenAdditionalData{
		Scope: NormalTokenScopeAll,
	}
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
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, actualValue)

	data.Scope = NormalTokenScope(255)
	actualValue, err = data.MarshalBinary()
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{255, 0, 0, 0, 0, 0, 0, 0, 0}, actualValue)
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
	err := data.UnmarshalBinary([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, nil, err)
	assert.Equal(t, NormalTokenScope(0), data.Scope)

	err = data.UnmarshalBinary([]byte{255, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, nil, err)
	assert.Equal(t, NormalTokenScope(255), data.Scope)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_InvalidDataKeepsScope(t *testing.T) {
	data := &NormalTokenAdditionalData{
		Scope: NormalTokenScopeAll,
	}
	err := data.UnmarshalBinary(nil)
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, NormalTokenScopeAll, data.Scope)
}

func TestNormalTokenAdditionalData_BinaryRoundTrip(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{
		Scope: NormalTokenScopeAll,
	}
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

func TestNormalTokenAdditionalData_MarshalBinary_UserAvatar(t *testing.T) {
	data := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}

	expectedValue := []byte{byte(NormalTokenScopeUserAvatar), 1, 2, 3, 4, 5, 6, 7, 8}
	actualValue, err := data.MarshalBinary()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_UserAvatar(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}
	actualValue := &NormalTokenAdditionalData{}

	err := actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeUserAvatar), 1, 2, 3, 4, 5, 6, 7, 8})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_UserAvatarInvalidLength(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}
	actualValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}

	err := actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeUserAvatar)})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	err = actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeUserAvatar), 1, 2, 3, 4, 5, 6, 7})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	err = actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeUserAvatar), 1, 2, 3, 4, 5, 6, 7, 8, 0})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_MarshalBinary_TransactionPicture(t *testing.T) {
	data := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}

	expectedValue := []byte{byte(NormalTokenScopeTransactionPicture), 1, 2, 3, 4, 5, 6, 7, 8}
	actualValue, err := data.MarshalBinary()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_TransactionPicture(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}
	actualValue := &NormalTokenAdditionalData{}

	err := actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeTransactionPicture), 1, 2, 3, 4, 5, 6, 7, 8})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_TransactionPictureInvalidLength(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}
	actualValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeTransactionPicture,
		DerivedTokenExpiresAt: 0x0102030405060708,
	}

	err := actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeTransactionPicture)})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	err = actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeTransactionPicture), 1, 2, 3, 4, 5, 6, 7})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)

	err = actualValue.UnmarshalBinary([]byte{byte(NormalTokenScopeTransactionPicture), 1, 2, 3, 4, 5, 6, 7, 8, 0})
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_MarshalBinary_AllOmitsExpiration(t *testing.T) {
	data := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeAll,
		DerivedTokenExpiresAt: 123,
	}

	expectedValue := []byte{1}
	actualValue, err := data.MarshalBinary()

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestNormalTokenAdditionalData_UnmarshalBinary_AllClearsExpiration(t *testing.T) {
	expectedValue := &NormalTokenAdditionalData{
		Scope: NormalTokenScopeAll,
	}
	actualValue := &NormalTokenAdditionalData{
		Scope:                 NormalTokenScopeUserAvatar,
		DerivedTokenExpiresAt: 123,
	}

	err := actualValue.UnmarshalBinary([]byte{1})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

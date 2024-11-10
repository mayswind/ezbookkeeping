package errs

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorError(t *testing.T) {
	err := New(CATEGORY_SYSTEM, 12, 34, http.StatusInternalServerError, "error message")
	assert.EqualError(t, err, "error message")
}

func TestErrorCode(t *testing.T) {
	err := New(CATEGORY_SYSTEM, 12, 34, http.StatusInternalServerError, "error message")
	assert.Equal(t, int32(112034), err.Code())
}

func TestMultiError(t *testing.T) {
	err1 := errors.New("error1 message")
	err2 := errors.New("error2 message")
	err := NewMultiErrorOrNil(err1, err2)
	assert.EqualError(t, err, "multi errors: error1 message, error2 message")
}

func TestNewMultiErrorOrNilWithOnlyOneParameter(t *testing.T) {
	err1 := errors.New("error1 message")
	err := NewMultiErrorOrNil(err1)
	assert.Equal(t, err1, err)
	assert.EqualError(t, err, "error1 message")
}

func TestNewMultiErrorOrNilWithoutOneParameter(t *testing.T) {
	err := NewMultiErrorOrNil()
	assert.Nil(t, err)
}

func TestOr(t *testing.T) {
	err1 := errors.New("test error")
	err2 := New(CATEGORY_SYSTEM, 12, 34, http.StatusInternalServerError, "test custom error")
	err := Or(err1, err2)
	assert.Equal(t, err2, err)
	assert.EqualError(t, err, "test custom error")

	err1 = New(CATEGORY_SYSTEM, 12, 34, http.StatusInternalServerError, "test custom error1")
	err2 = New(CATEGORY_SYSTEM, 23, 45, http.StatusInternalServerError, "test custom error2")
	err = Or(err1, err2)
	assert.Equal(t, err1, err)
	assert.EqualError(t, err, "test custom error1")
}

func TestIsCustomError(t *testing.T) {
	err1 := errors.New("test error")
	err2 := New(CATEGORY_SYSTEM, 12, 34, http.StatusInternalServerError, "test custom error")
	assert.False(t, IsCustomError(err1))
	assert.True(t, IsCustomError(err2))
}

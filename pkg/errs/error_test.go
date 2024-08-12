package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError(t *testing.T) {
	err1 := errors.New("error1 message")
	err2 := errors.New("error2 message")
	err := NewMultiErrorOrNil(err1, err2)
	assert.Equal(t, "multi errors: error1 message, error2 message", err.Error())
}

func TestNewMultiErrorOrNilWithOnlyOneParamerter(t *testing.T) {
	err1 := errors.New("error1 message")
	err := NewMultiErrorOrNil(err1)
	assert.Equal(t, err1, err)
}

func TestNewMultiErrorOrNilWithoutOneParamerter(t *testing.T) {
	err := NewMultiErrorOrNil()
	assert.Nil(t, err)
}

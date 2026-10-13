package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportError_DoesNotMutateSharedErrors(t *testing.T) {
	actualErr := NewImportError(ErrAmountInvalid, "record#1")
	assert.Nil(t, ErrAmountInvalid.Context)
	assert.Equal(t, ErrAmountInvalid.Code(), actualErr.Code())
	assert.Equal(t, ErrAmountInvalid.HttpStatusCode, actualErr.HttpStatusCode)
	assert.Equal(t, ErrAmountInvalid.Message, actualErr.Message)

	expectedContext := &ImportErrorContext{RecordId: "record#1"}
	assert.Equal(t, expectedContext, actualErr.Context)

	actualErr = NewImportError(ErrAmountInvalid, "record#2")
	assert.Nil(t, ErrAmountInvalid.Context)

	expectedContext = &ImportErrorContext{RecordId: "record#2"}
	assert.Equal(t, expectedContext, actualErr.Context)

	anotherActualErr := NewImportError(actualErr, "record#3")
	assert.Same(t, actualErr, anotherActualErr)

	assert.Same(t, ErrAmountInvalid, NewImportError(ErrAmountInvalid, ""))
}

func TestImportError_KeepsOriginalValueWhenAddingRecordId(t *testing.T) {
	actualErr := NewImportError(ErrAmountInvalid, "", "1,2..3")
	assert.Nil(t, ErrAmountInvalid.Context)

	expectedInnerErrContext := &ImportErrorContext{
		OriginalData: "1,2..3",
	}
	assert.Equal(t, expectedInnerErrContext, actualErr.Context)

	actualErr = NewImportError(actualErr, "line#1", []string{"fallback row"})
	expectedContext := &ImportErrorContext{
		RecordId:     "line#1",
		OriginalData: "1,2..3",
	}
	assert.Equal(t, expectedContext, actualErr.Context)

	expectedErr := NewImportError(ErrAmountInvalid, "line#1", "1,2..3")
	assert.Equal(t, expectedErr, actualErr)

	anotherActualErr := NewImportError(actualErr, "line#2", []string{"fallback row"})
	assert.Equal(t, &ImportErrorContext{RecordId: "line#1", OriginalData: "1,2..3"}, anotherActualErr.Context)
}

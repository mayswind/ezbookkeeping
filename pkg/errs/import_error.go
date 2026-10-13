package errs

// ImportErrorContext represents the context of an import error
type ImportErrorContext struct {
	RecordId     string `json:"recordId,omitempty"`
	OriginalData any    `json:"originalData,omitempty"`
}

// NewImportError returns a new import error with the given record ID and original data
func NewImportError(err *Error, recordId string, originalData ...any) *Error {
	context := ImportErrorContext{
		RecordId: recordId,
	}

	if existing, ok := err.Context.(*ImportErrorContext); ok {
		context = *existing

		if context.RecordId == "" {
			context.RecordId = recordId
		}

		if context.RecordId == existing.RecordId && (context.OriginalData != nil || len(originalData) == 0) {
			return err
		}
	}

	if context.OriginalData == nil && len(originalData) > 0 {
		context.OriginalData = originalData[0]
	}

	if context.RecordId == "" && context.OriginalData == nil {
		return err
	}

	return NewErrorWithContext(err, &context)
}

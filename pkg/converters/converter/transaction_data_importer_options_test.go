package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImporterOptions(t *testing.T) {
	actualValue := ParseImporterOptions("payeeAsTag,memberAsTag")
	expectedValue := TransactionDataImporterOptions{
		payeeAsTag:    true,
		memberAsTag:   true,
		projectAsTag:  false,
		merchantAsTag: false,
	}
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, true, actualValue.IsPayeeAsTag())
	assert.Equal(t, true, actualValue.IsMemberAsTag())
	assert.Equal(t, false, actualValue.IsProjectAsTag())
	assert.Equal(t, false, actualValue.IsMerchantAsTag())

	actualValue = ParseImporterOptions("")
	expectedValue = TransactionDataImporterOptions{
		payeeAsTag:    false,
		memberAsTag:   false,
		projectAsTag:  false,
		merchantAsTag: false,
	}
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, false, actualValue.IsPayeeAsTag())
	assert.Equal(t, false, actualValue.IsMemberAsTag())
	assert.Equal(t, false, actualValue.IsProjectAsTag())
	assert.Equal(t, false, actualValue.IsMerchantAsTag())
}

func TestParseImporterOptions_WithAllOptions(t *testing.T) {
	actualValue := ParseImporterOptions("payeeAsTag,payeeAsDescription,memberAsTag,projectAsTag,merchantAsTag")
	expectedValue := TransactionDataImporterOptions{
		payeeAsTag:         true,
		payeeAsDescription: true,
		memberAsTag:        true,
		projectAsTag:       true,
		merchantAsTag:      true,
	}
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, true, actualValue.IsPayeeAsTag())
	assert.Equal(t, true, actualValue.IsPayeeAsDescription())
	assert.Equal(t, true, actualValue.IsMemberAsTag())
	assert.Equal(t, true, actualValue.IsProjectAsTag())
	assert.Equal(t, true, actualValue.IsMerchantAsTag())
}

func TestParseImporterOptions_WithInvalidOptions(t *testing.T) {
	actualValue := ParseImporterOptions("invalidOption,payeeAsTag,memberAsTag")
	expectedValue := TransactionDataImporterOptions{
		payeeAsTag:    true,
		memberAsTag:   true,
		projectAsTag:  false,
		merchantAsTag: false,
	}
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, true, actualValue.IsPayeeAsTag())
	assert.Equal(t, true, actualValue.IsMemberAsTag())
	assert.Equal(t, false, actualValue.IsProjectAsTag())
	assert.Equal(t, false, actualValue.IsMerchantAsTag())

	actualValue = ParseImporterOptions("invalidOption")
	expectedValue = TransactionDataImporterOptions{
		payeeAsTag:    false,
		memberAsTag:   false,
		projectAsTag:  false,
		merchantAsTag: false,
	}
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, false, actualValue.IsPayeeAsTag())
	assert.Equal(t, false, actualValue.IsMemberAsTag())
	assert.Equal(t, false, actualValue.IsProjectAsTag())
	assert.Equal(t, false, actualValue.IsMerchantAsTag())
}

func TestParseImporterOptions_Clone(t *testing.T) {
	original := TransactionDataImporterOptions{
		payeeAsTag:         true,
		payeeAsDescription: false,
		memberAsTag:        false,
		projectAsTag:       true,
		merchantAsTag:      false,
	}

	cloned := original.Clone()
	assert.Equal(t, original, cloned)

	// Modify cloned options and verify original options are not affected
	cloned.payeeAsTag = false
	cloned.payeeAsDescription = true
	cloned.memberAsTag = true

	assert.Equal(t, true, original.payeeAsTag)
	assert.Equal(t, false, original.payeeAsDescription)
	assert.Equal(t, false, original.memberAsTag)
	assert.Equal(t, true, original.projectAsTag)
	assert.Equal(t, false, original.merchantAsTag)

	assert.Equal(t, false, cloned.payeeAsTag)
	assert.Equal(t, true, cloned.payeeAsDescription)
	assert.Equal(t, true, cloned.memberAsTag)
	assert.Equal(t, true, cloned.projectAsTag)
	assert.Equal(t, false, cloned.merchantAsTag)
}

package converter

import "strings"

// TransactionDataImporterOptions defines the options for transaction data importer
type TransactionDataImporterOptions struct {
	payeeAsTag         bool
	payeeAsDescription bool
	memberAsTag        bool
	projectAsTag       bool
	merchantAsTag      bool
}

// DefaultImporterOptions provides the default options for transaction data importer
var DefaultImporterOptions = TransactionDataImporterOptions{
	payeeAsTag:         false,
	payeeAsDescription: false,
	memberAsTag:        false,
	projectAsTag:       false,
	merchantAsTag:      false,
}

// IsPayeeAsTag returns whether to import payee as tag
func (o TransactionDataImporterOptions) IsPayeeAsTag() bool {
	return o.payeeAsTag
}

// IsPayeeAsDescription returns whether to import payee as description
func (o TransactionDataImporterOptions) IsPayeeAsDescription() bool {
	return o.payeeAsDescription
}

// IsMemberAsTag returns whether to import member as tag
func (o TransactionDataImporterOptions) IsMemberAsTag() bool {
	return o.memberAsTag
}

// IsProjectAsTag returns whether to import project as tag
func (o TransactionDataImporterOptions) IsProjectAsTag() bool {
	return o.projectAsTag
}

// IsMerchantAsTag returns whether to import merchant as tag
func (o TransactionDataImporterOptions) IsMerchantAsTag() bool {
	return o.merchantAsTag
}

// WithPayeeAsTag sets the option to import payee as tag
func (o TransactionDataImporterOptions) WithPayeeAsTag() TransactionDataImporterOptions {
	cloned := o.Clone()
	cloned.payeeAsTag = true
	return cloned
}

// WithPayeeAsDescription sets the option to import payee as description
func (o TransactionDataImporterOptions) WithPayeeAsDescription() TransactionDataImporterOptions {
	cloned := o.Clone()
	cloned.payeeAsDescription = true
	return cloned
}

// WithMemberAsTag sets the option to import member as tag
func (o TransactionDataImporterOptions) WithMemberAsTag() TransactionDataImporterOptions {
	cloned := o.Clone()
	cloned.memberAsTag = true
	return cloned
}

// WithProjectAsTag sets the option to import project as tag
func (o TransactionDataImporterOptions) WithProjectAsTag() TransactionDataImporterOptions {
	cloned := o.Clone()
	cloned.projectAsTag = true
	return cloned
}

// WithMerchantAsTag sets the option to import merchant as tag
func (o TransactionDataImporterOptions) WithMerchantAsTag() TransactionDataImporterOptions {
	cloned := o.Clone()
	cloned.merchantAsTag = true
	return cloned
}

// Clone creates a copy of the options instance
func (o TransactionDataImporterOptions) Clone() TransactionDataImporterOptions {
	return TransactionDataImporterOptions{
		payeeAsTag:         o.payeeAsTag,
		payeeAsDescription: o.payeeAsDescription,
		memberAsTag:        o.memberAsTag,
		projectAsTag:       o.projectAsTag,
		merchantAsTag:      o.merchantAsTag,
	}
}

// ParseImporterOptions parses the textual options to the instance
func ParseImporterOptions(s string) TransactionDataImporterOptions {
	options := TransactionDataImporterOptions{}

	if s == "" {
		return options
	}

	for _, option := range strings.Split(s, ",") {
		switch option {
		case "payeeAsTag":
			options.payeeAsTag = true
		case "payeeAsDescription":
			options.payeeAsDescription = true
		case "memberAsTag":
			options.memberAsTag = true
		case "projectAsTag":
			options.projectAsTag = true
		case "merchantAsTag":
			options.merchantAsTag = true
		}
	}

	return options
}

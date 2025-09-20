package models

// RecognizedReceiptImageResponse represents a view-object of recognized receipt image response
type RecognizedReceiptImageResponse struct {
	Type                 TransactionType `json:"type"`
	Time                 int64           `json:"time,omitempty"`
	CategoryId           int64           `json:"categoryId,string,omitempty"`
	SourceAccountId      int64           `json:"sourceAccountId,string,omitempty"`
	DestinationAccountId int64           `json:"destinationAccountId,string,omitempty"`
	SourceAmount         int64           `json:"sourceAmount,omitempty"`
	DestinationAmount    int64           `json:"destinationAmount,omitempty"`
	TagIds               []string        `json:"tagIds,omitempty"`
	Comment              string          `json:"comment,omitempty"`
}

// RecognizedReceiptImageResult represents the result of recognized receipt image
type RecognizedReceiptImageResult struct {
	Type                   string   `json:"type,omitempty" jsonschema:"enum=income,enum=expense,enum=transfer" jsonschema_description:"Transaction type (income, expense, transfer)"`
	Time                   string   `json:"time" jsonschema:"format=date-time" jsonschema_description:"Transaction time in long date time format (YYYY-MM-DD HH:mm:ss, e.g. 2023-01-01 12:00:00)"`
	Amount                 string   `json:"amount,omitempty" jsonschema_description:"Transaction amount"`
	AccountName            string   `json:"account,omitempty" jsonschema_description:"Account name for the transaction"`
	CategoryName           string   `json:"category,omitempty" jsonschema_description:"Category name for the transaction"`
	TagNames               []string `json:"tags,omitempty" jsonschema_description:"List of tags associated with the transaction (maximum 10 tags allowed)"`
	Description            string   `json:"description,omitempty" jsonschema_description:"Transaction description"`
	DestinationAmount      string   `json:"destination_amount,omitempty" jsonschema_description:"Destination amount for transfer transactions"`
	DestinationAccountName string   `json:"destination_account,omitempty" jsonschema_description:"Destination account name for transfer transactions"`
}

package camt

import "encoding/xml"

type camtCreditDebitIndicator string

const (
	CAMT_INDICATOR_CREDIT camtCreditDebitIndicator = "CRDT"
	CAMT_INDICATOR_DEBIT  camtCreditDebitIndicator = "DBIT"
)

type camt052File struct {
	XMLName                     xml.Name                         `xml:"Document"`
	BankToCustomerAccountReport *camtBankToCustomerAccountReport `xml:"BkToCstmrAcctRpt"`
}

type camt053File struct {
	XMLName                 xml.Name                     `xml:"Document"`
	BankToCustomerStatement *camtBankToCustomerStatement `xml:"BkToCstmrStmt"`
}

type camtBankToCustomerAccountReport struct {
	Statements []*camtStatement `xml:"Rpt"`
}

type camtBankToCustomerStatement struct {
	Statements []*camtStatement `xml:"Stmt"`
}

type camtStatement struct {
	Account *camtAccount `xml:"Acct"`
	Entries []*camtEntry `xml:"Ntry"`
}

type camtAccount struct {
	IBAN                string `xml:"Id>IBAN"`
	OtherIdentification string `xml:"Id>Othr>Id"`
	Currency            string `xml:"Ccy"`
}

type camtEntry struct {
	Amount                     *camtAmount              `xml:"Amt"`
	CreditDebitIndicator       camtCreditDebitIndicator `xml:"CdtDbtInd"`
	BookingDate                *camtDate                `xml:"BookgDt"`
	EntryDetails               *camtEntryDetails        `xml:"NtryDtls"`
	AdditionalEntryInformation string                   `xml:"AddtlNtryInf"`
}

type camtAmount struct {
	Value    string `xml:",chardata"`
	Currency string `xml:"Ccy,attr"`
}

type camtDate struct {
	Date     string `xml:"Dt"`
	DateTime string `xml:"DtTm"`
}

type camtEntryDetails struct {
	TransactionDetails []*camtTransactionDetails `xml:"TxDtls"`
}

type camtTransactionDetails struct {
	AmountDetails                    *camtAmountDetails         `xml:"AmtDtls"`
	RemittanceInformation            *camtRemittanceInformation `xml:"RmtInf"`
	AdditionalTransactionInformation string                     `xml:"AddtlTxInf"`
}

type camtAmountDetails struct {
	InstructedAmount  *camtAmount `xml:"InstdAmt>Amt"`
	TransactionAmount *camtAmount `xml:"TxAmt>Amt"`
}

type camtRemittanceInformation struct {
	Unstructured []string `xml:"Ustrd"`
}

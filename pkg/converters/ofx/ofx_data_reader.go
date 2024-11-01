package ofx

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const ofxUnicodeEncoding = "unicode"
const ofxUSAsciiEncoding = "usascii"
const ofx1SGMLDataFormat = "OFXSGML"

const ofxDataElementName = "OFX"

const ofxBankMessageResponseV1ElementName = "BANKMSGSRSV1"
const ofxCreditCardMessageResponseV1ElementName = "CREDITCARDMSGSRSV1"

const ofxBankStatementTransactionResponseElementName = "STMTTRNRS"
const ofxCreditCardStatementTransactionResponseElementName = "CCSTMTTRNRS"

const ofxBankStatementResponseElementName = "STMTRS"
const ofxCreditCardStatementResponseElementName = "CCSTMTRS"

const ofxBankStatementResponseDefaultCurrencyName = "CURDEF"
const ofxBankStatementResponseBankAccountFromName = "BANKACCTFROM"
const ofxBankStatementResponseBankTransactionListName = "BANKTRANLIST"

const ofxCreditCardStatementResponseDefaultCurrencyName = "CURDEF"
const ofxCreditCardStatementResponseCreditCardAccountFromName = "CCACCTFROM"
const ofxCreditCardStatementResponseCreditCardTransactionListName = "BANKTRANLIST"

const ofxBankTransactionListStartDateName = "DTSTART"
const ofxBankTransactionListEndDateName = "DTEND"
const ofxBankTransactionListStatementTransactionsName = "STMTTRN"

const ofxCreditCardTransactionListStartDateName = "DTSTART"
const ofxCreditCardTransactionListEndDateName = "DTEND"
const ofxCreditCardTransactionListStatementTransactionsName = "STMTTRN"

const ofxBankAccountBankIdName = "BANKID"
const ofxBankAccountBranchIdName = "BRANCHID"
const ofxBankAccountAccountIdName = "ACCTID"
const ofxBankAccountAccountTypeName = "ACCTTYPE"
const ofxBankAccountAccountKeyName = "ACCTKEY"

const ofxCreditCardAccountAccountIdName = "ACCTID"
const ofxCreditCardAccountAccountKeyName = "ACCTKEY"

const ofxTransactionTransactionIdName = "FITID"
const ofxTransactionTransactionTypeName = "TRNTYPE"
const ofxTransactionPostedDateName = "DTPOSTED"
const ofxTransactionAmountName = "TRNAMT"
const ofxTransactionNameName = "NAME"
const ofxTransactionMemoName = "MEMO"
const ofxTransactionCurrencyName = "CURRENCY"
const ofxTransactionOriginalCurrencyName = "ORIGCURRENCY"
const ofxTransactionPayeeName = "PAYEE"
const ofxTransactionBankAccountToName = "BANKACCTTO"
const ofxTransactionCreditCardAccountToName = "CCACCTTO"

const ofxPayeeNameName = "NAME"
const ofxPayeeAddress1Name = "ADDR1"
const ofxPayeeAddress2Name = "ADDR2"
const ofxPayeeAddress3Name = "ADDR3"
const ofxPayeeCityName = "CITY"
const ofxPayeeStateName = "STATE"
const ofxPayeePostalCodeName = "POSTALCODE"
const ofxPayeeCountryName = "COUNTRY"
const ofxPayeePhoneName = "PHONE"

var ofxBankStatementResponseChildrenNames = map[string]bool{
	ofxBankStatementResponseDefaultCurrencyName: true,
}

var ofxCreditCardStatementResponseChildrenNames = map[string]bool{
	ofxCreditCardStatementResponseDefaultCurrencyName: true,
}

var ofxBankTransactionListChildrenNames = map[string]bool{
	ofxBankTransactionListStartDateName: true,
	ofxBankTransactionListEndDateName:   true,
}

var ofxCreditCardTransactionListChildrenNames = map[string]bool{
	ofxCreditCardTransactionListStartDateName: true,
	ofxCreditCardTransactionListEndDateName:   true,
}

var ofxBankAccountChildrenNames = map[string]bool{
	ofxBankAccountBankIdName:      true,
	ofxBankAccountBranchIdName:    true,
	ofxBankAccountAccountIdName:   true,
	ofxBankAccountAccountTypeName: true,
	ofxBankAccountAccountKeyName:  true,
}

var ofxCreditCardAccountChildrenNames = map[string]bool{
	ofxCreditCardAccountAccountIdName:  true,
	ofxCreditCardAccountAccountKeyName: true,
}

var ofxTransactionChildrenNames = map[string]bool{
	ofxTransactionTransactionIdName:    true,
	ofxTransactionTransactionTypeName:  true,
	ofxTransactionPostedDateName:       true,
	ofxTransactionAmountName:           true,
	ofxTransactionNameName:             true,
	ofxTransactionMemoName:             true,
	ofxTransactionCurrencyName:         true,
	ofxTransactionOriginalCurrencyName: true,
}

var ofxPayeeChildrenNames = map[string]bool{
	ofxPayeeNameName:       true,
	ofxPayeeAddress1Name:   true,
	ofxPayeeAddress2Name:   true,
	ofxPayeeAddress3Name:   true,
	ofxPayeeCityName:       true,
	ofxPayeeStateName:      true,
	ofxPayeePostalCodeName: true,
	ofxPayeeCountryName:    true,
	ofxPayeePhoneName:      true,
}

var ofx2HeaderPattern = regexp.MustCompile("<\\?OFX( +[A-Z]+=\"[^=]*\")* *\\?>")
var ofx2HeaderAttributePattern = regexp.MustCompile(" +([A-Z]+)=\"([^=]*)\"")

// ofxFileReader defines the structure of open financial exchange (ofx) file reader
type ofxFileReader struct {
	fileHeader *ofxFileHeader
	xmlDecoder *xml.Decoder
}

// read returns the imported open financial exchange (ofx) file
func (r *ofxFileReader) read(ctx core.Context) (*ofxFile, error) {
	var file *ofxFile
	strictMode := true

	if r.fileHeader != nil && r.fileHeader.OFXDeclarationVersion == ofxVersion1 {
		strictMode = false
	}

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == ofxDataElementName {
				file, err = r.readOFXElement(ctx, strictMode, ofxDataElementName)

				if err != nil {
					return nil, err
				}
			}
		}
	}

	if file == nil {
		log.Errorf(ctx, "[ofxFileReader.read] cannot parse ofx file")
		return nil, errs.ErrInvalidOFXFile
	}

	file.FileHeader = r.fileHeader

	return file, nil
}

func (r *ofxFileReader) readOFXElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxFile, error) {
	file := &ofxFile{}
	hasEndElement := false

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == ofxBankMessageResponseV1ElementName {
				element, err := r.readBankMessageResponseV1Element(ctx, strictMode, ofxBankMessageResponseV1ElementName)

				if err != nil {
					return nil, err
				}

				file.BankMessageResponseV1 = element
			} else if token.Name.Local == ofxCreditCardMessageResponseV1ElementName {
				element, err := r.readCreditCardMessageResponseV1Element(ctx, strictMode, ofxCreditCardMessageResponseV1ElementName)

				if err != nil {
					return nil, err
				}

				file.CreditCardMessageResponseV1 = element
			}
		case xml.EndElement:
			if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readOFXElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	return file, nil
}

func (r *ofxFileReader) readBankMessageResponseV1Element(ctx core.Context, strictMode bool, parentElementName string) (*ofxBankMessageResponseV1, error) {
	response := &ofxBankMessageResponseV1{}
	hasEndElement := false

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == ofxBankStatementTransactionResponseElementName {
				element, err := r.readBankStatementTransactionResponseElement(ctx, strictMode, ofxBankStatementTransactionResponseElementName)

				if err != nil {
					return nil, err
				}

				response.StatementTransactionResponse = element
			}
		case xml.EndElement:
			if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readBankMessageResponseV1Element] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	return response, nil
}

func (r *ofxFileReader) readCreditCardMessageResponseV1Element(ctx core.Context, strictMode bool, parentElementName string) (*ofxCreditCardMessageResponseV1, error) {
	response := &ofxCreditCardMessageResponseV1{}
	hasEndElement := false

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == ofxCreditCardStatementTransactionResponseElementName {
				element, err := r.readCreditCardStatementTransactionResponseElement(ctx, strictMode, ofxCreditCardStatementTransactionResponseElementName)

				if err != nil {
					return nil, err
				}

				response.StatementTransactionResponse = element
			}
		case xml.EndElement:
			if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readCreditCardMessageResponseV1Element] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	return response, nil
}

func (r *ofxFileReader) readBankStatementTransactionResponseElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxBankStatementTransactionResponse, error) {
	response := &ofxBankStatementTransactionResponse{}
	hasEndElement := false

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == ofxBankStatementResponseElementName {
				element, err := r.readBankStatementResponseElement(ctx, strictMode, ofxBankStatementResponseElementName)

				if err != nil {
					return nil, err
				}

				response.StatementResponse = element
			}
		case xml.EndElement:
			if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readBankStatementTransactionResponseElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	return response, nil
}

func (r *ofxFileReader) readCreditCardStatementTransactionResponseElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxCreditCardStatementTransactionResponse, error) {
	response := &ofxCreditCardStatementTransactionResponse{}
	hasEndElement := false

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == ofxCreditCardStatementResponseElementName {
				element, err := r.readCreditCardStatementResponseElement(ctx, strictMode, ofxCreditCardStatementResponseElementName)

				if err != nil {
					return nil, err
				}

				response.StatementResponse = element
			}
		case xml.EndElement:
			if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readCreditCardStatementTransactionResponseElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	return response, nil
}

func (r *ofxFileReader) readBankStatementResponseElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxBankStatementResponse, error) {
	response := &ofxBankStatementResponse{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxBankStatementResponseChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			} else if token.Name.Local == ofxBankStatementResponseBankAccountFromName {
				element, err := r.readBankAccountElement(ctx, strictMode, ofxBankStatementResponseBankAccountFromName)

				if err != nil {
					return nil, err
				}

				response.AccountFrom = element
			} else if token.Name.Local == ofxBankStatementResponseBankTransactionListName {
				element, err := r.readBankTransactionListElement(ctx, strictMode, ofxBankStatementResponseBankTransactionListName)

				if err != nil {
					return nil, err
				}

				response.TransactionList = element
			}
		case xml.EndElement:
			if ofxBankStatementResponseChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxBankStatementResponseChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readBankStatementResponseElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readBankStatementResponseElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxBankStatementResponseDefaultCurrencyName {
			response.DefaultCurrency = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return response, nil
}

func (r *ofxFileReader) readCreditCardStatementResponseElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxCreditCardStatementResponse, error) {
	response := &ofxCreditCardStatementResponse{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxCreditCardStatementResponseChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			} else if token.Name.Local == ofxCreditCardStatementResponseCreditCardAccountFromName {
				element, err := r.readCreditAccountElement(ctx, strictMode, ofxCreditCardStatementResponseCreditCardAccountFromName)

				if err != nil {
					return nil, err
				}

				response.AccountFrom = element
			} else if token.Name.Local == ofxCreditCardStatementResponseCreditCardTransactionListName {
				element, err := r.readCreditCardTransactionListElement(ctx, strictMode, ofxCreditCardStatementResponseCreditCardTransactionListName)

				if err != nil {
					return nil, err
				}

				response.TransactionList = element
			}
		case xml.EndElement:
			if ofxCreditCardStatementResponseChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxCreditCardStatementResponseChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readCreditCardStatementResponseElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readCreditCardStatementResponseElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxCreditCardStatementResponseDefaultCurrencyName {
			response.DefaultCurrency = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return response, nil
}

func (r *ofxFileReader) readBankAccountElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxBankAccount, error) {
	account := &ofxBankAccount{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxBankAccountChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			}
		case xml.EndElement:
			if ofxBankAccountChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxBankAccountChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readBankAccountElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readBankAccountElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxBankAccountBankIdName {
			account.BankId = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxBankAccountBranchIdName {
			account.BranchId = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxBankAccountAccountIdName {
			account.AccountId = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxBankAccountAccountTypeName {
			account.AccountType = ofxAccountType(r.getActualElementValue(name, value, elementNotHasEndElement, strictMode))
		} else if name == ofxBankAccountAccountKeyName {
			account.AccountKey = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return account, nil
}

func (r *ofxFileReader) readCreditAccountElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxCreditCardAccount, error) {
	account := &ofxCreditCardAccount{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxCreditCardAccountChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			}
		case xml.EndElement:
			if ofxCreditCardAccountChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxCreditCardAccountChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readCreditAccountElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readCreditAccountElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxCreditCardAccountAccountIdName {
			account.AccountId = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxCreditCardAccountAccountKeyName {
			account.AccountKey = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return account, nil
}

func (r *ofxFileReader) readBankTransactionListElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxBankTransactionList, error) {
	transactionList := &ofxBankTransactionList{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxBankTransactionListChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			} else if token.Name.Local == ofxBankTransactionListStatementTransactionsName {
				ofxBaseStatementTransaction, backAccountTo, _, err := r.readStatementTransactionElement(ctx, strictMode, "STMTTRN")

				if err != nil {
					return nil, err
				}

				transaction := &ofxBankStatementTransaction{
					ofxBaseStatementTransaction: *ofxBaseStatementTransaction,
					AccountTo:                   backAccountTo,
				}

				transactionList.StatementTransactions = append(transactionList.StatementTransactions, transaction)
			}
		case xml.EndElement:
			if ofxBankTransactionListChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxBankTransactionListChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readBankTransactionListElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readBankTransactionListElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxBankTransactionListStartDateName {
			transactionList.StartDate = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxBankTransactionListEndDateName {
			transactionList.EndDate = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return transactionList, nil
}

func (r *ofxFileReader) readCreditCardTransactionListElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxCreditCardTransactionList, error) {
	transactionList := &ofxCreditCardTransactionList{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxCreditCardTransactionListChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			} else if token.Name.Local == ofxCreditCardTransactionListStatementTransactionsName {
				ofxBaseStatementTransaction, _, creditCardAccountTo, err := r.readStatementTransactionElement(ctx, strictMode, "STMTTRN")

				if err != nil {
					return nil, err
				}

				transaction := &ofxCreditCardStatementTransaction{
					ofxBaseStatementTransaction: *ofxBaseStatementTransaction,
					AccountTo:                   creditCardAccountTo,
				}

				transactionList.StatementTransactions = append(transactionList.StatementTransactions, transaction)
			}
		case xml.EndElement:
			if ofxCreditCardTransactionListChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxCreditCardTransactionListChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readCreditCardTransactionListElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readCreditCardTransactionListElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxCreditCardTransactionListStartDateName {
			transactionList.StartDate = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxCreditCardTransactionListEndDateName {
			transactionList.EndDate = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return transactionList, nil
}

func (r *ofxFileReader) readStatementTransactionElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxBaseStatementTransaction, *ofxBankAccount, *ofxCreditCardAccount, error) {
	var bankAccountTo *ofxBankAccount
	var creditCardAccountTo *ofxCreditCardAccount
	transaction := &ofxBaseStatementTransaction{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxTransactionChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			} else if token.Name.Local == ofxTransactionPayeeName {
				element, err := r.readPayeeElement(ctx, strictMode, ofxTransactionPayeeName)

				if err != nil {
					return nil, nil, nil, err
				}

				transaction.Payee = element
			} else if token.Name.Local == ofxTransactionBankAccountToName {
				element, err := r.readBankAccountElement(ctx, strictMode, ofxTransactionBankAccountToName)

				if err != nil {
					return nil, nil, nil, err
				}

				bankAccountTo = element
			} else if token.Name.Local == ofxTransactionCreditCardAccountToName {
				element, err := r.readCreditAccountElement(ctx, strictMode, ofxTransactionCreditCardAccountToName)

				if err != nil {
					return nil, nil, nil, err
				}

				creditCardAccountTo = element
			}
		case xml.EndElement:
			if ofxTransactionChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxTransactionChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readStatementTransactionElement] not found </%s> element", parentElementName)
		return nil, nil, nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readStatementTransactionElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, nil, nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxTransactionTransactionIdName {
			transaction.TransactionId = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxTransactionTransactionTypeName {
			transaction.TransactionType = ofxTransactionType(r.getActualElementValue(name, value, elementNotHasEndElement, strictMode))
		} else if name == ofxTransactionPostedDateName {
			transaction.PostedDate = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxTransactionAmountName {
			transaction.Amount = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxTransactionNameName {
			transaction.Name = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxTransactionMemoName {
			transaction.Memo = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxTransactionCurrencyName {
			transaction.Currency = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxTransactionOriginalCurrencyName {
			transaction.OriginalCurrency = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return transaction, bankAccountTo, creditCardAccountTo, nil
}

func (r *ofxFileReader) readPayeeElement(ctx core.Context, strictMode bool, parentElementName string) (*ofxPayee, error) {
	payee := &ofxPayee{}
	hasEndElement := false
	elementNotHasEndElement := make(map[string]bool)
	elementValues := make(map[string]string)
	currentElementName := ""

	for {
		token, err := r.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if ofxPayeeChildrenNames[token.Name.Local] {
				currentElementName = token.Name.Local
				elementNotHasEndElement[token.Name.Local] = true
			}
		case xml.EndElement:
			if ofxPayeeChildrenNames[token.Name.Local] {
				delete(elementNotHasEndElement, token.Name.Local)
			} else if token.Name.Local == parentElementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if ofxPayeeChildrenNames[currentElementName] {
				elementValues[currentElementName] = string(token)
			}

			currentElementName = ""
		}

		if hasEndElement {
			break
		}
	}

	if strictMode && !hasEndElement {
		log.Errorf(ctx, "[ofxFileReader.readPayeeElement] not found </%s> element", parentElementName)
		return nil, errs.ErrInvalidOFXFile
	}

	if strictMode && len(elementNotHasEndElement) > 0 {
		log.Errorf(ctx, "[ofxFileReader.readPayeeElement] not found end element for %s", r.getNotHasEndElementNames(elementNotHasEndElement))
		return nil, errs.ErrInvalidOFXFile
	}

	for name, value := range elementValues {
		if name == ofxPayeeNameName {
			payee.Name = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeeAddress1Name {
			payee.Address1 = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeeAddress2Name {
			payee.Address2 = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeeAddress3Name {
			payee.Address3 = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeeCityName {
			payee.City = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeeStateName {
			payee.State = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeePostalCodeName {
			payee.PostalCode = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeeCountryName {
			payee.Country = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		} else if name == ofxPayeePhoneName {
			payee.Phone = r.getActualElementValue(name, value, elementNotHasEndElement, strictMode)
		}
	}

	return payee, nil
}

func (r *ofxFileReader) getActualElementValue(name string, value string, elementNotHasEndElement map[string]bool, strictMode bool) string {
	if strictMode {
		return value
	}

	_, notHasEndElement := elementNotHasEndElement[name]

	if !notHasEndElement {
		return value
	}

	for i := 0; i < len(value); i++ {
		if value[i] == '\r' || value[i] == '\n' {
			return value[0:i]
		}
	}

	return value
}

func (r *ofxFileReader) getNotHasEndElementNames(elementNotHasEndElement map[string]bool) string {
	builder := strings.Builder{}

	for name := range elementNotHasEndElement {
		if builder.Len() > 0 {
			builder.WriteRune(',')
		}

		builder.WriteString(name)
	}

	return builder.String()
}

func createNewOFXFileReader(ctx core.Context, data []byte) (*ofxFileReader, error) {
	firstNonCrLfIndex := 0

	for i := 0; i < len(data); i++ {
		if data[i] != '\n' && data[i] != '\r' {
			firstNonCrLfIndex = i
			break
		}
	}

	if len(data) > 5 && string(data[firstNonCrLfIndex:firstNonCrLfIndex+5]) == "<?xml" { // ofx 2.x starts with <?xml
		return createNewOFX2FileReader(ctx, data, true)
	} else if len(data) > 10 && string(data[firstNonCrLfIndex:firstNonCrLfIndex+10]) == "OFXHEADER:" { // ofx 1.x starts with OFXHEADER:
		return createNewOFX1FileReader(ctx, data)
	} else if len(data) > 5 && string(data[firstNonCrLfIndex:firstNonCrLfIndex+5]) == "<OFX>" { // no ofx header
		return createNewOFX2FileReader(ctx, data, false)
	}

	return nil, errs.ErrInvalidOFXFile
}

func createNewOFX1FileReader(ctx core.Context, data []byte) (*ofxFileReader, error) {
	fileHeader, fileData, dataType, enc, err := readOFX1FileHeader(ctx, data)

	if err != nil {
		return nil, err
	}

	if fileHeader.OFXDeclarationVersion != ofxVersion1 {
		log.Errorf(ctx, "[ofx_data_reader.createNewOFX1FileReader] cannot parse ofx 1.x file header, because declaration version is \"%s\"", fileHeader.OFXDeclarationVersion)
		return nil, errs.ErrInvalidOFXFile
	}

	if dataType != ofx1SGMLDataFormat {
		log.Errorf(ctx, "[ofx_data_reader.createNewOFX1FileReader] cannot parse ofx 1.x file header, because data type is \"%s\"", dataType)
		return nil, errs.ErrInvalidOFXFile
	}

	reader := bytes.NewReader(fileData)
	buffer := &bytes.Buffer{}

	if enc != nil {
		transformReader := transform.NewReader(reader, enc.NewDecoder())
		_, err = buffer.ReadFrom(transformReader)
	} else {
		_, err = buffer.ReadFrom(reader)
	}

	if err != nil {
		log.Errorf(ctx, "[ofx_data_reader.createNewOFX1FileReader] cannot read ofx 1.x file content, because %s", err.Error())
		return nil, errs.ErrInvalidOFXFile
	}

	sgmlData := buffer.String()
	stringReader := strings.NewReader(sgmlData)

	xmlDecoder := xml.NewDecoder(stringReader)
	xmlDecoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	return &ofxFileReader{
		fileHeader: fileHeader,
		xmlDecoder: xmlDecoder,
	}, nil
}

func createNewOFX2FileReader(ctx core.Context, data []byte, withHeader bool) (*ofxFileReader, error) {
	var fileHeader *ofxFileHeader = nil
	var err error

	if withHeader {
		fileHeader, err = readOFX2FileHeader(ctx, data)

		if err != nil {
			return nil, err
		}

		if fileHeader.OFXDeclarationVersion != ofxVersion2 {
			log.Errorf(ctx, "[ofx_data_reader.createNewOFX2FileReader] cannot parse ofx 2.x file header, because declaration version is \"%s\"", fileHeader.OFXDeclarationVersion)
			return nil, errs.ErrInvalidOFXFile
		}
	}

	xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	return &ofxFileReader{
		fileHeader: fileHeader,
		xmlDecoder: xmlDecoder,
	}, nil
}

func readOFX1FileHeader(ctx core.Context, data []byte) (fileHeader *ofxFileHeader, fileData []byte, dataType string, enc encoding.Encoding, err error) {
	fileHeader = &ofxFileHeader{}
	dataType = ""
	fileEncoding := ""
	fileCharset := ""
	fileDataStartPosition := 0
	lastCrLf := -1

	for i := 0; i < len(data); i++ {
		if data[i] != '\n' && data[i] != '\r' {
			continue
		}

		if lastCrLf == i-1 {
			lastCrLf = i
			continue
		}

		line := string(data[lastCrLf+1 : i])

		if strings.Index(line, "<OFX>") == 0 {
			fileDataStartPosition = lastCrLf + 1
			break
		}

		lastCrLf = i

		if line == "" {
			continue
		}

		items := strings.Split(line, ":")

		if len(items) != 2 {
			log.Warnf(ctx, "[ofx_data_reader.readOFX1FileHeader] cannot parse line in ofx 1.x file header, because line is \"%s\"", line)
			continue
		}

		key := items[0]
		value := items[1]

		if key == "OFXHEADER" {
			fileHeader.OFXDeclarationVersion = oFXDeclarationVersion(value)
		} else if key == "DATA" {
			dataType = value
		} else if key == "VERSION" {
			fileHeader.OFXDataVersion = value
		} else if key == "SECURITY" {
			fileHeader.Security = value
		} else if key == "ENCODING" {
			fileEncoding = strings.ToLower(value)
		} else if key == "CHARSET" {
			fileCharset = strings.ToLower(value)
		} else if key == "COMPRESSION" {
			continue // ignore
		} else if key == "OLDFILEUID" {
			fileHeader.OldFileUid = value
		} else if key == "NEWFILEUID" {
			fileHeader.NewFileUid = value
		} else {
			log.Warnf(ctx, "[ofx_data_reader.readOFX1FileHeader] cannot parse unknown header line in ofx 1.x file header, because line is \"%s\"", line)
			continue
		}
	}

	if fileEncoding == ofxUSAsciiEncoding {
		if utils.IsStringOnlyContainsDigits(fileCharset) {
			fileCharset = "cp" + fileCharset
		}

		enc, _ = charset.Lookup(fileCharset)

		if enc == nil {
			enc, _ = charset.Lookup("us-ascii")
		}

		if enc == nil {
			enc = charmap.Windows1252
		}
	} else if fileEncoding == ofxUnicodeEncoding {
		enc, _ = charset.Lookup(ofxUnicodeEncoding)

		if enc == nil {
			enc = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		}
	} else {
		log.Errorf(ctx, "[ofx_data_reader.readOFX1FileHeader] cannot parse ofx 1.x file, because encoding \"%s\" is unknown", fileEncoding)
		return nil, nil, "", nil, errs.ErrInvalidOFXFile
	}

	return fileHeader, data[fileDataStartPosition:], dataType, enc, nil
}

func readOFX2FileHeader(ctx core.Context, data []byte) (fileHeader *ofxFileHeader, err error) {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	fileHeader = &ofxFileHeader{}
	headerLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		ofxHeaderStartIndex := strings.Index(line, "<?OFX ")

		if ofxHeaderStartIndex >= 0 {
			headerLine = ofx2HeaderPattern.FindString(line)
			break
		}
	}

	if headerLine == "" {
		log.Errorf(ctx, "[ofx_data_reader.readOFX2FileHeader] cannot find ofx 2.x file header")
		return nil, errs.ErrInvalidOFXFile
	}

	headerAttributes := ofx2HeaderAttributePattern.FindAllStringSubmatch(headerLine, -1)

	for _, attributeItems := range headerAttributes {
		if len(attributeItems) != 3 {
			log.Warnf(ctx, "[ofx_data_reader.readOFX2FileHeader] cannot parse line in ofx 2.x file header, because item is \"%s\"", attributeItems)
			continue
		}

		name := attributeItems[1]
		value := attributeItems[2]

		if name == "OFXHEADER" {
			fileHeader.OFXDeclarationVersion = oFXDeclarationVersion(value)
		} else if name == "VERSION" {
			fileHeader.OFXDataVersion = value
		} else if name == "SECURITY" {
			fileHeader.Security = value
		} else if name == "OLDFILEUID" {
			fileHeader.OldFileUid = value
		} else if name == "NEWFILEUID" {
			fileHeader.NewFileUid = value
		} else {
			log.Warnf(ctx, "[ofx_data_reader.readOFX2FileHeader] cannot parse unknown header line in ofx 2.x file header, because item is \"%s\"", attributeItems)
			continue
		}
	}

	return fileHeader, nil
}

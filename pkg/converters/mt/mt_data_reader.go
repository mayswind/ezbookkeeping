package mt

import (
	"bufio"
	"bytes"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

const mtBasicHeaderBlockPrefix = "{1:"
const mtTextBlockStartPrefix = "{4:"
const mtTextBlockEndPrefix = "-}"
const mtTagPrefix = ':'
const mtInformationToAccountOwnerMaxLines = 6

const (
	mtTagStatementReferenceNumber  = ":20:"
	mtTagRelatedReference          = ":21:"
	mtTagAccountId                 = ":25:"
	mtTagSequentialNumber          = ":28C:"
	mtTagOpeningBalanceF           = ":60F:"
	mtTagOpeningBalanceM           = ":60M:"
	mtTagClosingBalanceF           = ":62F:"
	mtTagClosingBalanceM           = ":62M:"
	mtTagClosingAvailableBalance   = ":64:"
	mtTagStatementLine             = ":61:"
	mtTagInformationToAccountOwner = ":86:"
)

const (
	mtTransactionTypeSwiftTransfer    = 'S'
	mtTransactionTypeNonSwiftTransfer = 'N'
	mtTransactionTypeFirstAdvice      = 'F'
)

// mt940DataReader defines the structure of mt940 data reader
type mt940DataReader struct {
	allLines []string
}

// read returns the imported mt940 data
// Reference: https://www2.swift.com/knowledgecentre/publications/us9m_20230720/2.0?topic=mt940-format-spec.htm
func (r *mt940DataReader) read(ctx core.Context) (*mt940Data, error) {
	if len(r.allLines) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	data := &mt940Data{}
	var currentStatement *mtStatement
	var lastTag string

	for i := 0; i < len(r.allLines); i++ {
		line := strings.TrimSpace(r.allLines[i])

		if len(line) < 1 {
			continue
		}

		if strings.HasPrefix(line, mtBasicHeaderBlockPrefix) && strings.HasSuffix(line, mtTextBlockStartPrefix) {
			data = &mt940Data{}
			currentStatement = nil
			lastTag = ""
			continue
		} else if strings.HasPrefix(line, mtTextBlockEndPrefix) {
			break
		}

		if strings.HasPrefix(line, mtTagStatementReferenceNumber) {
			data.StatementReferenceNumber = line[len(mtTagStatementReferenceNumber):]
			lastTag = mtTagStatementReferenceNumber
		} else if strings.HasPrefix(line, mtTagRelatedReference) {
			data.RelatedReference = line[len(mtTagRelatedReference):]
			lastTag = mtTagRelatedReference
		} else if strings.HasPrefix(line, mtTagAccountId) {
			data.AccountId = line[len(mtTagAccountId):]
			lastTag = mtTagAccountId
		} else if strings.HasPrefix(line, mtTagSequentialNumber) {
			data.SequentialNumber = line[len(mtTagSequentialNumber):]
			lastTag = mtTagSequentialNumber
		} else if strings.HasPrefix(line, mtTagOpeningBalanceF) || strings.HasPrefix(line, mtTagOpeningBalanceM) {
			balance, err := r.parseBalance(ctx, line[len(mtTagOpeningBalanceF):])

			if err != nil {
				return nil, err
			}

			data.OpeningBalance = balance
			lastTag = line[:len(mtTagOpeningBalanceF)]
		} else if strings.HasPrefix(line, mtTagClosingBalanceF) || strings.HasPrefix(line, mtTagClosingBalanceM) {
			balance, err := r.parseBalance(ctx, line[len(mtTagClosingBalanceF):])

			if err != nil {
				return nil, err
			}

			data.ClosingBalance = balance
			lastTag = line[:len(mtTagClosingBalanceF)]
		} else if strings.HasPrefix(line, mtTagClosingAvailableBalance) {
			balance, err := r.parseBalance(ctx, line[len(mtTagClosingAvailableBalance):])

			if err != nil {
				return nil, err
			}

			data.ClosingAvailableBalance = balance
			lastTag = mtTagClosingAvailableBalance
		} else if strings.HasPrefix(line, mtTagStatementLine) {
			if currentStatement != nil {
				data.Statements = append(data.Statements, currentStatement)
			}

			statement, err := r.parseStatement(ctx, line[len(mtTagStatementLine):])

			if err != nil {
				return nil, err
			}

			currentStatement = statement
			lastTag = mtTagStatementLine
		} else if strings.HasPrefix(line, mtTagInformationToAccountOwner) && currentStatement != nil {
			currentStatement.InformationToAccountOwner = make([]string, 1)
			currentStatement.InformationToAccountOwner[0] = line[len(mtTagInformationToAccountOwner):]
			lastTag = mtTagInformationToAccountOwner
		} else if line[0] != mtTagPrefix && lastTag == mtTagStatementLine && currentStatement != nil {
			currentStatement.ReferenceForAccountOwner += line
			lastTag = ""
		} else if line[0] != mtTagPrefix && lastTag == mtTagInformationToAccountOwner && currentStatement != nil && len(currentStatement.InformationToAccountOwner) < mtInformationToAccountOwnerMaxLines {
			currentStatement.InformationToAccountOwner = append(currentStatement.InformationToAccountOwner, line)
			lastTag = mtTagInformationToAccountOwner
		} else {
			log.Warnf(ctx, "[mt_data_reader.read] unsupported line \"%s\" and skip this line", line)
		}
	}

	if currentStatement != nil {
		data.Statements = append(data.Statements, currentStatement)
	}

	return data, nil
}

func (r *mt940DataReader) parseBalance(ctx core.Context, data string) (*mtBalance, error) {
	// 1!a (debit/credit mark)
	// 6!n (date)
	// 3!a (currency)
	// 15d (amount)
	if len(data) < 9 {
		return nil, errs.ErrInvalidMT940File
	}

	if data[0] != MT_MARK_DEBIT[0] && data[0] != MT_MARK_CREDIT[0] {
		log.Errorf(ctx, "[mt_data_reader.parseBalance] cannot parse unknown debit/credit mark, current line is %s", data)
		return nil, errs.ErrTransactionTypeInvalid
	}

	balance := &mtBalance{
		DebitCreditMark: mtCreditDebitMark(data[0:1]),
		Date:            data[1:7],
		Currency:        data[7:10],
		Amount:          data[10:],
	}

	return balance, nil
}

func (r *mt940DataReader) parseStatement(ctx core.Context, data string) (*mtStatement, error) {
	// 6!n (value date)
	// [4!n] (entry date, optional)
	// 2a (debit/credit mark)
	// [1!a] (funds code, optional)
	// 15d (amount)
	// 1!a3!c (transaction type identification code)
	// 16x (reference for account owner)
	// [//16x] (reference of account servicing institution, optional)
	// [34x] (supplementary details, optional)
	if len(data) < 6 {
		return nil, errs.ErrInvalidMT940File
	}

	statement := &mtStatement{
		ValueDate: data[0:6],
	}

	currentIndex := 6

	// parse entry date if available
	if len(data) >= currentIndex+4 && '0' <= data[currentIndex] && data[currentIndex] <= '9' {
		statement.EntryDate = data[6:10]
		currentIndex += 4
	}

	// parse debit/credit indicator
	if len(data) >= currentIndex+1 && (data[currentIndex] == MT_MARK_DEBIT[0] || data[currentIndex] == MT_MARK_CREDIT[0]) {
		statement.CreditDebitMark = mtCreditDebitMark(data[currentIndex])
		currentIndex++
	} else if len(data) >= currentIndex+2 && (data[currentIndex:currentIndex+2] == string(MT_MARK_REVERSAL_CREDIT) || data[currentIndex:currentIndex+2] == string(MT_MARK_REVERSAL_DEBIT)) {
		statement.CreditDebitMark = mtCreditDebitMark(data[currentIndex : currentIndex+2])
		currentIndex += 2
	} else {
		log.Errorf(ctx, "[mt_data_reader.parseStatement] cannot parse unknown debit/credit mark, current line is %s", data)
		return nil, errs.ErrTransactionTypeInvalid
	}

	// parse funds code if available
	if len(data) >= currentIndex+1 && ('A' <= data[currentIndex] && data[currentIndex] <= 'Z') {
		statement.FundsCode = string(data[currentIndex])
		currentIndex++
	}

	// parse amount
	amountValue := ""
	for i := currentIndex; i < len(data); i++ {
		if len(amountValue) < 15 && ('0' <= data[i] && data[i] <= '9' || data[i] == ',') {
			amountValue += string(data[i])
		} else {
			currentIndex = i
			break
		}
	}
	statement.Amount = amountValue

	if len(statement.Amount) < 1 {
		log.Errorf(ctx, "[mt_data_reader.parseStatement] cannot parse amount, current line is %s", data)
		return nil, errs.ErrAmountInvalid
	}

	// parse transaction type identification code
	if len(data) >= currentIndex+4 && (data[currentIndex] == uint8(mtTransactionTypeSwiftTransfer) || data[currentIndex] == uint8(mtTransactionTypeNonSwiftTransfer) || data[currentIndex] == uint8(mtTransactionTypeFirstAdvice)) {
		statement.TransactionTypeIdentificationCode = data[currentIndex : currentIndex+4]
		currentIndex += 4
	} else {
		log.Errorf(ctx, "[mt_data_reader.parseStatement] cannot parse transaction type identification code, current line is %s", data)
		return nil, errs.ErrInvalidMT940File
	}

	// parse reference for account owner if available
	accountOwnerReference := ""
	for i := currentIndex; i < len(data); i++ {
		if len(accountOwnerReference) < 16 && (data[i] != '/' || (data[i] == '/' && (i >= len(data)-1 || data[i+1] != '/'))) {
			accountOwnerReference += string(data[i])
		} else {
			currentIndex = i
			break
		}
	}
	statement.ReferenceForAccountOwner = accountOwnerReference

	if len(statement.ReferenceForAccountOwner) < 1 {
		log.Errorf(ctx, "[mt_data_reader.parseStatement] cannot parse reference for account owner, current line is %s", data)
		return nil, errs.ErrInvalidMT940File
	}

	// parse reference of account servicing institution if available
	if len(data) >= currentIndex+3 && data[currentIndex] == '/' && data[currentIndex+1] == '/' {
		accountServicingInstitutionReference := ""
		currentIndex += 2
		for i := currentIndex; i < len(data); i++ {
			if len(accountServicingInstitutionReference) < 16 {
				accountServicingInstitutionReference += string(data[i])
			} else {
				currentIndex = i
				break
			}
		}
		statement.ReferenceOfAccountServicingInstitution = accountServicingInstitutionReference
	}

	return statement, nil
}

func createNewMT940FileReader(data []byte) *mt940DataReader {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))
	scanner := bufio.NewScanner(reader)
	allLines := make([]string, 0)

	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	return &mt940DataReader{
		allLines: allLines,
	}
}

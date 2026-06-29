package ai

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/llm"
	"github.com/mayswind/ezbookkeeping/pkg/llm/data"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/templates"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// aiTransactionDataParser defines the interface for parsing transaction data using AI
type aiTransactionDataParser struct {
	currentConfig *settings.Config
}

// aiTransactionDataParsedResult defines the structure of parsed transaction data result
type aiTransactionDataParsedResult struct {
	Transactions []*models.RecognizedTransactionResult `json:"transactions"`
}

// parse processes the input file data and returns the recognized transaction results using AI
func (p *aiTransactionDataParser) parse(c core.Context, user *models.User, fileData string, additionalPrompt string, defaultTimezone *time.Location, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) ([]*models.RecognizedTransactionResult, error) {
	if p.currentConfig == nil || p.currentConfig.TextRecognitionLLMConfig == nil || p.currentConfig.TextRecognitionLLMConfig.LLMProvider == "" || !p.currentConfig.TransactionFromAITextRecognition {
		return nil, errs.ErrLargeLanguageModelProviderNotEnabled
	}

	text := strings.TrimSpace(fileData)

	if len(text) == 0 {
		log.Warnf(c, "[ai_recognized_transaction_data_parser.parse] input text is empty for user \"uid:%d\"", user.Uid)
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	accountNames := p.getAccountNames(accountMap)
	expenseCategoryNames := p.getCategoryNames(expenseCategoryMap)
	incomeCategoryNames := p.getCategoryNames(incomeCategoryMap)
	transferCategoryNames := p.getCategoryNames(transferCategoryMap)
	tagNames := p.getTagNames(tagMap)

	systemPrompt, err := templates.GetTemplate(templates.SYSTEM_PROMPT_BATCH_TRANSACTION_TEXT_RECOGNITION)

	if err != nil {
		log.Errorf(c, "[ai_recognized_transaction_data_parser.parse] failed to get batch prompt template for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	systemPromptParams := map[string]any{
		"CurrentDateTime":          utils.FormatUnixTimeToLongDateTime(time.Now().Unix(), defaultTimezone),
		"AllExpenseCategoryNames":  strings.Join(expenseCategoryNames, "\n"),
		"AllIncomeCategoryNames":   strings.Join(incomeCategoryNames, "\n"),
		"AllTransferCategoryNames": strings.Join(transferCategoryNames, "\n"),
		"AllAccountNames":          strings.Join(accountNames, "\n"),
		"AllTagNames":              strings.Join(tagNames, "\n"),
		"AdditionalNotes":          additionalPrompt,
	}

	var bodyBuffer bytes.Buffer
	err = systemPrompt.Execute(&bodyBuffer, systemPromptParams)

	if err != nil {
		log.Errorf(c, "[ai_recognized_transaction_data_parser.parse] failed to render batch prompt template for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	llmRequest := &data.LargeLanguageModelRequest{
		Stream:         false,
		SystemPrompt:   strings.ReplaceAll(bodyBuffer.String(), "\r\n", "\n"),
		UserPrompt:     []byte(text),
		UserPromptType: data.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_TEXT,
	}

	llmResponse, err := llm.Container.GetJsonResponseByTextRecognitionModel(c, user.Uid, p.currentConfig, llmRequest)

	if err != nil {
		log.Errorf(c, "[ai_recognized_transaction_data_parser.parse] failed to get llm response for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	if llmResponse == nil || len(llmResponse.Content) == 0 || strings.HasPrefix(llmResponse.Content, "{}") {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	var result *aiTransactionDataParsedResult

	if err := json.Unmarshal([]byte(llmResponse.Content), &result); err != nil {
		log.Errorf(c, "[ai_recognized_transaction_data_parser.parse] failed to unmarshal batch llm response \"%s\" for user \"uid:%d\", because %s", llmResponse.Content, user.Uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	if result == nil || len(result.Transactions) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return result.Transactions, nil
}

func (p *aiTransactionDataParser) getAccountNames(accountMap map[string]*models.Account) []string {
	names := make([]string, 0, len(accountMap))

	for _, account := range accountMap {
		names = append(names, account.Name)
	}

	return names
}

func (p *aiTransactionDataParser) getCategoryNames(categoryMap map[string]map[string]*models.TransactionCategory) []string {
	names := make([]string, 0)

	for _, subCategoryMap := range categoryMap {
		for _, category := range subCategoryMap {
			names = append(names, category.Name)
		}
	}

	return names
}

func (p *aiTransactionDataParser) getTagNames(tagMap map[string]*models.TransactionTag) []string {
	names := make([]string, 0, len(tagMap))

	for _, tag := range tagMap {
		names = append(names, tag.Name)
	}

	return names
}

func createNewAITransactionDataParser(currentConfig *settings.Config) (*aiTransactionDataParser, error) {
	if currentConfig == nil || currentConfig.TextRecognitionLLMConfig == nil || currentConfig.TextRecognitionLLMConfig.LLMProvider == "" || !currentConfig.TransactionFromAITextRecognition {
		return nil, errs.ErrLargeLanguageModelProviderNotEnabled
	}

	return &aiTransactionDataParser{
		currentConfig: currentConfig,
	}, nil
}

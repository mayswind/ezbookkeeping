package mcp

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// MCPQueryTransactionsRequest represents all parameters of the query transactions request
type MCPQueryTransactionsRequest struct {
	StartTime             string `json:"start_time" jsonschema:"format=date-time" jsonschema_description:"Start time for the query in RFC 3339 format (e.g. 2023-01-01T12:00:00Z)"`
	EndTime               string `json:"end_time" jsonschema:"format=date-time" jsonschema_description:"End time for the query in RFC 3339 format or (e.g. 2023-01-01T12:00:00Z)"`
	Type                  string `json:"type,omitempty" jsonschema:"enum=income,enum=expense,enum=transfer" jsonschema_description:"Transaction type to filter by (income, expense, transfer) (optional)"`
	SecondaryCategoryName string `json:"category_name,omitempty" jsonschema_description:"Primary or secondary category name to filter transactions by (optional)"`
	AccountName           string `json:"account_name,omitempty" jsonschema_description:"Account name to filter transactions by (optional)"`
	Keyword               string `json:"keyword,omitempty" jsonschema_description:"Keyword to search in transaction description (optional)"`
	Count                 int32  `json:"count,omitempty" jsonschema:"default=100" jsonschema_description:"Maximum number of results to return (default: 100)"`
	Page                  int32  `json:"page,omitempty" jsonschema:"default=1" jsonschema_description:"Page number for pagination (default: 1)"`
	ResponseFields        string `json:"response_fields,omitempty" jsonschema_description:"Comma-separated list of fields to include in the response (optional, leave empty for all fields, available fields: time, currency, category_name, account_name, comment)"`
}

// MCPQueryTransactionsResponse represents the response structure for querying transactions
type MCPQueryTransactionsResponse struct {
	TotalCount   int64                 `json:"total_count" jsonschema_description:"Total number of transactions matching the query"`
	CurrentPage  int32                 `json:"current_page" jsonschema_description:"Current page number of the results"`
	TotalPage    int32                 `json:"total_page" jsonschema_description:"Total number of pages available for the query, calculated based on total_count and count"`
	Transactions []*MCPTransactionInfo `json:"transactions" jsonschema_description:"List of transactions matching the query"`
}

// MCPTransactionInfo defines the structure of transaction information
type MCPTransactionInfo struct {
	Time                   string `json:"time,omitempty" jsonschema_description:"Time of the transaction in RFC 3339 format (e.g. 2023-01-01T12:00:00Z)"`
	Type                   string `json:"type" jsonschema:"enum=income,enum=expense,enum=transfer" jsonschema_description:"Transaction type (income, expense, transfer)"`
	Amount                 string `json:"amount" jsonschema_description:"Amount of the transaction in the specified currency"`
	Currency               string `json:"currency,omitempty" jsonschema_description:"Currency code of the transaction (e.g. USD, EUR)"`
	SecondaryCategoryName  string `json:"category_name,omitempty" jsonschema_description:"Secondary category name for the transaction"`
	AccountName            string `json:"account_name,omitempty" jsonschema_description:"Account name for the transaction"`
	DestinationAmount      string `json:"destination_amount,omitempty" jsonschema_description:"Destination amount for transfer transactions (optional)"`
	DestinationCurrency    string `json:"destination_currency,omitempty" jsonschema_description:"Currency code of the destination amount for transfer transactions (optional)"`
	DestinationAccountName string `json:"destination_account_name,omitempty" jsonschema_description:"Destination account name for transfer transactions (optional)"`
	Comment                string `json:"comment,omitempty" jsonschema_description:"Description of the transaction"`
}

type mcpQueryTransactionsToolHandler struct{}

var MCPQueryTransactionsToolHandler = &mcpQueryTransactionsToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpQueryTransactionsToolHandler) Name() string {
	return "query_transactions"
}

// Description returns the description of the MCP tool
func (h *mcpQueryTransactionsToolHandler) Description() string {
	return "Query transactions based on various filters."
}

// InputType returns the input type for the MCP tool request
func (h *mcpQueryTransactionsToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPQueryTransactionsRequest{})
}

// OutputType returns the output type for the MCP tool response
func (h *mcpQueryTransactionsToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPQueryTransactionsResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpQueryTransactionsToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var queryTransactionsRequest MCPQueryTransactionsRequest

	if callToolReq.Arguments != nil {
		if err := json.Unmarshal(callToolReq.Arguments, &queryTransactionsRequest); err != nil {
			return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	uid := user.Uid
	maxTime, err := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(queryTransactionsRequest.EndTime)

	if err != nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	minTime, err := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(queryTransactionsRequest.StartTime)

	if err != nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(maxTime.Unix())
	minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(minTime.Unix())

	if queryTransactionsRequest.Count <= 0 {
		queryTransactionsRequest.Count = 100
	}

	if queryTransactionsRequest.Page <= 0 {
		queryTransactionsRequest.Page = 1
	}

	transactionType := models.TransactionType(byte(0))

	if queryTransactionsRequest.Type == transactionTypeExpense {
		transactionType = models.TRANSACTION_TYPE_EXPENSE
	} else if queryTransactionsRequest.Type == transactionTypeIncome {
		transactionType = models.TRANSACTION_TYPE_INCOME
	} else if queryTransactionsRequest.Type == transactionTypeTransfer {
		transactionType = models.TRANSACTION_TYPE_TRANSFER
	}

	allAccounts, err := services.GetAccountService().GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Warnf(c, "[add_transaction.Handle] get account error, because %s", err.Error())
		return nil, nil, err
	}

	filterAccountIds := make([]int64, 0)

	if queryTransactionsRequest.AccountName != "" {
		filterAccountIds = services.GetAccountService().GetAccountOrSubAccountIdsByAccountName(allAccounts, queryTransactionsRequest.AccountName)

		if len(filterAccountIds) < 1 {
			return nil, nil, errs.ErrAccountNotFound
		}
	}

	allCategories, err := services.GetTransactionCategoryService().GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.Warnf(c, "[add_transaction.Handle] get transaction category error, because %s", err.Error())
		return nil, nil, err
	}

	filterCategoryIds := make([]int64, 0)

	if queryTransactionsRequest.SecondaryCategoryName != "" {
		filterCategoryIds = services.GetTransactionCategoryService().GetCategoryOrSubCategoryIdsByCategoryName(allCategories, queryTransactionsRequest.SecondaryCategoryName)

		if len(filterCategoryIds) < 1 {
			return nil, nil, errs.ErrTransactionCategoryNotFound
		}
	}

	totalCount, err := services.GetTransactionService().GetTransactionCount(c, uid, maxTransactionTime, minTransactionTime, transactionType, filterCategoryIds, filterAccountIds, nil, false, "", queryTransactionsRequest.Keyword)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionListHandler] failed to get transaction count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, err
	}

	transactions, err := services.GetTransactionService().GetTransactionsByMaxTime(c, uid, maxTransactionTime, minTransactionTime, transactionType, filterCategoryIds, filterAccountIds, nil, false, "", queryTransactionsRequest.Keyword, queryTransactionsRequest.Page, queryTransactionsRequest.Count, false, true)
	structuredResponse, response, err := h.createNewMCPQueryTransactionsResponse(c, &queryTransactionsRequest, transactions, totalCount, services.GetAccountService().GetAccountMapByList(allAccounts), services.GetTransactionCategoryService().GetCategoryMapByList(allCategories))

	if err != nil {
		return nil, nil, err
	}

	return structuredResponse, response, nil
}

func (h *mcpQueryTransactionsToolHandler) createNewMCPQueryTransactionsResponse(c *core.WebContext, queryTransactionsRequest *MCPQueryTransactionsRequest, transactions []*models.Transaction, totalCount int64, accountsMap map[int64]*models.Account, categoriesMap map[int64]*models.TransactionCategory) (any, []*MCPTextContent, error) {
	response := MCPQueryTransactionsResponse{
		TotalCount:   totalCount,
		CurrentPage:  queryTransactionsRequest.Page,
		TotalPage:    int32((totalCount + int64(queryTransactionsRequest.Count) - 1) / int64(queryTransactionsRequest.Count)),
		Transactions: make([]*MCPTransactionInfo, 0, len(transactions)),
	}

	filteredFields := make(map[string]bool)

	if queryTransactionsRequest.ResponseFields != "" {
		for _, field := range strings.Split(queryTransactionsRequest.ResponseFields, ",") {
			filteredFields[field] = true
		}
	}

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		transactionInfo := MCPTransactionInfo{
			Amount: utils.FormatAmount(transaction.Amount),
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			transactionInfo.Type = transactionTypeExpense
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			transactionInfo.Type = transactionTypeIncome
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			transactionInfo.Type = transactionTypeTransfer
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			transactionInfo.DestinationAmount = utils.FormatAmount(transaction.RelatedAccountAmount)
		}

		if _, exists := filteredFields["time"]; exists || len(filteredFields) == 0 {
			transactionUnixTime := utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime)
			transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
			transactionInfo.Time = utils.FormatUnixTimeToLongDateTimeWithTimezoneRFC3339Format(transactionUnixTime, transactionTimeZone)
		}

		if _, exists := filteredFields["currency"]; exists || len(filteredFields) == 0 {
			if account, exists := accountsMap[transaction.AccountId]; exists && account != nil {
				transactionInfo.Currency = account.Currency
			}

			if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT && transaction.RelatedAccountId > 0 {
				if destinationAccount, exists := accountsMap[transaction.RelatedAccountId]; exists && destinationAccount != nil {
					transactionInfo.DestinationCurrency = destinationAccount.Currency
				}
			}
		}

		if _, exists := filteredFields["category_name"]; exists || len(filteredFields) == 0 {
			if category, exists := categoriesMap[transaction.CategoryId]; exists && category != nil {
				transactionInfo.SecondaryCategoryName = category.Name
			}
		}

		if _, exists := filteredFields["account_name"]; exists || len(filteredFields) == 0 {
			if account, exists := accountsMap[transaction.AccountId]; exists && account != nil {
				transactionInfo.AccountName = account.Name
			}

			if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT && transaction.RelatedAccountId > 0 {
				if destinationAccount, exists := accountsMap[transaction.RelatedAccountId]; exists && destinationAccount != nil {
					transactionInfo.DestinationAccountName = destinationAccount.Name
				}
			}
		}

		if _, exists := filteredFields["comment"]; exists || len(filteredFields) == 0 {
			transactionInfo.Comment = transaction.Comment
		}

		response.Transactions = append(response.Transactions, &transactionInfo)
	}

	content, err := json.Marshal(response)

	if err != nil {
		return nil, nil, err
	}

	return response, []*MCPTextContent{
		NewMCPTextContent(string(content)),
	}, nil
}

package mcp

import (
	"encoding/json"
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const transactionTypeIncome = "income"
const transactionTypeExpense = "expense"
const transactionTypeTransfer = "transfer"

// MCPAddTransactionRequest represents all parameters of the add transaction request
type MCPAddTransactionRequest struct {
	Type                   string   `json:"type" jsonschema:"enum=income,enum=expense,enum=transfer" jsonschema_description:"Transaction type (income, expense, transfer)"`
	Time                   string   `json:"time" jsonschema:"format=date-time" jsonschema_description:"Transaction time in RFC 3339 format (e.g. 2023-01-01T12:00:00Z)"`
	SecondaryCategoryName  string   `json:"category_name" jsonschema_description:"Secondary category name for the transaction"`
	AccountName            string   `json:"account_name" jsonschema_description:"Account name for the transaction"`
	Amount                 string   `json:"amount" jsonschema_description:"Transaction amount"`
	DestinationAccountName string   `json:"destination_account_name,omitempty" jsonschema_description:"Destination account name for transfer transactions (optional)"`
	DestinationAmount      string   `json:"destination_amount,omitempty" jsonschema_description:"Destination amount for transfer transactions (optional)"`
	Tags                   []string `json:"tags,omitempty" jsonschema_description:"List of tags associated with the transaction (optional, maximum 10 tags allowed)"`
	Comment                string   `json:"comment,omitempty" jsonschema_description:"Transaction description"`
	DryRun                 bool     `json:"dry_run,omitempty" jsonschema_description:"If true, the transaction will not be saved, only validated (optional)"`
}

// MCPAddTransactionResponse represents the response structure for add transaction
type MCPAddTransactionResponse struct {
	Success                   bool   `json:"success" jsonschema_description:"Indicates whether the transaction was added successfully"`
	DryRun                    bool   `json:"dry_run,omitempty" jsonschema_description:"Indicates whether this is a dry run (transaction not saved actually)"`
	AccountBalance            string `json:"account_balance,omitempty" jsonschema_description:"Account balance (or outstanding balance for debt accounts) after the transaction"`
	DestinationAccountBalance string `json:"destination_account_balance,omitempty" jsonschema_description:"Destination account balance (or outstanding balance for debt accounts) after the transaction (only for transfer transactions)"`
}

type mcpAddTransactionToolHandler struct{}

var MCPAddTransactionToolHandler = &mcpAddTransactionToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpAddTransactionToolHandler) Name() string {
	return "add_transaction"
}

// Description returns the description of the MCP tool
func (h *mcpAddTransactionToolHandler) Description() string {
	return "Add a new transaction in ezBookkeeping."
}

// InputType returns the input type for the MCP tool request
func (h *mcpAddTransactionToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPAddTransactionRequest{})
}

// OutputType returns the output type for the MCP tool response
func (h *mcpAddTransactionToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPAddTransactionResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpAddTransactionToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var addTransactionRequest MCPAddTransactionRequest

	if callToolReq.Arguments != nil {
		if err := json.Unmarshal(callToolReq.Arguments, &addTransactionRequest); err != nil {
			return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if addTransactionRequest.Type == transactionTypeTransfer {
		if addTransactionRequest.DestinationAccountName == "" || addTransactionRequest.DestinationAmount == "" {
			return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
		}
	}

	if len(addTransactionRequest.Tags) > models.MaximumTagsCountOfTransaction {
		return nil, nil, errs.ErrTransactionHasTooManyTags
	}

	uid := user.Uid
	allAccounts, err := services.GetAccountService().GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Warnf(c, "[add_transaction.Handle] get account error, because %s", err.Error())
		return nil, nil, err
	}

	accountsMap := services.GetAccountService().GetVisibleAccountNameMapByList(allAccounts)
	sourceAccount, exists := accountsMap[addTransactionRequest.AccountName]

	if !exists {
		log.Warnf(c, "[add_transaction.Handle] source account \"%s\" not found for user \"uid:%d\"", addTransactionRequest.AccountName, uid)
		return nil, nil, errs.ErrSourceAccountNotFound
	}

	var destinationAccount *models.Account
	destinationAccountId := int64(0)

	if addTransactionRequest.Type == transactionTypeTransfer {
		destinationAccount, exists = accountsMap[addTransactionRequest.DestinationAccountName]

		if !exists {
			log.Warnf(c, "[add_transaction.Handle] destination account \"%s\" not found for user \"uid:%d\"", addTransactionRequest.DestinationAccountName, uid)
			return nil, nil, errs.ErrDestinationAccountNotFound
		}

		destinationAccountId = destinationAccount.AccountId
	}

	allCategories, err := services.GetTransactionCategoryService().GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.Warnf(c, "[add_transaction.Handle] get transaction category error, because %s", err.Error())
		return nil, nil, err
	}

	categoriesMap := services.GetTransactionCategoryService().GetVisibleCategoryNameMapByList(allCategories)
	category, exists := categoriesMap[addTransactionRequest.SecondaryCategoryName]

	if !exists {
		log.Warnf(c, "[add_transaction.Handle] secondary category \"%s\" not found for user \"uid:%d\"", addTransactionRequest.SecondaryCategoryName, uid)
		return nil, nil, errs.ErrTransactionCategoryNotFound
	}

	var tagIds []int64

	if len(addTransactionRequest.Tags) > 0 {
		allTags, err := services.GetTransactionTagService().GetAllTagsByUid(c, uid)

		if err != nil {
			log.Warnf(c, "[add_transaction.Handle] get transaction tag ids error, because %s", err.Error())
			return nil, nil, err
		}

		tagMaps := services.GetTransactionTagService().GetTagNameMapByList(allTags)
		tagIds = make([]int64, 0, len(addTransactionRequest.Tags))

		for _, tagName := range addTransactionRequest.Tags {
			if tag, exists := tagMaps[tagName]; exists {
				tagIds = append(tagIds, tag.TagId)
			} else {
				log.Warnf(c, "[add_transaction.Handle] transaction tag \"%s\" not found for user \"uid:%d\"", tagName, uid)
			}
		}
	}

	transaction := h.createNewTransactionModel(uid, &addTransactionRequest, category.CategoryId, sourceAccount.AccountId, destinationAccountId, c.ClientIP())
	transactionEditable := user.CanEditTransactionByTransactionTime(transaction.TransactionTime, transaction.TimezoneUtcOffset)

	if !transactionEditable {
		return nil, nil, errs.ErrCannotCreateTransactionWithThisTransactionTime
	}

	if !addTransactionRequest.DryRun {
		err = services.GetTransactionService().CreateTransaction(c, transaction, tagIds, nil)

		if err != nil {
			log.Errorf(c, "[add_transaction.Handle] failed to create transaction \"id:%d\" for user \"uid:%d\", because %s", transaction.TransactionId, uid, err.Error())
			return nil, nil, err
		}

		log.Infof(c, "[add_transaction.Handle] user \"uid:%d\" has created a new transaction \"id:%d\" successfully", uid, transaction.TransactionId)

		accountIds := []int64{sourceAccount.AccountId}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			accountIds = append(accountIds, destinationAccountId)
		}

		newAccounts, err := services.GetAccountService().GetAccountsByAccountIds(c, uid, accountIds)

		if err != nil {
			log.Warnf(c, "[add_transaction.Handle] failed to get latest accounts info after transaction created, because %s", err.Error())
		}

		structuredResponse, response, err := h.createNewMCPAddTransactionResponse(c, transaction, newAccounts, false)

		if err != nil {
			return nil, nil, err
		}

		return structuredResponse, response, nil
	} else {
		newAccounts := make(map[int64]*models.Account)
		newAccounts[sourceAccount.AccountId] = sourceAccount

		if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			sourceAccount.Balance -= transaction.Amount
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			sourceAccount.Balance += transaction.Amount
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT && destinationAccount != nil {
			newAccounts[destinationAccount.AccountId] = destinationAccount
			destinationAccount.Balance += transaction.RelatedAccountAmount
		}

		structuredResponse, response, err := h.createNewMCPAddTransactionResponse(c, transaction, newAccounts, true)

		if err != nil {
			return nil, nil, err
		}

		return structuredResponse, response, nil
	}
}

func (h *mcpAddTransactionToolHandler) createNewTransactionModel(uid int64, addTransactionRequest *MCPAddTransactionRequest, categoryId int64, sourceAccountId int64, destinationAccountId int64, clientIp string) *models.Transaction {
	var transactionDbType models.TransactionDbType

	if addTransactionRequest.Type == transactionTypeExpense {
		transactionDbType = models.TRANSACTION_DB_TYPE_EXPENSE
	} else if addTransactionRequest.Type == transactionTypeIncome {
		transactionDbType = models.TRANSACTION_DB_TYPE_INCOME
	} else if addTransactionRequest.Type == transactionTypeTransfer {
		transactionDbType = models.TRANSACTION_DB_TYPE_TRANSFER_OUT
	}

	transactionTime, err := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(addTransactionRequest.Time)

	if err != nil {
		return nil
	}

	amount, err := utils.ParseAmount(addTransactionRequest.Amount)

	if err != nil {
		return nil
	}

	transaction := &models.Transaction{
		Uid:               uid,
		Type:              transactionDbType,
		CategoryId:        categoryId,
		TransactionTime:   utils.GetMinTransactionTimeFromUnixTime(transactionTime.Unix()),
		TimezoneUtcOffset: utils.GetTimezoneOffsetMinutes(transactionTime.Location()),
		AccountId:         sourceAccountId,
		Amount:            amount,
		HideAmount:        false,
		Comment:           addTransactionRequest.Comment,
		CreatedIp:         clientIp,
	}

	if addTransactionRequest.Type == transactionTypeTransfer {
		transaction.RelatedAccountId = destinationAccountId

		destinationAmount, err := utils.ParseAmount(addTransactionRequest.DestinationAmount)

		if err != nil {
			return nil
		}

		transaction.RelatedAccountAmount = destinationAmount
	}

	return transaction
}

func (h *mcpAddTransactionToolHandler) createNewMCPAddTransactionResponse(c *core.WebContext, transaction *models.Transaction, accountsMap map[int64]*models.Account, dryRun bool) (any, []*MCPTextContent, error) {
	var sourceAccountInfo *models.AccountInfoResponse
	var destinationAccountInfo *models.AccountInfoResponse

	if sourceAccount, exists := accountsMap[transaction.AccountId]; exists {
		sourceAccountInfo = sourceAccount.ToAccountInfoResponse()
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if destinationAccount, exists := accountsMap[transaction.RelatedAccountId]; exists {
			destinationAccountInfo = destinationAccount.ToAccountInfoResponse()
		}
	}

	response := MCPAddTransactionResponse{
		Success: true,
		DryRun:  dryRun,
	}

	if sourceAccountInfo != nil {
		if sourceAccountInfo.IsAsset {
			response.AccountBalance = utils.FormatAmount(sourceAccountInfo.Balance)
		} else if sourceAccountInfo.IsLiability {
			response.AccountBalance = utils.FormatAmount(-sourceAccountInfo.Balance)
		}
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT && destinationAccountInfo != nil {
		if destinationAccountInfo.IsAsset {
			response.DestinationAccountBalance = utils.FormatAmount(destinationAccountInfo.Balance)
		} else if destinationAccountInfo.IsLiability {
			response.DestinationAccountBalance = utils.FormatAmount(-destinationAccountInfo.Balance)
		}
	}

	content, err := json.Marshal(response)

	if err != nil {
		return nil, nil, err
	}

	return response, []*MCPTextContent{
		NewMCPTextContent(string(content)),
	}, nil
}

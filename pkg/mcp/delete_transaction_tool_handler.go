package mcp

import (
	"encoding/json/v2"
	"reflect"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// MCPDeleteTransactionRequest represents all parameters of the delete transaction request
type MCPDeleteTransactionRequest struct {
	TransactionId string `json:"transaction_id" jsonschema_description:"Id of the transaction to delete (as a string to avoid precision loss)"`
	DryRun        bool   `json:"dry_run,omitempty" jsonschema_description:"If true, the transaction will not be deleted, only validated (optional)"`
}

// MCPDeleteTransactionResponse represents the response structure for delete transaction
type MCPDeleteTransactionResponse struct {
	Success       bool   `json:"success" jsonschema_description:"Indicates whether this operation is successful"`
	DryRun        bool   `json:"dry_run,omitzero" jsonschema_description:"Indicates whether this operation is a dry run (transaction not deleted actually)"`
	TransactionId string `json:"transaction_id" jsonschema_description:"Id of the deleted transaction"`
}

type mcpDeleteTransactionToolHandler struct{}

var MCPDeleteTransactionToolHandler = &mcpDeleteTransactionToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpDeleteTransactionToolHandler) Name() string {
	return "delete_transaction"
}

// Description returns the description of the MCP tool
func (h *mcpDeleteTransactionToolHandler) Description() string {
	return "Delete an existing transaction in ezBookkeeping. For a transfer, use the id of the transfer out transaction; the paired transfer in transaction is deleted as well. Set dry_run to validate without deleting."
}

// InputType returns the input type for the MCP tool request
func (h *mcpDeleteTransactionToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTransactionRequest{})
}

// OutputType returns the output type for the MCP tool response
func (h *mcpDeleteTransactionToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTransactionResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpDeleteTransactionToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var deleteTransactionRequest MCPDeleteTransactionRequest

	if callToolReq.Arguments == nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if err := json.Unmarshal(callToolReq.Arguments, &deleteTransactionRequest); err != nil {
		return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	transactionId, err := utils.StringToInt64(deleteTransactionRequest.TransactionId)

	if err != nil || transactionId <= 0 {
		return nil, nil, errs.ErrTransactionIdInvalid
	}

	uid := user.Uid
	transaction, err := services.GetTransactionService().GetTransactionByTransactionId(c, uid, transactionId)

	if err != nil {
		log.Errorf(c, "[delete_transaction_tool_handler.Handle] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", transactionId, uid, err.Error())
		return nil, nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		log.Warnf(c, "[delete_transaction_tool_handler.Handle] cannot delete transaction \"id:%d\" for user \"uid:%d\", because transaction type is transfer in", transactionId, uid)
		return nil, nil, errs.ErrTransactionTypeInvalid
	}

	usedAccounts, err := h.getTransactionUsedAccounts(c, uid, transaction, services)

	if err != nil {
		return nil, nil, err
	}

	transactionTimezone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
	transactionEditable := user.CanEditTransactionByTransactionTime(transaction.TransactionTime, transactionTimezone, usedAccounts[transaction.AccountId], usedAccounts[transaction.RelatedAccountId])

	if !transactionEditable {
		log.Warnf(c, "[delete_transaction_tool_handler.Handle] transaction \"id:%d\" is not editable for user \"uid:%d\"", transactionId, uid)
		return nil, nil, errs.ErrCannotDeleteTransactionWithThisTransactionTime
	}

	if !deleteTransactionRequest.DryRun {
		err = services.GetTransactionService().DeleteTransaction(c, uid, transactionId)

		if err != nil {
			log.Errorf(c, "[delete_transaction_tool_handler.Handle] failed to delete transaction \"id:%d\" for user \"uid:%d\", because %s", transactionId, uid, err.Error())
			return nil, nil, errs.Or(err, errs.ErrOperationFailed)
		}

		log.Infof(c, "[delete_transaction_tool_handler.Handle] user \"uid:%d\" has deleted transaction \"id:%d\"", uid, transactionId)
	}

	return h.createMCPDeleteTransactionResponse(transaction.TransactionId, deleteTransactionRequest.DryRun)
}

func (h *mcpDeleteTransactionToolHandler) getTransactionUsedAccounts(c *core.WebContext, uid int64, transaction *models.Transaction, services MCPAvailableServices) (map[int64]*models.Account, error) {
	accountIds := []int64{transaction.AccountId}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		accountIds = append(accountIds, transaction.RelatedAccountId)
	}

	accountMap, err := services.GetAccountService().GetAccountsByAccountIds(c, uid, accountIds)

	if err != nil {
		log.Errorf(c, "[delete_transaction_tool_handler.getTransactionUsedAccounts] failed to get accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if _, exists := accountMap[transaction.AccountId]; !exists {
		log.Warnf(c, "[delete_transaction_tool_handler.getTransactionUsedAccounts] account of transaction \"id:%d\" does not exist for user \"uid:%d\"", transaction.TransactionId, uid)
		return nil, errs.ErrSourceAccountNotFound
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if _, exists := accountMap[transaction.RelatedAccountId]; !exists {
			log.Warnf(c, "[delete_transaction_tool_handler.getTransactionUsedAccounts] related account of transaction \"id:%d\" does not exist for user \"uid:%d\"", transaction.TransactionId, uid)
			return nil, errs.ErrDestinationAccountNotFound
		}
	}

	return accountMap, nil
}

func (h *mcpDeleteTransactionToolHandler) createMCPDeleteTransactionResponse(transactionId int64, dryRun bool) (any, []*MCPTextContent, error) {
	response := MCPDeleteTransactionResponse{
		Success:       true,
		DryRun:        dryRun,
		TransactionId: utils.Int64ToString(transactionId),
	}

	content, err := json.Marshal(response)

	if err != nil {
		return nil, nil, err
	}

	return response, []*MCPTextContent{
		NewMCPTextContent(string(content)),
	}, nil
}

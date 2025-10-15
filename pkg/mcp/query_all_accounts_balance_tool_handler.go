package mcp

import (
	"encoding/json"
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// MCPQueryAllAccountsBalanceResponse represents the response structure for querying accounts balance
type MCPQueryAllAccountsBalanceResponse struct {
	CashAccounts                 []*MCPAccountBalanceInfo `json:"cashAccounts,omitempty" jsonschema_description:"List of cash account balances"`
	CheckingAccounts             []*MCPAccountBalanceInfo `json:"checkingAccounts,omitempty" jsonschema_description:"List of checking account balances"`
	SavingsAccounts              []*MCPAccountBalanceInfo `json:"savingsAccounts,omitempty" jsonschema_description:"List of savings account balances"`
	CreditCardAccounts           []*MCPAccountBalanceInfo `json:"creditCardAccounts,omitempty" jsonschema_description:"List of credit card account outstanding balances"`
	VirtualAccounts              []*MCPAccountBalanceInfo `json:"virtualAccounts,omitempty" jsonschema_description:"List of virtual account balances"`
	DebtAccounts                 []*MCPAccountBalanceInfo `json:"debtAccounts,omitempty" jsonschema_description:"List of debt account outstanding balances"`
	ReceivableAccounts           []*MCPAccountBalanceInfo `json:"receivableAccounts,omitempty" jsonschema_description:"List of receivable account balances"`
	CertificateOfDepositAccounts []*MCPAccountBalanceInfo `json:"certificateOfDepositAccounts,omitempty" jsonschema_description:"List of certificate of deposit account balances"`
	InvestmentAccounts           []*MCPAccountBalanceInfo `json:"investmentAccounts,omitempty" jsonschema_description:"List of investment account balances"`
}

// MCPAccountBalanceInfo defines the structure of account balance information
type MCPAccountBalanceInfo struct {
	Name               string `json:"name" jsonschema_description:"Account name"`
	Type               string `json:"type" jsonschema:"enum=asset,enum=liability" jsonschema_description:"Account type (asset or liability)"`
	Balance            string `json:"balance,omitempty" jsonschema_description:"Current balance of the account"`
	OutstandingBalance string `json:"outstandingBalance,omitempty" jsonschema_description:"Current outstanding balance of the account (positive value indicates amount owed)"`
	Currency           string `json:"currency" jsonschema_description:"Currency code of the account (e.g. USD, EUR)"`
}

type mcpQueryAllAccountsBalanceToolHandler struct{}

var MCPQueryAllAccountsBalanceToolHandler = &mcpQueryAllAccountsBalanceToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpQueryAllAccountsBalanceToolHandler) Name() string {
	return "query_all_accounts_balance"
}

// Description returns the description of the MCP tool
func (h *mcpQueryAllAccountsBalanceToolHandler) Description() string {
	return "Query all accounts balance for the current user in ezBookkeeping."
}

// InputType returns the input type for the MCP tool request
func (h *mcpQueryAllAccountsBalanceToolHandler) InputType() reflect.Type {
	return nil
}

// OutputType returns the output type for the MCP tool response
func (h *mcpQueryAllAccountsBalanceToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPQueryAllAccountsBalanceResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpQueryAllAccountsBalanceToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	uid := user.Uid
	accounts, err := services.GetAccountService().GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[query_all_accounts_balance_tool_handler.Handle] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, err
	}

	structuredResponse, response, err := h.createNewMCPQueryAllAccountsBalanceResponse(c, accounts)

	if err != nil {
		return nil, nil, err
	}

	return structuredResponse, response, nil
}

func (h *mcpQueryAllAccountsBalanceToolHandler) createNewMCPQueryAllAccountsBalanceResponse(c *core.WebContext, accounts []*models.Account) (any, []*MCPTextContent, error) {
	response := MCPQueryAllAccountsBalanceResponse{}

	for i := 0; i < len(accounts); i++ {
		account := accounts[i]

		if account.Hidden || (account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS && account.ParentAccountId == models.LevelOneAccountParentId) {
			continue
		}

		if account.Category == models.ACCOUNT_CATEGORY_CASH {
			if response.CashAccounts == nil {
				response.CashAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.CashAccounts = append(response.CashAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_CHECKING_ACCOUNT {
			if response.CheckingAccounts == nil {
				response.CheckingAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.CheckingAccounts = append(response.CheckingAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_SAVINGS_ACCOUNT {
			if response.SavingsAccounts == nil {
				response.SavingsAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.SavingsAccounts = append(response.SavingsAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
			if response.CreditCardAccounts == nil {
				response.CreditCardAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.CreditCardAccounts = append(response.CreditCardAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_VIRTUAL {
			if response.VirtualAccounts == nil {
				response.VirtualAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.VirtualAccounts = append(response.VirtualAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_DEBT {
			if response.DebtAccounts == nil {
				response.DebtAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.DebtAccounts = append(response.DebtAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_RECEIVABLES {
			if response.ReceivableAccounts == nil {
				response.ReceivableAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.ReceivableAccounts = append(response.ReceivableAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT {
			if response.CertificateOfDepositAccounts == nil {
				response.CertificateOfDepositAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.CertificateOfDepositAccounts = append(response.CertificateOfDepositAccounts, h.createNewMCPAccountBalanceInfo(account))
		} else if account.Category == models.ACCOUNT_CATEGORY_INVESTMENT {
			if response.InvestmentAccounts == nil {
				response.InvestmentAccounts = make([]*MCPAccountBalanceInfo, 0)
			}

			response.InvestmentAccounts = append(response.InvestmentAccounts, h.createNewMCPAccountBalanceInfo(account))
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

func (h *mcpQueryAllAccountsBalanceToolHandler) createNewMCPAccountBalanceInfo(account *models.Account) *MCPAccountBalanceInfo {
	accountResp := account.ToAccountInfoResponse()

	balanceInfo := &MCPAccountBalanceInfo{
		Name:     accountResp.Name,
		Currency: accountResp.Currency,
	}

	if accountResp.IsAsset {
		balanceInfo.Type = "asset"
		balanceInfo.Balance = utils.FormatAmount(accountResp.Balance)
	} else if accountResp.IsLiability {
		balanceInfo.Type = "liability"
		balanceInfo.OutstandingBalance = utils.FormatAmount(-accountResp.Balance)
	}

	return balanceInfo
}

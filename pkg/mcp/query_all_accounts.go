package mcp

import (
	"encoding/json"
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPQueryAllAccountsResponse represents the response structure for querying accounts
type MCPQueryAllAccountsResponse struct {
	CashAccounts                 []string `json:"cashAccounts,omitempty" jsonschema_description:"List of cash account names"`
	CheckingAccounts             []string `json:"checkingAccounts,omitempty" jsonschema_description:"List of checking account names"`
	SavingsAccounts              []string `json:"savingsAccounts,omitempty" jsonschema_description:"List of savings account names"`
	CreditCardAccounts           []string `json:"creditCardAccounts,omitempty" jsonschema_description:"List of credit card account names"`
	VirtualAccounts              []string `json:"virtualAccounts,omitempty" jsonschema_description:"List of virtual account names"`
	DebtAccounts                 []string `json:"debtAccounts,omitempty" jsonschema_description:"List of debt account names"`
	ReceivableAccounts           []string `json:"receivableAccounts,omitempty" jsonschema_description:"List of receivable account names"`
	CertificateOfDepositAccounts []string `json:"certificateOfDepositAccounts,omitempty" jsonschema_description:"List of certificate of deposit account names"`
	InvestmentAccounts           []string `json:"investmentAccounts,omitempty" jsonschema_description:"List of investment account names"`
}

type mcpQueryAllAccountsToolHandler struct{}

var MCPQueryAllAccountsToolHandler = &mcpQueryAllAccountsToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpQueryAllAccountsToolHandler) Name() string {
	return "query_all_accounts"
}

// Description returns the description of the MCP tool
func (h *mcpQueryAllAccountsToolHandler) Description() string {
	return "Query all accounts for the current user in ezBookkeeping."
}

// InputType returns the input type for the MCP tool request
func (h *mcpQueryAllAccountsToolHandler) InputType() reflect.Type {
	return nil
}

// OutputType returns the output type for the MCP tool response
func (h *mcpQueryAllAccountsToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPQueryAllAccountsResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpQueryAllAccountsToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, *errs.Error) {
	uid := c.GetCurrentUid()
	accounts, err := services.GetAccountService().GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[query_all_accounts.Handle] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, errs.Or(err, errs.ErrOperationFailed)
	}

	structuredResponse, response, err := h.createNewMCPQueryAllAccountsResponse(c, accounts)

	if err != nil {
		return nil, nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return structuredResponse, response, nil
}

func (h *mcpQueryAllAccountsToolHandler) createNewMCPQueryAllAccountsResponse(c *core.WebContext, accounts []*models.Account) (any, []*MCPTextContent, error) {
	response := MCPQueryAllAccountsResponse{}

	for i := 0; i < len(accounts); i++ {
		account := accounts[i]

		if account.Hidden || (account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS && account.ParentAccountId == models.LevelOneAccountParentId) {
			continue
		}

		if account.Category == models.ACCOUNT_CATEGORY_CASH {
			if response.CashAccounts == nil {
				response.CashAccounts = make([]string, 0)
			}

			response.CashAccounts = append(response.CashAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_CHECKING_ACCOUNT {
			if response.CheckingAccounts == nil {
				response.CheckingAccounts = make([]string, 0)
			}

			response.CheckingAccounts = append(response.CheckingAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_SAVINGS_ACCOUNT {
			if response.SavingsAccounts == nil {
				response.SavingsAccounts = make([]string, 0)
			}

			response.SavingsAccounts = append(response.SavingsAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
			if response.CreditCardAccounts == nil {
				response.CreditCardAccounts = make([]string, 0)
			}

			response.CreditCardAccounts = append(response.CreditCardAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_VIRTUAL {
			if response.VirtualAccounts == nil {
				response.VirtualAccounts = make([]string, 0)
			}

			response.VirtualAccounts = append(response.VirtualAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_DEBT {
			if response.DebtAccounts == nil {
				response.DebtAccounts = make([]string, 0)
			}

			response.DebtAccounts = append(response.DebtAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_RECEIVABLES {
			if response.ReceivableAccounts == nil {
				response.ReceivableAccounts = make([]string, 0)
			}

			response.ReceivableAccounts = append(response.ReceivableAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT {
			if response.CertificateOfDepositAccounts == nil {
				response.CertificateOfDepositAccounts = make([]string, 0)
			}

			response.CertificateOfDepositAccounts = append(response.CertificateOfDepositAccounts, account.Name)
		} else if account.Category == models.ACCOUNT_CATEGORY_INVESTMENT {
			if response.InvestmentAccounts == nil {
				response.InvestmentAccounts = make([]string, 0)
			}

			response.InvestmentAccounts = append(response.InvestmentAccounts, account.Name)
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

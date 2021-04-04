package cli

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/exporters"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
)

const pageCountForGettingTransactions = 1000
const pageCountForDataExport = 1000

// UserDataCli represents user data cli
type UserDataCli struct {
	csvExporter  *exporters.CSVFileExporter
	accounts     *services.AccountService
	transactions *services.TransactionService
	categories   *services.TransactionCategoryService
	tags         *services.TransactionTagService
	users        *services.UserService
}

// Initialize an user data cli singleton instance
var (
	UserData = &UserDataCli{
		csvExporter:  &exporters.CSVFileExporter{},
		accounts:     services.Accounts,
		transactions: services.Transactions,
		users:        services.Users,
		categories:   services.TransactionCategories,
		tags:         services.TransactionTags,
	}
)

// GetUserByUsername returns user by user name
func (a *UserDataCli) GetUserByUsername(c *cli.Context, username string) (*models.User, error) {
	if username == "" {
		log.BootErrorf("[user_data.GetUserByUsername] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	user, err := a.users.GetUserByUsername(username)

	if err != nil {
		log.BootErrorf("[user_data.GetUserByUsername] failed to get user by user name \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes user according to the specified user name
func (a *UserDataCli) DeleteUser(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.DeleteUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := a.users.DeleteUser(username)

	if err != nil {
		log.BootErrorf("[user_data.DeleteUser] failed to delete user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// CheckTransactionAndAccount checks whether all user transactions and all user accounts are correct
func (a *UserDataCli) CheckTransactionAndAccount(c *cli.Context, uid int64) (bool, error) {
	accountMap, categoryMap, tagMap, tagIndexs, err := a.getUserEssentialData(uid)

	if err != nil {
		log.BootErrorf("[user_data.CheckTransactionAndAccount] failed to get essential data for user \"uid:%d\", because %s", uid, err.Error())
		return false, err
	}

	accountHasChild := make(map[int64]bool)

	for _, account := range accountMap {
		if account.ParentAccountId > models.LevelOneAccountParentId {
			accountHasChild[account.ParentAccountId] = true
		}
	}

	allTransactions, err := a.transactions.GetAllTransactions(uid, pageCountForGettingTransactions, false)

	if err != nil {
		log.BootErrorf("[user_data.CheckTransactionAndAccount] failed to all transactions for user \"uid:%d\", because %s", uid, err.Error())
		return false, err
	}

	transactionMap := a.transactions.GetTransactionMapByList(allTransactions)
	accountBalance := make(map[int64]int64)

	for i := len(allTransactions) - 1; i >= 0; i-- {
		transaction := allTransactions[i]

		err := a.checkTransactionAccount(c, transaction, accountMap, accountHasChild)

		if err != nil {
			return false, err
		}

		err = a.checkTransactionCategory(c, transaction, categoryMap)

		if err != nil {
			return false, err
		}

		err = a.checkTransactionTag(c, transaction.TransactionId, tagIndexs, tagMap)

		if err != nil {
			return false, err
		}

		err = a.checkTransactionRelatedTransaction(c, transaction, transactionMap, accountMap)

		if err != nil {
			return false, err
		}

		balance, exists := accountBalance[transaction.AccountId]

		if !exists {
			balance = 0
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			balance = balance + transaction.RelatedAccountAmount
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			balance = balance + transaction.Amount
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			balance = balance - transaction.Amount
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			balance = balance - transaction.Amount
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			balance = balance + transaction.Amount
		} else {
			log.BootErrorf("[user_data.CheckAccountBalance] transaction type of transaction \"id:%d\" is invalid", transaction.TransactionId)
			return false, errs.ErrOperationFailed
		}

		accountBalance[transaction.AccountId] = balance
	}

	for _, account := range accountMap {
		actualBalance, exists := accountBalance[account.AccountId]

		if !exists && account.Balance == 0 {
			continue
		}

		if !exists && account.Balance != 0 {
			log.BootErrorf("[user_data.CheckAccountBalance] account \"id:%d\" balance is not correct, expected balance is %d, but there is no transaction actually", account.AccountId, account.Balance)
			return false, errs.ErrOperationFailed
		}

		if account.Balance != actualBalance {
			log.BootErrorf("[user_data.CheckAccountBalance] account \"id:%d\" balance is not correct, expected balance is %d, but actual balance is %d", account.AccountId, account.Balance, actualBalance)
			return false, errs.ErrOperationFailed
		}
	}

	for accountId, actualBalance := range accountBalance {
		_, exists := accountMap[accountId]

		if !exists {
			log.BootErrorf("[user_data.CheckAccountBalance] account \"id:%d\" does not exist, but there are some transactions of this account actually, and actual balance is %d", accountId, actualBalance)
			return false, errs.ErrOperationFailed
		}
	}

	return true, nil
}

// ExportTransaction returns csv file content according user all transactions
func (a *UserDataCli) ExportTransaction(c *cli.Context, uid int64) ([]byte, error) {
	accountMap, categoryMap, tagMap, tagIndexs, err := a.getUserEssentialData(uid)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to get essential data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, err
	}

	allTransactions, err := a.transactions.GetAllTransactions(uid, pageCountForDataExport, true)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to all transactions for user \"uid:%d\", because %s", uid, err.Error())
		return nil, err
	}

	result, err := a.csvExporter.GetOutputContent(uid, time.Local, allTransactions, accountMap, categoryMap, tagMap, tagIndexs)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to get csv format exported data for \"uid:%d\", because %s", uid, err.Error())
		return nil, err
	}

	return result, nil
}

// GetUserIdByUsername returns user id by user name
func (a *UserDataCli) GetUserIdByUsername(c *cli.Context, username string) (int64, error) {
	user, err := a.GetUserByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.GetUserIdByUsername] failed to get user by user name \"%s\", because %s", username, err.Error())
		return 0, err
	}

	return user.Uid, nil
}

func (a *UserDataCli) getUserEssentialData(uid int64) (accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, tagIndexs map[int64][]int64, err error) {
	if uid <= 0 {
		log.BootErrorf("[user_data.getUserEssentialData] user uid \"%d\" is invalid", uid)
		return nil, nil, nil, nil, errs.ErrUserIdInvalid
	}

	accounts, err := a.accounts.GetAllAccountsByUid(uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, nil, nil, err
	}

	accountMap = a.accounts.GetAccountMapByList(accounts)

	categories, err := a.categories.GetAllCategoriesByUid(uid, 0, -1)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, nil, nil, err
	}

	categoryMap = a.categories.GetCategoryMapByList(categories)

	tags, err := a.tags.GetAllTagsByUid(uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, nil, nil, err
	}

	tagMap = a.tags.GetTagMapByList(tags)

	tagIndexs, err = a.tags.GetAllTagIdsOfAllTransactions(uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get tag index for user \"uid:%d\", because %s", uid, err.Error())
		return nil, nil, nil, nil, err
	}

	return accountMap, categoryMap, tagMap, tagIndexs, nil
}

func (a *UserDataCli) checkTransactionAccount(c *cli.Context, transaction *models.Transaction, accountMap map[int64]*models.Account, accountHasChild map[int64]bool) error {
	account, exists := accountMap[transaction.AccountId]

	if !exists {
		log.BootErrorf("[user_data.checkTransactionAccount] the account \"id:%d\" of transaction \"id:%d\" does not exist", transaction.AccountId, transaction.TransactionId)
		return errs.ErrAccountNotFound
	}

	if account.ParentAccountId == models.LevelOneAccountParentId && accountHasChild[account.AccountId] {
		log.BootErrorf("[user_data.checkTransactionAccount] the account \"id:%d\" of transaction \"id:%d\" is not a sub account", transaction.AccountId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountId > 0 {
		relatedAccount, exists := accountMap[transaction.RelatedAccountId]

		if !exists {
			log.BootErrorf("[user_data.checkTransactionAccount] the related account \"id:%d\" of transaction \"id:%d\" does not exist", transaction.RelatedAccountId, transaction.TransactionId)
			return errs.ErrAccountNotFound
		}

		if relatedAccount.ParentAccountId == models.LevelOneAccountParentId && accountHasChild[relatedAccount.AccountId] {
			log.BootErrorf("[user_data.checkTransactionAccount] the related account \"id:%d\" of transaction \"id:%d\" is not a sub account", transaction.RelatedAccountId, transaction.TransactionId)
			return errs.ErrOperationFailed
		}
	}

	return nil
}

func (a *UserDataCli) checkTransactionCategory(c *cli.Context, transaction *models.Transaction, categoryMap map[int64]*models.TransactionCategory) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.CategoryId > 0 {
			log.BootErrorf("[user_data.checkTransactionCategory] transaction \"id:%d\" is balance modification transaction, but has category \"id:%d\"", transaction.TransactionId, transaction.CategoryId)
			return errs.ErrBalanceModificationTransactionCannotSetCategory
		} else {
			return nil
		}
	}

	category, exists := categoryMap[transaction.CategoryId]

	if !exists {
		log.BootErrorf("[user_data.checkTransactionCategory] the transaction category \"id:%d\" of transaction \"id:%d\" does not exist", transaction.CategoryId, transaction.TransactionId)
		return errs.ErrTransactionCategoryNotFound
	}

	if category.ParentCategoryId == models.LevelOneTransactionParentId {
		log.BootErrorf("[user_data.checkTransactionCategory] the transaction category \"id:%d\" of transaction \"id:%d\" is not a sub category", transaction.CategoryId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	return nil
}

func (a *UserDataCli) checkTransactionTag(c *cli.Context, transactionId int64, allTagIndexs map[int64][]int64, tagMap map[int64]*models.TransactionTag) error {
	tagIndexs, exists := allTagIndexs[transactionId]

	if !exists {
		return nil
	}

	for i := 0; i < len(tagIndexs); i++ {
		tagIndex := tagIndexs[i]
		tag, exists := tagMap[tagIndex]

		if !exists {
			log.BootErrorf("[user_data.checkTransactionTag] the transaction tag \"id:%d\" of transaction \"id:%d\" does not exist", tag.TagId, transactionId)
			return errs.ErrTransactionTagNotFound
		}
	}

	return nil
}

func (a *UserDataCli) checkTransactionRelatedTransaction(c *cli.Context, transaction *models.Transaction, transactionMap map[int64]*models.Transaction, accountMap map[int64]*models.Account) error {
	if transaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_OUT && transaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return nil
	}

	relatedTransaction, exists := transactionMap[transaction.RelatedId]

	if !exists {
		log.BootErrorf("[user_data.checkTransactionRelatedTransaction] the related transaction \"id:%d\" of transaction \"id:%d\" does not exist", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrTransactionNotFound
	}

	if transaction.RelatedId != relatedTransaction.TransactionId || transaction.TransactionId != relatedTransaction.RelatedId {
		log.BootErrorf("[user_data.checkTransactionRelatedTransaction] related ids of transaction \"id:%d\" and transaction \"id:%d\" are not equal", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountId != relatedTransaction.AccountId || transaction.AccountId != relatedTransaction.RelatedAccountId {
		log.BootErrorf("[user_data.checkTransactionRelatedTransaction] related account ids of transaction \"id:%d\" and transaction \"id:%d\" are not equal", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountAmount != relatedTransaction.Amount || transaction.Amount != relatedTransaction.RelatedAccountAmount {
		log.BootErrorf("[user_data.checkTransactionRelatedTransaction] related amounts of transaction \"id:%d\" and transaction \"id:%d\" are not equal", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	account := accountMap[transaction.AccountId]
	relatedAccount := accountMap[transaction.RelatedAccountId]

	if account.Currency == relatedAccount.Currency && transaction.Amount != transaction.RelatedAccountAmount {
		log.BootWarnf("[user_data.checkTransactionRelatedTransaction] transfer-in amount and transfer-out amount of transaction \"id:%d\" are not equal", transaction.TransactionId)
	}

	return nil
}

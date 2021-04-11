package cli

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/exporters"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const pageCountForGettingTransactions = 1000
const pageCountForDataExport = 1000

// UserDataCli represents user data cli
type UserDataCli struct {
	csvExporter             *exporters.CSVFileExporter
	accounts                *services.AccountService
	transactions            *services.TransactionService
	categories              *services.TransactionCategoryService
	tags                    *services.TransactionTagService
	users                   *services.UserService
	twoFactorAuthorizations *services.TwoFactorAuthorizationService
	tokens                  *services.TokenService
}

// Initialize an user data cli singleton instance
var (
	UserData = &UserDataCli{
		csvExporter:             &exporters.CSVFileExporter{},
		accounts:                services.Accounts,
		transactions:            services.Transactions,
		categories:              services.TransactionCategories,
		tags:                    services.TransactionTags,
		users:                   services.Users,
		twoFactorAuthorizations: services.TwoFactorAuthorizations,
		tokens:                  services.Tokens,
	}
)

// AddNewUser adds a new user according to specified info
func (l *UserDataCli) AddNewUser(c *cli.Context, username string, email string, nickname string, password string, defaultCurrency string) (*models.User, error) {
	if username == "" {
		log.BootErrorf("[user_data.AddNewUser] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	if email == "" {
		log.BootErrorf("[user_data.AddNewUser] user email is empty")
		return nil, errs.ErrEmailIsEmpty
	}

	if nickname == "" {
		log.BootErrorf("[user_data.AddNewUser] user nickname is empty")
		return nil, errs.ErrNicknameIsEmpty
	}

	if password == "" {
		log.BootErrorf("[user_data.AddNewUser] user password is empty")
		return nil, errs.ErrPasswordIsEmpty
	}

	if defaultCurrency == "" {
		log.BootErrorf("[user_data.AddNewUser] user default currency is empty")
		return nil, errs.ErrUserDefaultCurrencyIsEmpty
	}

	if _, ok := validators.AllCurrencyNames[defaultCurrency]; !ok {
		log.BootErrorf("[user_data.AddNewUser] user default currency is invalid")
		return nil, errs.ErrUserDefaultCurrencyIsInvalid
	}

	user := &models.User{
		Username:             username,
		Email:                email,
		Nickname:             nickname,
		Password:             password,
		DefaultCurrency:      defaultCurrency,
		FirstDayOfWeek:       models.WEEKDAY_SUNDAY,
		TransactionEditScope: models.TRANSACTION_EDIT_SCOPE_ALL,
	}

	err := l.users.CreateUser(user)

	if err != nil {
		log.BootErrorf("[user_data.AddNewUser] failed to create user \"%s\", because %s", user.Username, err.Error())
		return nil, err
	}

	log.BootInfof("[user_data.AddNewUser] user \"%s\" has add successfully, uid is %d", user.Username, user.Uid)

	return user, nil
}

// GetUserByUsername returns user by user name
func (l *UserDataCli) GetUserByUsername(c *cli.Context, username string) (*models.User, error) {
	if username == "" {
		log.BootErrorf("[user_data.GetUserByUsername] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	user, err := l.users.GetUserByUsername(username)

	if err != nil {
		log.BootErrorf("[user_data.GetUserByUsername] failed to get user by user name \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return user, nil
}

// ModifyUserPassword modifies user password
func (l *UserDataCli) ModifyUserPassword(c *cli.Context, username string, password string) error {
	if username == "" {
		log.BootErrorf("[user_data.ModifyUserPassword] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	if password == "" {
		log.BootErrorf("[user_data.ModifyUserPassword] user password is empty")
		return errs.ErrPasswordIsEmpty
	}

	user, err := l.users.GetUserByUsername(username)

	if err != nil {
		log.BootErrorf("[user_data.ModifyUserPassword] failed to get user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	if l.users.IsPasswordEqualsUserPassword(password, user) {
		return errs.ErrNothingWillBeUpdated
	}

	userNew := &models.User{
		Uid:      user.Uid,
		Salt:     user.Salt,
		Password: password,
	}

	_, err = l.users.UpdateUser(userNew)

	if err != nil {
		log.BootErrorf("[user_data.ModifyUserPassword] failed to update user \"%s\" password, because %s", user.Username, err.Error())
		return err
	}

	now := time.Now().Unix()
	err = l.tokens.DeleteTokensBeforeTime(user.Uid, now)

	if err == nil {
		log.BootInfof("[user_data.ModifyUserPassword] revoke old tokens before unix time \"%d\" for user \"%s\"", now, user.Username)
	} else {
		log.BootWarnf("[user_data.ModifyUserPassword] failed to revoke old tokens for user \"%s\", because %s", user.Username, err.Error())
	}

	return nil
}

// DeleteUser deletes user according to the specified user name
func (l *UserDataCli) DeleteUser(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.DeleteUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.DeleteUser(username)

	if err != nil {
		log.BootErrorf("[user_data.DeleteUser] failed to delete user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// ClearUserTokens clears all tokens of the specified user
func (l *UserDataCli) ClearUserTokens(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.ClearUserTokens] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.ClearUserTokens] error occurs when getting user id by user name")
		return err
	}

	now := time.Now().Unix()
	err = l.tokens.DeleteTokensBeforeTime(uid, now)

	if err != nil {
		log.BootErrorf("[user_data.ClearUserTokens] failed to delete tokens of user \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// DisableUserTwoFactorAuthorization disables 2fa for the specified user
func (l *UserDataCli) DisableUserTwoFactorAuthorization(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] error occurs when getting user id by user name")
		return err
	}

	enableTwoFactor, err := l.twoFactorAuthorizations.ExistsTwoFactorSetting(uid)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] failed to check two factor setting, because %s", err.Error())
		return err
	}

	if !enableTwoFactor {
		return errs.ErrTwoFactorIsNotEnabled
	}

	err = l.twoFactorAuthorizations.DeleteTwoFactorRecoveryCodes(uid)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] failed to delete two factor recovery codes for user \"%s\"", username)
		return err
	}

	err = l.twoFactorAuthorizations.DeleteTwoFactorSetting(uid)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] failed to delete two factor setting for user \"%s\"", username)
		return err
	}

	return nil
}

// CheckTransactionAndAccount checks whether all user transactions and all user accounts are correct
func (l *UserDataCli) CheckTransactionAndAccount(c *cli.Context, username string) (bool, error) {
	if username == "" {
		log.BootErrorf("[user_data.CheckTransactionAndAccount] user name is empty")
		return false, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.CheckTransactionAndAccount] error occurs when getting user id by user name")
		return false, err
	}

	accountMap, categoryMap, tagMap, tagIndexs, err := l.getUserEssentialData(uid, username)

	if err != nil {
		log.BootErrorf("[user_data.CheckTransactionAndAccount] failed to get essential data for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	accountHasChild := make(map[int64]bool)

	for _, account := range accountMap {
		if account.ParentAccountId > models.LevelOneAccountParentId {
			accountHasChild[account.ParentAccountId] = true
		}
	}

	allTransactions, err := l.transactions.GetAllTransactions(uid, pageCountForGettingTransactions, false)

	if err != nil {
		log.BootErrorf("[user_data.CheckTransactionAndAccount] failed to all transactions for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	transactionMap := l.transactions.GetTransactionMapByList(allTransactions)
	accountBalance := make(map[int64]int64)

	for i := len(allTransactions) - 1; i >= 0; i-- {
		transaction := allTransactions[i]

		err := l.checkTransactionAccount(c, transaction, accountMap, accountHasChild)

		if err != nil {
			return false, err
		}

		err = l.checkTransactionCategory(c, transaction, categoryMap)

		if err != nil {
			return false, err
		}

		err = l.checkTransactionTag(c, transaction.TransactionId, tagIndexs, tagMap)

		if err != nil {
			return false, err
		}

		err = l.checkTransactionRelatedTransaction(c, transaction, transactionMap, accountMap)

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
func (l *UserDataCli) ExportTransaction(c *cli.Context, username string) ([]byte, error) {
	if username == "" {
		log.BootErrorf("[user_data.ExportTransaction] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] error occurs when getting user id by user name")
		return nil, err
	}

	accountMap, categoryMap, tagMap, tagIndexs, err := l.getUserEssentialData(uid, username)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to get essential data for user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	allTransactions, err := l.transactions.GetAllTransactions(uid, pageCountForDataExport, true)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to all transactions for user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	result, err := l.csvExporter.GetOutputContent(uid, time.Local, allTransactions, accountMap, categoryMap, tagMap, tagIndexs)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to get csv format exported data for \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return result, nil
}

func (l *UserDataCli) getUserIdByUsername(c *cli.Context, username string) (int64, error) {
	user, err := l.GetUserByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.getUserIdByUsername] failed to get user by user name \"%s\", because %s", username, err.Error())
		return 0, err
	}

	return user.Uid, nil
}

func (l *UserDataCli) getUserEssentialData(uid int64, username string) (accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, tagIndexs map[int64][]int64, err error) {
	if uid <= 0 {
		log.BootErrorf("[user_data.getUserEssentialData] user uid \"%d\" is invalid", uid)
		return nil, nil, nil, nil, errs.ErrUserIdInvalid
	}

	accounts, err := l.accounts.GetAllAccountsByUid(uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get accounts for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, err
	}

	accountMap = l.accounts.GetAccountMapByList(accounts)

	categories, err := l.categories.GetAllCategoriesByUid(uid, 0, -1)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get categories for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, err
	}

	categoryMap = l.categories.GetCategoryMapByList(categories)

	tags, err := l.tags.GetAllTagsByUid(uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get tags for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, err
	}

	tagMap = l.tags.GetTagMapByList(tags)

	tagIndexs, err = l.tags.GetAllTagIdsOfAllTransactions(uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get tag index for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, err
	}

	return accountMap, categoryMap, tagMap, tagIndexs, nil
}

func (l *UserDataCli) checkTransactionAccount(c *cli.Context, transaction *models.Transaction, accountMap map[int64]*models.Account, accountHasChild map[int64]bool) error {
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

func (l *UserDataCli) checkTransactionCategory(c *cli.Context, transaction *models.Transaction, categoryMap map[int64]*models.TransactionCategory) error {
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

func (l *UserDataCli) checkTransactionTag(c *cli.Context, transactionId int64, allTagIndexs map[int64][]int64, tagMap map[int64]*models.TransactionTag) error {
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

func (l *UserDataCli) checkTransactionRelatedTransaction(c *cli.Context, transaction *models.Transaction, transactionMap map[int64]*models.Transaction, accountMap map[int64]*models.Account) error {
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

package cli

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/ezbookkeeping/pkg/converters"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const pageCountForGettingTransactions = 1000
const pageCountForDataExport = 1000

// UserDataCli represents user data cli
type UserDataCli struct {
	ezBookKeepingCsvExporter *converters.EzBookKeepingCSVFileExporter
	ezBookKeepingTsvExporter *converters.EzBookKeepingTSVFileExporter
	accounts                 *services.AccountService
	transactions             *services.TransactionService
	categories               *services.TransactionCategoryService
	tags                     *services.TransactionTagService
	users                    *services.UserService
	twoFactorAuthorizations  *services.TwoFactorAuthorizationService
	tokens                   *services.TokenService
	forgetPasswords          *services.ForgetPasswordService
}

// Initialize an user data cli singleton instance
var (
	UserData = &UserDataCli{
		ezBookKeepingCsvExporter: &converters.EzBookKeepingCSVFileExporter{},
		ezBookKeepingTsvExporter: &converters.EzBookKeepingTSVFileExporter{},
		accounts:                 services.Accounts,
		transactions:             services.Transactions,
		categories:               services.TransactionCategories,
		tags:                     services.TransactionTags,
		users:                    services.Users,
		twoFactorAuthorizations:  services.TwoFactorAuthorizations,
		tokens:                   services.Tokens,
		forgetPasswords:          services.ForgetPasswords,
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

	err := l.users.CreateUser(nil, user)

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

	user, err := l.users.GetUserByUsername(nil, username)

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

	user, err := l.users.GetUserByUsername(nil, username)

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

	_, _, err = l.users.UpdateUser(nil, userNew, false)

	if err != nil {
		log.BootErrorf("[user_data.ModifyUserPassword] failed to update user \"%s\" password, because %s", user.Username, err.Error())
		return err
	}

	now := time.Now().Unix()
	err = l.tokens.DeleteTokensBeforeTime(nil, user.Uid, now)

	if err == nil {
		log.BootInfof("[user_data.ModifyUserPassword] revoke old tokens before unix time \"%d\" for user \"%s\"", now, user.Username)
	} else {
		log.BootWarnf("[user_data.ModifyUserPassword] failed to revoke old tokens for user \"%s\", because %s", user.Username, err.Error())
	}

	return nil
}

// SendPasswordResetMail sends an email with password reset link
func (l *UserDataCli) SendPasswordResetMail(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.SendPasswordResetMail] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	user, err := l.users.GetUserByUsername(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.SendPasswordResetMail] failed to get user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	if settings.Container.Current.ForgetPasswordRequireVerifyEmail && !user.EmailVerified {
		log.BootWarnf("[user_data.SendPasswordResetMail] user \"uid:%d\" has not verified email", user.Uid)
		return errs.ErrEmailIsNotVerified
	}

	token, _, err := l.tokens.CreatePasswordResetToken(nil, user)

	if err != nil {
		log.BootErrorf("[user_data.SendPasswordResetMail] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return err
	}

	err = l.forgetPasswords.SendPasswordResetEmail(nil, user, token, "")

	if err != nil {
		log.BootWarnf("[user_data.SendPasswordResetMail] cannot send email to \"%s\", because %s", user.Email, err.Error())
		return err
	}

	return nil
}

// EnableUser sets user enabled according to the specified user name
func (l *UserDataCli) EnableUser(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.EnableUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.EnableUser(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.EnableUser] failed to set user enabled by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// DisableUser sets user disabled according to the specified user name
func (l *UserDataCli) DisableUser(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.DisableUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.DisableUser(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.DisableUser] failed to set user disabled by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// ResendVerifyEmail resends an email with account activation link
func (l *UserDataCli) ResendVerifyEmail(c *cli.Context, username string) error {
	if !settings.Container.Current.EnableUserVerifyEmail {
		return errs.ErrEmailValidationNotAllowed
	}

	if username == "" {
		log.BootErrorf("[user_data.ResendVerifyEmail] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	user, err := l.users.GetUserByUsername(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.ResendVerifyEmail] failed to get user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	if user.EmailVerified {
		log.BootWarnf("[user_data.ResendVerifyEmail] user \"uid:%d\" email has been verified", user.Uid)
		return errs.ErrEmailIsVerified
	}

	token, _, err := l.tokens.CreateEmailVerifyToken(nil, user)

	if err != nil {
		log.BootErrorf("[user_data.ResendVerifyEmail] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return errs.ErrTokenGenerating
	}

	err = l.users.SendVerifyEmail(user, token, "")

	if err != nil {
		log.BootErrorf("[user_data.ResendVerifyEmail] cannot send email to \"%s\", because %s", user.Email, err.Error())
		return err
	}

	return nil
}

// SetUserEmailVerified sets user email address verified
func (l *UserDataCli) SetUserEmailVerified(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.SetUserEmailVerified] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.SetUserEmailVerified(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.SetUserEmailVerified] failed to set user email address verified by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// SetUserEmailUnverified sets user email address unverified
func (l *UserDataCli) SetUserEmailUnverified(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.SetUserEmailUnverified] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.SetUserEmailUnverified(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.SetUserEmailUnverified] failed to set user email address unverified by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// DeleteUser deletes user according to the specified user name
func (l *UserDataCli) DeleteUser(c *cli.Context, username string) error {
	if username == "" {
		log.BootErrorf("[user_data.DeleteUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.DeleteUser(nil, username)

	if err != nil {
		log.BootErrorf("[user_data.DeleteUser] failed to delete user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// ListUserTokens returns all tokens of the specified user
func (l *UserDataCli) ListUserTokens(c *cli.Context, username string) ([]*models.TokenRecord, error) {
	if username == "" {
		log.BootErrorf("[user_data.ListUserTokens] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.ListUserTokens] error occurs when getting user id by user name")
		return nil, err
	}

	tokens, err := l.tokens.GetAllUnexpiredNormalTokensByUid(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.ListUserTokens] failed to get tokens of user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return tokens, nil
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
	err = l.tokens.DeleteTokensBeforeTime(nil, uid, now)

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

	enableTwoFactor, err := l.twoFactorAuthorizations.ExistsTwoFactorSetting(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] failed to check two-factor setting, because %s", err.Error())
		return err
	}

	if !enableTwoFactor {
		return errs.ErrTwoFactorIsNotEnabled
	}

	err = l.twoFactorAuthorizations.DeleteTwoFactorRecoveryCodes(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] failed to delete two-factor recovery codes for user \"%s\"", username)
		return err
	}

	err = l.twoFactorAuthorizations.DeleteTwoFactorSetting(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.DisableUserTwoFactorAuthorization] failed to delete two-factor setting for user \"%s\"", username)
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

	accountMap, categoryMap, tagMap, tagIndexes, tagIndexesMap, err := l.getUserEssentialData(uid, username)

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

	allTransactions, err := l.transactions.GetAllTransactions(nil, uid, pageCountForGettingTransactions, false)

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

		err = l.checkTransactionTag(c, transaction.TransactionId, tagIndexesMap, tagMap)

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
			log.BootErrorf("[user_data.CheckTransactionAndAccount] transaction type of transaction \"id:%d\" is invalid", transaction.TransactionId)
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
			log.BootErrorf("[user_data.CheckTransactionAndAccount] account \"id:%d\" balance is not correct, expected balance is %d, but there is no transaction actually", account.AccountId, account.Balance)
			return false, errs.ErrOperationFailed
		}

		if account.Balance != actualBalance {
			log.BootErrorf("[user_data.CheckTransactionAndAccount] account \"id:%d\" balance is not correct, expected balance is %d, but actual balance is %d", account.AccountId, account.Balance, actualBalance)
			return false, errs.ErrOperationFailed
		}
	}

	for accountId, actualBalance := range accountBalance {
		_, exists := accountMap[accountId]

		if !exists {
			log.BootErrorf("[user_data.CheckTransactionAndAccount] account \"id:%d\" does not exist, but there are some transactions of this account actually, and actual balance is %d", accountId, actualBalance)
			return false, errs.ErrOperationFailed
		}
	}

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]

		if tagIndex.TransactionTime < 1 {
			log.BootErrorf("[user_data.CheckTransactionAndAccount] transaction tag index \"id:%d\" does not have transaction time", tagIndex.TagIndexId)
			return false, errs.ErrOperationFailed
		}
	}

	return true, nil
}

// FixTransactionTagIndexWithTransactionTime fixes user transaction tag index data with transaction time
func (l *UserDataCli) FixTransactionTagIndexWithTransactionTime(c *cli.Context, username string) (bool, error) {
	if username == "" {
		log.BootErrorf("[user_data.FixTransactionTagIndexWithTransactionTime] user name is empty")
		return false, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.FixTransactionTagIndexWithTransactionTime] error occurs when getting user id by user name")
		return false, err
	}

	tagIndexes, err := l.tags.GetAllTagIdsOfAllTransactions(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.FixTransactionTagIndexWithTransactionTime] failed to get tag index for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	invalidTagIndexes := make([]*models.TransactionTagIndex, 0, len(tagIndexes))

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]

		if tagIndex.TransactionTime < 1 {
			invalidTagIndexes = append(invalidTagIndexes, tagIndex)
		}
	}

	if len(invalidTagIndexes) < 1 {
		log.BootErrorf("[user_data.FixTransactionTagIndexWithTransactionTime] all user transaction tag index data has been checked, there is no problem with user data")
		return false, errs.ErrOperationFailed
	}

	allTransactions, err := l.transactions.GetAllTransactions(nil, uid, pageCountForGettingTransactions, false)

	if err != nil {
		log.BootErrorf("[user_data.FixTransactionTagIndexWithTransactionTime] failed to all transactions for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	transactionMap := l.transactions.GetTransactionMapByList(allTransactions)

	for i := 0; i < len(invalidTagIndexes); i++ {
		tagIndex := invalidTagIndexes[i]
		transaction, exists := transactionMap[tagIndex.TransactionId]

		if !exists {
			continue
		}

		tagIndex.TransactionTime = transaction.TransactionTime
	}

	err = l.tags.ModifyTagIndexTransactionTime(nil, uid, invalidTagIndexes)

	if err != nil {
		log.BootErrorf("[user_data.FixTransactionTagIndexWithTransactionTime] failed to update transaction tag index for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	return true, nil
}

// ExportTransaction returns csv file content according user all transactions
func (l *UserDataCli) ExportTransaction(c *cli.Context, username string, fileType string) ([]byte, error) {
	if username == "" {
		log.BootErrorf("[user_data.ExportTransaction] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] error occurs when getting user id by user name")
		return nil, err
	}

	accountMap, categoryMap, tagMap, _, tagIndexesMap, err := l.getUserEssentialData(uid, username)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to get essential data for user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	allTransactions, err := l.transactions.GetAllTransactions(nil, uid, pageCountForDataExport, true)

	if err != nil {
		log.BootErrorf("[user_data.ExportTransaction] failed to all transactions for user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	var dataExporter converters.DataConverter

	if fileType == "tsv" {
		dataExporter = l.ezBookKeepingTsvExporter
	} else {
		dataExporter = l.ezBookKeepingCsvExporter
	}

	result, err := dataExporter.ToExportedContent(uid, allTransactions, accountMap, categoryMap, tagMap, tagIndexesMap)

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

func (l *UserDataCli) getUserEssentialData(uid int64, username string) (accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, tagIndexes []*models.TransactionTagIndex, tagIndexesMap map[int64][]int64, err error) {
	if uid <= 0 {
		log.BootErrorf("[user_data.getUserEssentialData] user uid \"%d\" is invalid", uid)
		return nil, nil, nil, nil, nil, errs.ErrUserIdInvalid
	}

	accounts, err := l.accounts.GetAllAccountsByUid(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get accounts for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	accountMap = l.accounts.GetAccountMapByList(accounts)

	categories, err := l.categories.GetAllCategoriesByUid(nil, uid, 0, -1)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get categories for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	categoryMap = l.categories.GetCategoryMapByList(categories)

	tags, err := l.tags.GetAllTagsByUid(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get tags for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	tagMap = l.tags.GetTagMapByList(tags)

	tagIndexes, err = l.tags.GetAllTagIdsOfAllTransactions(nil, uid)

	if err != nil {
		log.BootErrorf("[user_data.getUserEssentialData] failed to get tag index for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	tagIndexesMap = l.tags.GetGroupedTransactionTagIds(tagIndexes)

	return accountMap, categoryMap, tagMap, tagIndexes, tagIndexesMap, nil
}

func (l *UserDataCli) checkTransactionAccount(c *cli.Context, transaction *models.Transaction, accountMap map[int64]*models.Account, accountHasChild map[int64]bool) error {
	account, exists := accountMap[transaction.AccountId]

	if !exists {
		log.BootErrorf("[user_data.checkTransactionAccount] the account \"id:%d\" of transaction \"id:%d\" does not exist", transaction.AccountId, transaction.TransactionId)
		return errs.ErrAccountNotFound
	}

	if account.ParentAccountId == models.LevelOneAccountParentId && accountHasChild[account.AccountId] {
		log.BootErrorf("[user_data.checkTransactionAccount] the account \"id:%d\" of transaction \"id:%d\" is not a sub-account", transaction.AccountId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountId > 0 {
		relatedAccount, exists := accountMap[transaction.RelatedAccountId]

		if !exists {
			log.BootErrorf("[user_data.checkTransactionAccount] the related account \"id:%d\" of transaction \"id:%d\" does not exist", transaction.RelatedAccountId, transaction.TransactionId)
			return errs.ErrAccountNotFound
		}

		if relatedAccount.ParentAccountId == models.LevelOneAccountParentId && accountHasChild[relatedAccount.AccountId] {
			log.BootErrorf("[user_data.checkTransactionAccount] the related account \"id:%d\" of transaction \"id:%d\" is not a sub-account", transaction.RelatedAccountId, transaction.TransactionId)
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

func (l *UserDataCli) checkTransactionTag(c *cli.Context, transactionId int64, allTagIndexesMap map[int64][]int64, tagMap map[int64]*models.TransactionTag) error {
	tagIndexes, exists := allTagIndexesMap[transactionId]

	if !exists {
		return nil
	}

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]
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

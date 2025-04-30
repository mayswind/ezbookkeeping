package cli

import (
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const pageCountForGettingTransactions = 1000
const pageCountForDataExport = 1000

// UserDataCli represents user data cli
type UserDataCli struct {
	CliUsingConfig
	accounts                *services.AccountService
	transactions            *services.TransactionService
	categories              *services.TransactionCategoryService
	tags                    *services.TransactionTagService
	users                   *services.UserService
	twoFactorAuthorizations *services.TwoFactorAuthorizationService
	tokens                  *services.TokenService
	forgetPasswords         *services.ForgetPasswordService
}

// Initialize a user data cli singleton instance
var (
	UserData = &UserDataCli{
		CliUsingConfig: CliUsingConfig{
			container: settings.Container,
		},
		accounts:                services.Accounts,
		transactions:            services.Transactions,
		categories:              services.TransactionCategories,
		tags:                    services.TransactionTags,
		users:                   services.Users,
		twoFactorAuthorizations: services.TwoFactorAuthorizations,
		tokens:                  services.Tokens,
		forgetPasswords:         services.ForgetPasswords,
	}
)

// AddNewUser adds a new user according to specified info
func (l *UserDataCli) AddNewUser(c *core.CliContext, username string, email string, nickname string, password string, defaultCurrency string) (*models.User, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.AddNewUser] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	if email == "" {
		log.CliErrorf(c, "[user_data.AddNewUser] user email is empty")
		return nil, errs.ErrEmailIsEmpty
	}

	if nickname == "" {
		log.CliErrorf(c, "[user_data.AddNewUser] user nickname is empty")
		return nil, errs.ErrNicknameIsEmpty
	}

	if password == "" {
		log.CliErrorf(c, "[user_data.AddNewUser] user password is empty")
		return nil, errs.ErrPasswordIsEmpty
	}

	if defaultCurrency == "" {
		log.CliErrorf(c, "[user_data.AddNewUser] user default currency is empty")
		return nil, errs.ErrUserDefaultCurrencyIsEmpty
	}

	if _, ok := validators.AllCurrencyNames[defaultCurrency]; !ok {
		log.CliErrorf(c, "[user_data.AddNewUser] user default currency is invalid")
		return nil, errs.ErrUserDefaultCurrencyIsInvalid
	}

	user := &models.User{
		Username:             username,
		Email:                email,
		Nickname:             nickname,
		Password:             password,
		DefaultCurrency:      defaultCurrency,
		FirstDayOfWeek:       core.WEEKDAY_SUNDAY,
		TransactionEditScope: models.TRANSACTION_EDIT_SCOPE_ALL,
		FeatureRestriction:   l.CurrentConfig().DefaultFeatureRestrictions,
	}

	err := l.users.CreateUser(c, user)

	if err != nil {
		log.CliErrorf(c, "[user_data.AddNewUser] failed to create user \"%s\", because %s", user.Username, err.Error())
		return nil, err
	}

	log.CliInfof(c, "[user_data.AddNewUser] user \"%s\" has add successfully, uid is %d", user.Username, user.Uid)

	return user, nil
}

// GetUserByUsername returns user by user name
func (l *UserDataCli) GetUserByUsername(c *core.CliContext, username string) (*models.User, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.GetUserByUsername] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	user, err := l.users.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.GetUserByUsername] failed to get user by user name \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return user, nil
}

// ModifyUserPassword modifies user password
func (l *UserDataCli) ModifyUserPassword(c *core.CliContext, username string, password string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.ModifyUserPassword] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	if password == "" {
		log.CliErrorf(c, "[user_data.ModifyUserPassword] user password is empty")
		return errs.ErrPasswordIsEmpty
	}

	user, err := l.users.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ModifyUserPassword] failed to get user by user name \"%s\", because %s", username, err.Error())
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

	_, _, err = l.users.UpdateUser(c, userNew, false)

	if err != nil {
		log.CliErrorf(c, "[user_data.ModifyUserPassword] failed to update user \"%s\" password, because %s", user.Username, err.Error())
		return err
	}

	now := time.Now().Unix()
	err = l.tokens.DeleteTokensBeforeTime(c, user.Uid, now)

	if err == nil {
		log.CliInfof(c, "[user_data.ModifyUserPassword] revoke old tokens before unix time \"%d\" for user \"%s\"", now, user.Username)
	} else {
		log.CliWarnf(c, "[user_data.ModifyUserPassword] failed to revoke old tokens for user \"%s\", because %s", user.Username, err.Error())
	}

	return nil
}

// SendPasswordResetMail sends an email with password reset link
func (l *UserDataCli) SendPasswordResetMail(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.SendPasswordResetMail] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	user, err := l.users.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.SendPasswordResetMail] failed to get user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	if l.CurrentConfig().ForgetPasswordRequireVerifyEmail && !user.EmailVerified {
		log.CliWarnf(c, "[user_data.SendPasswordResetMail] user \"uid:%d\" has not verified email", user.Uid)
		return errs.ErrEmailIsNotVerified
	}

	token, _, err := l.tokens.CreatePasswordResetTokenWithoutUserAgent(c, user)

	if err != nil {
		log.CliErrorf(c, "[user_data.SendPasswordResetMail] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return err
	}

	err = l.forgetPasswords.SendPasswordResetEmail(c, user, token, "")

	if err != nil {
		log.CliWarnf(c, "[user_data.SendPasswordResetMail] cannot send email to \"%s\", because %s", user.Email, err.Error())
		return err
	}

	return nil
}

// EnableUser sets user enabled according to the specified user name
func (l *UserDataCli) EnableUser(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.EnableUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.EnableUser(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.EnableUser] failed to set user enabled by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// DisableUser sets user disabled according to the specified user name
func (l *UserDataCli) DisableUser(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.DisableUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.DisableUser(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.DisableUser] failed to set user disabled by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// SetUserFeatureRestrictions sets user feature restrictions according to the specified user name
func (l *UserDataCli) SetUserFeatureRestrictions(c *core.CliContext, username string, featureRestriction core.UserFeatureRestrictions) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.SetUserFeatureRestrictions] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.UpdateUserFeatureRestriction(c, username, featureRestriction)

	if err != nil {
		log.CliErrorf(c, "[user_data.SetUserFeatureRestrictions] failed to set user feature restrictions by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// AddUserFeatureRestrictions adds user feature restrictions according to the specified user name
func (l *UserDataCli) AddUserFeatureRestrictions(c *core.CliContext, username string, featureRestriction core.UserFeatureRestrictions) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.AddUserFeatureRestrictions] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.AddUserFeatureRestriction(c, username, featureRestriction)

	if err != nil {
		log.CliErrorf(c, "[user_data.AddUserFeatureRestrictions] failed to add user feature restrictions by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// RemoveUserFeatureRestrictions removes user feature restrictions according to the specified user name
func (l *UserDataCli) RemoveUserFeatureRestrictions(c *core.CliContext, username string, featureRestriction core.UserFeatureRestrictions) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.RemoveUserFeatureRestrictions] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.RemoveUserFeatureRestriction(c, username, featureRestriction)

	if err != nil {
		log.CliErrorf(c, "[user_data.RemoveUserFeatureRestrictions] failed to remove user feature restrictions by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// ResendVerifyEmail resends an email with account activation link
func (l *UserDataCli) ResendVerifyEmail(c *core.CliContext, username string) error {
	if !l.CurrentConfig().EnableUserVerifyEmail {
		return errs.ErrEmailValidationNotAllowed
	}

	if username == "" {
		log.CliErrorf(c, "[user_data.ResendVerifyEmail] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	user, err := l.users.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ResendVerifyEmail] failed to get user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	if user.EmailVerified {
		log.CliWarnf(c, "[user_data.ResendVerifyEmail] user \"uid:%d\" email has been verified", user.Uid)
		return errs.ErrEmailIsVerified
	}

	token, _, err := l.tokens.CreateEmailVerifyTokenWithoutUserAgent(c, user)

	if err != nil {
		log.CliErrorf(c, "[user_data.ResendVerifyEmail] failed to create token for user \"uid:%d\", because %s", user.Uid, err.Error())
		return errs.ErrTokenGenerating
	}

	err = l.users.SendVerifyEmail(user, token, "")

	if err != nil {
		log.CliErrorf(c, "[user_data.ResendVerifyEmail] cannot send email to \"%s\", because %s", user.Email, err.Error())
		return err
	}

	return nil
}

// SetUserEmailVerified sets user email address verified
func (l *UserDataCli) SetUserEmailVerified(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.SetUserEmailVerified] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.SetUserEmailVerified(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.SetUserEmailVerified] failed to set user email address verified by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// SetUserEmailUnverified sets user email address unverified
func (l *UserDataCli) SetUserEmailUnverified(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.SetUserEmailUnverified] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.SetUserEmailUnverified(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.SetUserEmailUnverified] failed to set user email address unverified by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// DeleteUser deletes user according to the specified user name
func (l *UserDataCli) DeleteUser(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.DeleteUser] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	err := l.users.DeleteUser(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.DeleteUser] failed to delete user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// ListUserTokens returns all tokens of the specified user
func (l *UserDataCli) ListUserTokens(c *core.CliContext, username string) ([]*models.TokenRecord, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.ListUserTokens] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ListUserTokens] error occurs when getting user id by user name")
		return nil, err
	}

	tokens, err := l.tokens.GetAllUnexpiredNormalTokensByUid(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.ListUserTokens] failed to get tokens of user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return tokens, nil
}

// CreateNewUserToken returns a new token for the specified user
func (l *UserDataCli) CreateNewUserToken(c *core.CliContext, username string) (*models.TokenRecord, string, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.CreateNewUserToken] user name is empty")
		return nil, "", errs.ErrUsernameIsEmpty
	}

	user, err := l.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.CreateNewUserToken] error occurs when getting user by user name")
		return nil, "", err
	}

	token, tokenRecord, err := l.tokens.CreateTokenViaCli(c, user)

	if err != nil {
		log.CliErrorf(c, "[user_data.CreateNewUserToken] failed to create token for user \"%s\", because %s", username, err.Error())
		return nil, "", err
	}

	return tokenRecord, token, nil
}

// ClearUserTokens clears all tokens of the specified user
func (l *UserDataCli) ClearUserTokens(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.ClearUserTokens] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ClearUserTokens] error occurs when getting user id by user name")
		return err
	}

	now := time.Now().Unix()
	err = l.tokens.DeleteTokensBeforeTime(c, uid, now)

	if err != nil {
		log.CliErrorf(c, "[user_data.ClearUserTokens] failed to delete tokens of user \"%s\", because %s", username, err.Error())
		return err
	}

	return nil
}

// DisableUserTwoFactorAuthorization disables 2fa for the specified user
func (l *UserDataCli) DisableUserTwoFactorAuthorization(c *core.CliContext, username string) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.DisableUserTwoFactorAuthorization] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.DisableUserTwoFactorAuthorization] error occurs when getting user id by user name")
		return err
	}

	enableTwoFactor, err := l.twoFactorAuthorizations.ExistsTwoFactorSetting(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.DisableUserTwoFactorAuthorization] failed to check two-factor setting, because %s", err.Error())
		return err
	}

	if !enableTwoFactor {
		return errs.ErrTwoFactorIsNotEnabled
	}

	err = l.twoFactorAuthorizations.DeleteTwoFactorRecoveryCodes(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.DisableUserTwoFactorAuthorization] failed to delete two-factor recovery codes for user \"%s\"", username)
		return err
	}

	err = l.twoFactorAuthorizations.DeleteTwoFactorSetting(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.DisableUserTwoFactorAuthorization] failed to delete two-factor setting for user \"%s\"", username)
		return err
	}

	return nil
}

// CheckTransactionAndAccount checks whether all user transactions and all user accounts are correct
func (l *UserDataCli) CheckTransactionAndAccount(c *core.CliContext, username string) (bool, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] user name is empty")
		return false, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] error occurs when getting user id by user name")
		return false, err
	}

	accountMap, categoryMap, tagMap, tagIndexes, tagIndexesMap, err := l.getUserEssentialData(c, uid, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] failed to get essential data for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	accountHasChild := make(map[int64]bool)

	for _, account := range accountMap {
		if account.ParentAccountId > models.LevelOneAccountParentId {
			accountHasChild[account.ParentAccountId] = true
		}
	}

	allTransactions, err := l.transactions.GetAllTransactions(c, uid, pageCountForGettingTransactions, false)

	if err != nil {
		log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] failed to all transactions for user \"%s\", because %s", username, err.Error())
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
			log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] transaction type of transaction \"id:%d\" is invalid", transaction.TransactionId)
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
			log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] account \"id:%d\" balance is not correct, expected balance is %d, but there is no transaction actually", account.AccountId, account.Balance)
			return false, errs.ErrOperationFailed
		}

		if account.Balance != actualBalance {
			log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] account \"id:%d\" balance is not correct, expected balance is %d, but actual balance is %d", account.AccountId, account.Balance, actualBalance)
			return false, errs.ErrOperationFailed
		}
	}

	for accountId, actualBalance := range accountBalance {
		_, exists := accountMap[accountId]

		if !exists {
			log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] account \"id:%d\" does not exist, but there are some transactions of this account actually, and actual balance is %d", accountId, actualBalance)
			return false, errs.ErrOperationFailed
		}
	}

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]

		if tagIndex.TransactionTime < 1 {
			log.CliErrorf(c, "[user_data.CheckTransactionAndAccount] transaction tag index \"id:%d\" does not have transaction time", tagIndex.TagIndexId)
			return false, errs.ErrOperationFailed
		}
	}

	return true, nil
}

// FixTransactionTagIndexWithTransactionTime fixes user transaction tag index data with transaction time
func (l *UserDataCli) FixTransactionTagIndexWithTransactionTime(c *core.CliContext, username string) (bool, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.FixTransactionTagIndexWithTransactionTime] user name is empty")
		return false, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.FixTransactionTagIndexWithTransactionTime] error occurs when getting user id by user name")
		return false, err
	}

	tagIndexes, err := l.tags.GetAllTagIdsOfAllTransactions(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.FixTransactionTagIndexWithTransactionTime] failed to get tag index for user \"%s\", because %s", username, err.Error())
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
		log.CliErrorf(c, "[user_data.FixTransactionTagIndexWithTransactionTime] all user transaction tag index data has been checked, there is no problem with user data")
		return false, errs.ErrOperationFailed
	}

	allTransactions, err := l.transactions.GetAllTransactions(c, uid, pageCountForGettingTransactions, false)

	if err != nil {
		log.CliErrorf(c, "[user_data.FixTransactionTagIndexWithTransactionTime] failed to all transactions for user \"%s\", because %s", username, err.Error())
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

	err = l.tags.ModifyTagIndexTransactionTime(c, uid, invalidTagIndexes)

	if err != nil {
		log.CliErrorf(c, "[user_data.FixTransactionTagIndexWithTransactionTime] failed to update transaction tag index for user \"%s\", because %s", username, err.Error())
		return false, err
	}

	return true, nil
}

// ExportTransaction returns csv file content according user all transactions
func (l *UserDataCli) ExportTransaction(c *core.CliContext, username string, fileType string) ([]byte, error) {
	if username == "" {
		log.CliErrorf(c, "[user_data.ExportTransaction] user name is empty")
		return nil, errs.ErrUsernameIsEmpty
	}

	uid, err := l.getUserIdByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ExportTransaction] error occurs when getting user id by user name")
		return nil, err
	}

	accountMap, categoryMap, tagMap, _, tagIndexesMap, err := l.getUserEssentialData(c, uid, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ExportTransaction] failed to get essential data for user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	allTransactions, err := l.transactions.GetAllTransactions(c, uid, pageCountForDataExport, true)

	if err != nil {
		log.CliErrorf(c, "[user_data.ExportTransaction] failed to all transactions for user \"%s\", because %s", username, err.Error())
		return nil, err
	}

	dataExporter := converters.GetTransactionDataExporter(fileType)

	if dataExporter == nil {
		return nil, errs.ErrNotImplemented
	}

	result, err := dataExporter.ToExportedContent(c, uid, allTransactions, accountMap, categoryMap, tagMap, tagIndexesMap)

	if err != nil {
		log.CliErrorf(c, "[user_data.ExportTransaction] failed to get csv format exported data for \"%s\", because %s", username, err.Error())
		return nil, err
	}

	return result, nil
}

func (l *UserDataCli) ImportTransaction(c *core.CliContext, username string, fileType string, data []byte) error {
	if username == "" {
		log.CliErrorf(c, "[user_data.ImportTransaction] user name is empty")
		return errs.ErrUsernameIsEmpty
	}

	dataImporter, err := converters.GetTransactionDataImporter(fileType)

	if err != nil {
		return err
	}

	user, err := l.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ImportTransaction] failed to get user by user name \"%s\", because %s", username, err.Error())
		return err
	}

	accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap, err := l.getUserEssentialDataForImport(c, user.Uid, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.ImportTransaction] failed to get essential data for user \"%s\", because %s", username, err.Error())
		return err
	}

	parsedTransactions, newAccounts, newSubExpenseCategories, newSubIncomeCategories, newSubTransferCategories, newTags, err := dataImporter.ParseImportedData(c, user, data, utils.GetTimezoneOffsetMinutes(time.Local), accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)

	if err != nil {
		log.CliErrorf(c, "[user_data.ImportTransaction] failed to parse imported data for \"%s\", because %s", username, err.Error())
		return err
	}

	if len(parsedTransactions) < 1 {
		log.CliErrorf(c, "[user_data.ImportTransaction] there are no transactions in import file")
		return errs.ErrOperationFailed
	}

	if len(newAccounts) > 0 {
		accountNames := l.accounts.GetAccountNames(newAccounts)
		log.CliErrorf(c, "[user_data.ImportTransaction] there are %d accounts (%s) need to be created, please create them manually", len(newAccounts), strings.Join(accountNames, ","))
		return errs.ErrOperationFailed
	}

	if len(newSubExpenseCategories) > 0 {
		categoryNames := l.categories.GetCategoryNames(newSubExpenseCategories)
		log.CliErrorf(c, "[user_data.ImportTransaction] there are %d expense categories (%s) need to be created, please create them manually", len(newSubExpenseCategories), strings.Join(categoryNames, ","))
		return errs.ErrOperationFailed
	}

	if len(newSubIncomeCategories) > 0 {
		categoryNames := l.categories.GetCategoryNames(newSubIncomeCategories)
		log.CliErrorf(c, "[user_data.ImportTransaction] there are %d income categories (%s) need to be created, please create them manually", len(newSubIncomeCategories), strings.Join(categoryNames, ","))
		return errs.ErrOperationFailed
	}

	if len(newSubTransferCategories) > 0 {
		categoryNames := l.categories.GetCategoryNames(newSubTransferCategories)
		log.CliErrorf(c, "[user_data.ImportTransaction] there are %d transfer categories (%s) need to be created, please create them manually", len(newSubTransferCategories), strings.Join(categoryNames, ","))
		return errs.ErrOperationFailed
	}

	if len(newTags) > 0 {
		tagNames := l.tags.GetTagNames(newTags)
		log.CliErrorf(c, "[user_data.ImportTransaction] there are %d transaction tags (%s) need to be created, please create them manually", len(newTags), strings.Join(tagNames, ","))
		return errs.ErrOperationFailed
	}

	newTransactions := parsedTransactions.ToTransactionsList()
	newTransactionTagIdsMap, err := parsedTransactions.ToTransactionTagIdsMap()

	if err != nil {
		log.CliErrorf(c, "[user_data.ImportTransaction] failed to get transaction tag ids map, because %s", err.Error())
		return errs.ErrOperationFailed
	}

	err = l.transactions.BatchCreateTransactions(c, user.Uid, newTransactions, newTransactionTagIdsMap, nil)

	if err != nil {
		log.CliErrorf(c, "[user_data.ImportTransaction] failed to create transaction, because %s", err.Error())
		return err
	}

	return nil
}

func (l *UserDataCli) getUserIdByUsername(c *core.CliContext, username string) (int64, error) {
	user, err := l.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserIdByUsername] failed to get user by user name \"%s\", because %s", username, err.Error())
		return 0, err
	}

	return user.Uid, nil
}

func (l *UserDataCli) getUserEssentialData(c *core.CliContext, uid int64, username string) (accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, tagIndexes []*models.TransactionTagIndex, tagIndexesMap map[int64][]int64, err error) {
	if uid <= 0 {
		log.CliErrorf(c, "[user_data.getUserEssentialData] user uid \"%d\" is invalid", uid)
		return nil, nil, nil, nil, nil, errs.ErrUserIdInvalid
	}

	accounts, err := l.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialData] failed to get accounts for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	accountMap = l.accounts.GetAccountMapByList(accounts)

	categories, err := l.categories.GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialData] failed to get categories for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	categoryMap = l.categories.GetCategoryMapByList(categories)

	tags, err := l.tags.GetAllTagsByUid(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialData] failed to get tags for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	tagMap = l.tags.GetTagMapByList(tags)

	tagIndexes, err = l.tags.GetAllTagIdsOfAllTransactions(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialData] failed to get tag index for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	tagIndexesMap = l.tags.GetGroupedTransactionTagIds(tagIndexes)

	return accountMap, categoryMap, tagMap, tagIndexes, tagIndexesMap, nil
}

func (l *UserDataCli) getUserEssentialDataForImport(c *core.CliContext, uid int64, username string) (accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag, err error) {
	if uid <= 0 {
		log.CliErrorf(c, "[user_data.getUserEssentialDataForImport] user uid \"%d\" is invalid", uid)
		return nil, nil, nil, nil, nil, errs.ErrUserIdInvalid
	}

	accounts, err := l.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialDataForImport] failed to get accounts for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	accountMap = l.accounts.GetVisibleAccountNameMapByList(accounts)

	categories, err := l.categories.GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialDataForImport] failed to get categories for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	expenseCategoryMap, incomeCategoryMap, transferCategoryMap = l.categories.GetVisibleSubCategoryNameMapByList(categories)

	tags, err := l.tags.GetAllTagsByUid(c, uid)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserEssentialDataForImport] failed to get tags for user \"%s\", because %s", username, err.Error())
		return nil, nil, nil, nil, nil, err
	}

	tagMap = l.tags.GetTagNameMapByList(tags)

	return accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap, nil
}

func (l *UserDataCli) checkTransactionAccount(c *core.CliContext, transaction *models.Transaction, accountMap map[int64]*models.Account, accountHasChild map[int64]bool) error {
	account, exists := accountMap[transaction.AccountId]

	if !exists {
		log.CliErrorf(c, "[user_data.checkTransactionAccount] the account \"id:%d\" of transaction \"id:%d\" does not exist", transaction.AccountId, transaction.TransactionId)
		return errs.ErrAccountNotFound
	}

	if account.ParentAccountId == models.LevelOneAccountParentId && accountHasChild[account.AccountId] {
		log.CliErrorf(c, "[user_data.checkTransactionAccount] the account \"id:%d\" of transaction \"id:%d\" is not a sub-account", transaction.AccountId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountId > 0 {
		relatedAccount, exists := accountMap[transaction.RelatedAccountId]

		if !exists {
			log.CliErrorf(c, "[user_data.checkTransactionAccount] the related account \"id:%d\" of transaction \"id:%d\" does not exist", transaction.RelatedAccountId, transaction.TransactionId)
			return errs.ErrAccountNotFound
		}

		if relatedAccount.ParentAccountId == models.LevelOneAccountParentId && accountHasChild[relatedAccount.AccountId] {
			log.CliErrorf(c, "[user_data.checkTransactionAccount] the related account \"id:%d\" of transaction \"id:%d\" is not a sub-account", transaction.RelatedAccountId, transaction.TransactionId)
			return errs.ErrOperationFailed
		}
	}

	return nil
}

func (l *UserDataCli) checkTransactionCategory(c *core.CliContext, transaction *models.Transaction, categoryMap map[int64]*models.TransactionCategory) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.CategoryId > 0 {
			log.CliErrorf(c, "[user_data.checkTransactionCategory] transaction \"id:%d\" is balance modification transaction, but has category \"id:%d\"", transaction.TransactionId, transaction.CategoryId)
			return errs.ErrBalanceModificationTransactionCannotSetCategory
		} else {
			return nil
		}
	}

	category, exists := categoryMap[transaction.CategoryId]

	if !exists {
		log.CliErrorf(c, "[user_data.checkTransactionCategory] the transaction category \"id:%d\" of transaction \"id:%d\" does not exist", transaction.CategoryId, transaction.TransactionId)
		return errs.ErrTransactionCategoryNotFound
	}

	if category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
		log.CliErrorf(c, "[user_data.checkTransactionCategory] the transaction category \"id:%d\" of transaction \"id:%d\" is not a sub category", transaction.CategoryId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	return nil
}

func (l *UserDataCli) checkTransactionTag(c *core.CliContext, transactionId int64, allTagIndexesMap map[int64][]int64, tagMap map[int64]*models.TransactionTag) error {
	tagIndexes, exists := allTagIndexesMap[transactionId]

	if !exists {
		return nil
	}

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]
		tag, exists := tagMap[tagIndex]

		if !exists {
			log.CliErrorf(c, "[user_data.checkTransactionTag] the transaction tag \"id:%d\" of transaction \"id:%d\" does not exist", tag.TagId, transactionId)
			return errs.ErrTransactionTagNotFound
		}
	}

	return nil
}

func (l *UserDataCli) checkTransactionRelatedTransaction(c *core.CliContext, transaction *models.Transaction, transactionMap map[int64]*models.Transaction, accountMap map[int64]*models.Account) error {
	if transaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_OUT && transaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return nil
	}

	relatedTransaction, exists := transactionMap[transaction.RelatedId]

	if !exists {
		log.CliErrorf(c, "[user_data.checkTransactionRelatedTransaction] the related transaction \"id:%d\" of transaction \"id:%d\" does not exist", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrTransactionNotFound
	}

	if transaction.RelatedId != relatedTransaction.TransactionId || transaction.TransactionId != relatedTransaction.RelatedId {
		log.CliErrorf(c, "[user_data.checkTransactionRelatedTransaction] related ids of transaction \"id:%d\" and transaction \"id:%d\" are not equal", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountId != relatedTransaction.AccountId || transaction.AccountId != relatedTransaction.RelatedAccountId {
		log.CliErrorf(c, "[user_data.checkTransactionRelatedTransaction] related account ids of transaction \"id:%d\" and transaction \"id:%d\" are not equal", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	if transaction.RelatedAccountAmount != relatedTransaction.Amount || transaction.Amount != relatedTransaction.RelatedAccountAmount {
		log.CliErrorf(c, "[user_data.checkTransactionRelatedTransaction] related amounts of transaction \"id:%d\" and transaction \"id:%d\" are not equal", transaction.RelatedId, transaction.TransactionId)
		return errs.ErrOperationFailed
	}

	account := accountMap[transaction.AccountId]
	relatedAccount := accountMap[transaction.RelatedAccountId]

	if account.Currency == relatedAccount.Currency && transaction.Amount != transaction.RelatedAccountAmount {
		log.CliWarnf(c, "[user_data.checkTransactionRelatedTransaction] transfer-in amount and transfer-out amount of transaction \"id:%d\" are not equal", transaction.TransactionId)
	}

	return nil
}

package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	clis "github.com/mayswind/ezbookkeeping/pkg/cli"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// UserData represents the data command
var UserData = &cli.Command{
	Name:  "userdata",
	Usage: "ezBookkeeping user data maintenance",
	Subcommands: []*cli.Command{
		{
			Name:   "user-add",
			Usage:  "Add new user",
			Action: bindAction(addNewUser),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "New user name",
				},
				&cli.StringFlag{
					Name:     "email",
					Aliases:  []string{"m"},
					Required: true,
					Usage:    "New user email",
				},
				&cli.StringFlag{
					Name:     "nickname",
					Aliases:  []string{"i"},
					Required: true,
					Usage:    "New user nickname",
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"p"},
					Required: true,
					Usage:    "New user password",
				},
				&cli.StringFlag{
					Name:     "default-currency",
					Aliases:  []string{"c"},
					Required: true,
					Usage:    "New user default currency",
				},
			},
		},
		{
			Name:   "user-get",
			Usage:  "Get specified user info",
			Action: bindAction(getUserInfo),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-modify-password",
			Usage:  "Modify user password",
			Action: bindAction(modifyUserPassword),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"p"},
					Required: true,
					Usage:    "User new password",
				},
			},
		},
		{
			Name:   "user-enable",
			Usage:  "Enable specified user",
			Action: bindAction(enableUser),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-disable",
			Usage:  "Disable specified user",
			Action: bindAction(disableUser),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-set-restrict-features",
			Usage:  "Set restrictions of user features",
			Action: bindAction(setUserFeatureRestriction),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
				&cli.StringFlag{
					Name:     "features",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "Specific feature types (feature types separated by commas)",
				},
			},
		},
		{
			Name:   "user-add-restrict-features",
			Usage:  "Add restrictions of user features",
			Action: bindAction(addUserFeatureRestriction),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
				&cli.StringFlag{
					Name:     "features",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "Specific feature types (feature types separated by commas)",
				},
			},
		},
		{
			Name:   "user-remove-restrict-features",
			Usage:  "Remove restrictions of user features",
			Action: bindAction(removeUserFeatureRestriction),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
				&cli.StringFlag{
					Name:     "features",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "Specific feature types (feature types separated by commas)",
				},
			},
		},
		{
			Name:   "user-resend-verify-email",
			Usage:  "Resend user verify email",
			Action: bindAction(resendUserVerifyEmail),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-set-email-verified",
			Usage:  "Set user email address verified",
			Action: bindAction(setUserEmailVerified),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-set-email-unverified",
			Usage:  "Set user email address unverified",
			Action: bindAction(setUserEmailUnverified),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-delete",
			Usage:  "Delete specified user",
			Action: bindAction(deleteUser),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-2fa-disable",
			Usage:  "Disable user 2fa setting",
			Action: bindAction(disableUser2FA),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-session-list",
			Usage:  "List all user sessions",
			Action: bindAction(listUserTokens),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-session-new",
			Usage:  "Create new session for user",
			Action: bindAction(createNewUserToken),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "user-session-clear",
			Usage:  "Clear user all sessions",
			Action: bindAction(clearUserTokens),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "send-password-reset-mail",
			Usage:  "Send password reset mail",
			Action: bindAction(sendPasswordResetMail),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "transaction-check",
			Usage:  "Check whether user all transactions and accounts are correct",
			Action: bindAction(checkUserTransactionAndAccount),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "transaction-tag-index-fix-transaction-time",
			Usage:  "Fix the transaction tag index data which does not have transaction time",
			Action: bindAction(fixTransactionTagIndexNotHaveTransactionTime),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
			},
		},
		{
			Name:   "transaction-import",
			Usage:  "Import transactions to specified user",
			Action: bindAction(importUserTransaction),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
				&cli.StringFlag{
					Name:     "file",
					Aliases:  []string{"f"},
					Required: true,
					Usage:    "Specific import file path (e.g. transaction.csv)",
				},
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "Import file type (supports \"ezbookkeeping_csv\", \"ezbookkeeping_tsv\")",
				},
			},
		},
		{
			Name:   "transaction-export",
			Usage:  "Export user all transactions to file",
			Action: bindAction(exportUserTransaction),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Specific user name",
				},
				&cli.StringFlag{
					Name:     "file",
					Aliases:  []string{"f"},
					Required: true,
					Usage:    "Specific exported file path (e.g. transaction.csv)",
				},
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Required: false,
					Usage:    "Export file type, support csv or tsv, default is csv",
				},
			},
		},
	},
}

func addNewUser(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	email := c.String("email")
	nickname := c.String("nickname")
	password := c.String("password")
	defaultCurrency := c.String("default-currency")

	user, err := clis.UserData.AddNewUser(c, username, email, nickname, password, defaultCurrency)

	if err != nil {
		log.CliErrorf(c, "[user_data.addNewUser] error occurs when adding new user")
		return err
	}

	printUserInfo(user)

	return nil
}

func getUserInfo(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	user, err := clis.UserData.GetUserByUsername(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.getUserInfo] error occurs when getting user data")
		return err
	}

	printUserInfo(user)

	return nil
}

func modifyUserPassword(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	password := c.String("password")
	err = clis.UserData.ModifyUserPassword(c, username, password)

	if err != nil {
		log.CliErrorf(c, "[user_data.modifyUserPassword] error occurs when modifying user password")
		return err
	}

	log.CliInfof(c, "[user_data.modifyUserPassword] password of user \"%s\" has been changed", username)

	return nil
}

func sendPasswordResetMail(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.SendPasswordResetMail(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.sendPasswordResetMail] error occurs when sending password reset email")
		return err
	}

	log.CliInfof(c, "[user_data.sendPasswordResetMail] a password reset email for user \"%s\" has been sent", username)

	return nil
}

func enableUser(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.EnableUser(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.enableUser] error occurs when setting user enabled")
		return err
	}

	log.CliInfof(c, "[user_data.enableUser] user \"%s\" has been set enabled", username)

	return nil
}

func disableUser(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.DisableUser(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.disableUser] error occurs when setting user disabled")
		return err
	}

	log.CliInfof(c, "[user_data.disableUser] user \"%s\" has been set disabled", username)

	return nil
}

func setUserFeatureRestriction(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	featureRestriction := core.ParseUserFeatureRestrictions(c.String("features"))
	err = clis.UserData.SetUserFeatureRestrictions(c, username, featureRestriction)

	if err != nil {
		log.CliErrorf(c, "[user_data.setUserFeatureRestriction] error occurs when setting user feature restriction")
		return err
	}

	log.CliInfof(c, "[user_data.setUserFeatureRestriction] user \"%s\" has been set new feature restriction", username)

	return nil
}

func addUserFeatureRestriction(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	featureRestriction := core.ParseUserFeatureRestrictions(c.String("features"))

	if featureRestriction < 1 {
		log.CliErrorf(c, "[user_data.addUserFeatureRestriction] nothing has been modified")
		return nil
	}

	err = clis.UserData.AddUserFeatureRestrictions(c, username, featureRestriction)

	if err != nil {
		log.CliErrorf(c, "[user_data.addUserFeatureRestriction] error occurs when adding user feature restriction")
		return err
	}

	log.CliInfof(c, "[user_data.addUserFeatureRestriction] user \"%s\" has been add new feature restriction", username)

	return nil
}

func removeUserFeatureRestriction(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	featureRestriction := core.ParseUserFeatureRestrictions(c.String("features"))

	if featureRestriction < 1 {
		log.CliErrorf(c, "[user_data.removeUserFeatureRestriction] nothing has been modified")
		return nil
	}

	err = clis.UserData.RemoveUserFeatureRestrictions(c, username, featureRestriction)

	if err != nil {
		log.CliErrorf(c, "[user_data.removeUserFeatureRestriction] error occurs when removing user feature restriction")
		return err
	}

	log.CliInfof(c, "[user_data.removeUserFeatureRestriction] user \"%s\" has been removed new feature restriction", username)

	return nil
}

func resendUserVerifyEmail(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.ResendVerifyEmail(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.resendUserVerifyEmail] error occurs when resending user verify email")
		return err
	}

	log.CliInfof(c, "[user_data.resendUserVerifyEmail] verify email for user \"%s\" has been resent", username)

	return nil
}

func setUserEmailVerified(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.SetUserEmailVerified(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.setUserEmailVerified] error occurs when setting user email address verified")
		return err
	}

	log.CliInfof(c, "[user_data.setUserEmailVerified] user \"%s\" email address has been set verified", username)

	return nil
}

func setUserEmailUnverified(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.SetUserEmailUnverified(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.setUserEmailUnverified] error occurs when setting user email address unverified")
		return err
	}

	log.CliInfof(c, "[user_data.setUserEmailUnverified] user \"%s\" email address has been set unverified", username)

	return nil
}

func deleteUser(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.DeleteUser(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.deleteUser] error occurs when deleting user")
		return err
	}

	log.CliInfof(c, "[user_data.deleteUser] user \"%s\" has been deleted", username)

	return nil
}

func disableUser2FA(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.DisableUserTwoFactorAuthorization(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.disableUser2FA] error occurs when disabling user two-factor authorization")
		return err
	}

	log.CliInfof(c, "[user_data.disableUser2FA] two-factor authorization of user \"%s\" has been disabled", username)

	return nil
}

func listUserTokens(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	tokens, err := clis.UserData.ListUserTokens(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.listUserTokens] error occurs when getting user tokens")
		return err
	}

	for i := 0; i < len(tokens); i++ {
		printTokenInfo(tokens[i])

		if i < len(tokens)-1 {
			fmt.Printf("---\n")
		}
	}

	return nil
}

func createNewUserToken(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	token, tokenString, err := clis.UserData.CreateNewUserToken(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.createNewUserToken] error occurs when creating user token")
		return err
	}

	printTokenInfo(token)
	fmt.Printf("[NewToken] %s\n", tokenString)

	return nil
}

func clearUserTokens(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.ClearUserTokens(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.clearUserTokens] error occurs when clearing user tokens")
		return err
	}

	log.CliInfof(c, "[user_data.clearUserTokens] all tokens of user \"%s\" has been cleared", username)

	return nil
}

func checkUserTransactionAndAccount(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")

	log.CliInfof(c, "[user_data.checkUserTransactionAndAccount] starting checking user \"%s\" data", username)

	_, err = clis.UserData.CheckTransactionAndAccount(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.checkUserTransactionAndAccount] error occurs when checking user data")
		return err
	}

	log.CliInfof(c, "[user_data.checkUserTransactionAndAccount] user transactions and accounts data has been checked successfully, there is no problem with user data")

	return nil
}

func fixTransactionTagIndexNotHaveTransactionTime(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")

	log.CliInfof(c, "[user_data.fixTransactionTagIndexNotHaveTransactionTime] starting fixing user \"%s\" transaction tag index data", username)

	_, err = clis.UserData.FixTransactionTagIndexWithTransactionTime(c, username)

	if err != nil {
		log.CliErrorf(c, "[user_data.fixTransactionTagIndexNotHaveTransactionTime] error occurs when fixing user data")
		return err
	}

	log.CliInfof(c, "[user_data.fixTransactionTagIndexNotHaveTransactionTime] user transaction tag index data has been fixed successfully")

	return nil
}

func exportUserTransaction(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	filePath := c.String("file")
	fileType := c.String("type")

	if fileType == "" {
		fileType = "csv"
	}

	if fileType != "csv" && fileType != "tsv" {
		log.CliErrorf(c, "[user_data.exportUserTransaction] export file type is not supported")
		return errs.ErrNotSupported
	}

	if filePath == "" {
		log.CliErrorf(c, "[user_data.exportUserTransaction] export file path is unspecified")
		return os.ErrNotExist
	}

	fileExists, err := utils.IsExists(filePath)

	if fileExists {
		log.CliErrorf(c, "[user_data.exportUserTransaction] specified file path already exists")
		return os.ErrExist
	}

	log.CliInfof(c, "[user_data.exportUserTransaction] starting exporting user \"%s\" data", username)

	content, err := clis.UserData.ExportTransaction(c, username, fileType)

	if err != nil {
		log.CliErrorf(c, "[user_data.exportUserTransaction] error occurs when exporting user data")
		return err
	}

	err = utils.WriteFile(filePath, content)

	if err != nil {
		log.CliErrorf(c, "[user_data.exportUserTransaction] failed to write to %s", filePath)
		return err
	}

	log.CliInfof(c, "[user_data.exportUserTransaction] user transactions have been exported to %s", filePath)

	return nil
}

func importUserTransaction(c *core.CliContext) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	filePath := c.String("file")
	filetype := c.String("type")

	if filePath == "" {
		log.CliErrorf(c, "[user_data.importUserTransaction] import file path is not specified")
		return os.ErrNotExist
	}

	fileExists, err := utils.IsExists(filePath)

	if !fileExists {
		log.CliErrorf(c, "[user_data.importUserTransaction] import file does not exist")
		return os.ErrExist
	}

	if filetype != "ezbookkeeping_csv" && filetype != "ezbookkeeping_tsv" {
		log.CliErrorf(c, "[user_data.importUserTransaction] unknown file type \"%s\"", filetype)
		return errs.ErrImportFileTypeNotSupported
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		log.CliErrorf(c, "[user_data.importUserTransaction] failed to load import file")
		return err
	}

	log.CliInfof(c, "[user_data.importUserTransaction] start importing transactions to user \"%s\"", username)

	err = clis.UserData.ImportTransaction(c, username, filetype, data)

	if err != nil {
		log.CliErrorf(c, "[user_data.importUserTransaction] error occurs when importing user data")
		return err
	}

	log.CliInfof(c, "[user_data.importUserTransaction] transactions have been imported to user \"%s\"", username)

	return nil
}

func printUserInfo(user *models.User) {
	fmt.Printf("[Uid] %d\n", user.Uid)
	fmt.Printf("[Username] %s\n", user.Username)
	fmt.Printf("[Email] %s\n", user.Email)
	fmt.Printf("[Nickname] %s\n", user.Nickname)
	fmt.Printf("[Password] %s\n", user.Password)
	fmt.Printf("[Salt] %s\n", user.Salt)
	fmt.Printf("[DefaultAccountId] %d\n", user.DefaultAccountId)
	fmt.Printf("[TransactionEditScope] %s (%d)\n", user.TransactionEditScope, user.TransactionEditScope)
	fmt.Printf("[Language] %s\n", user.Language)
	fmt.Printf("[DefaultCurrency] %s\n", user.DefaultCurrency)
	fmt.Printf("[FirstDayOfWeek] %s (%d)\n", user.FirstDayOfWeek, user.FirstDayOfWeek)
	fmt.Printf("[FiscalYearStart] %s (%d)\n", user.FiscalYearStart, user.FiscalYearStart)
	fmt.Printf("[LongDateFormat] %s (%d)\n", user.LongDateFormat, user.LongDateFormat)
	fmt.Printf("[ShortDateFormat] %s (%d)\n", user.ShortDateFormat, user.ShortDateFormat)
	fmt.Printf("[LongTimeFormat] %s (%d)\n", user.LongTimeFormat, user.LongTimeFormat)
	fmt.Printf("[ShortTimeFormat] %s (%d)\n", user.ShortTimeFormat, user.ShortTimeFormat)
	fmt.Printf("[DecimalSeparator] %s (%d)\n", user.DecimalSeparator, user.DecimalSeparator)
	fmt.Printf("[DigitGroupingSymbol] %s (%d)\n", user.DigitGroupingSymbol, user.DigitGroupingSymbol)
	fmt.Printf("[DigitGrouping] %s (%d)\n", user.DigitGrouping, user.DigitGrouping)
	fmt.Printf("[CurrencyDisplayType] %s (%d)\n", user.CurrencyDisplayType, user.CurrencyDisplayType)
	fmt.Printf("[ExpenseAmountColor] %s (%d)\n", user.ExpenseAmountColor, user.ExpenseAmountColor)
	fmt.Printf("[IncomeAmountColor] %s (%d)\n", user.IncomeAmountColor, user.IncomeAmountColor)
	fmt.Printf("[FeatureRestriction] %s (%d)\n", user.FeatureRestriction, user.FeatureRestriction)
	fmt.Printf("[Deleted] %t\n", user.Deleted)
	fmt.Printf("[EmailVerified] %t\n", user.EmailVerified)
	fmt.Printf("[CreatedAt] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(user.CreatedUnixTime), user.CreatedUnixTime)

	if user.UpdatedUnixTime > 0 {
		fmt.Printf("[UpdatedAt] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(user.UpdatedUnixTime), user.UpdatedUnixTime)
	}

	if user.DeletedUnixTime > 0 {
		fmt.Printf("[DeletedAt] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(user.DeletedUnixTime), user.DeletedUnixTime)
	}

	if user.LastLoginUnixTime > 0 {
		fmt.Printf("[LastLoginAt] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(user.LastLoginUnixTime), user.LastLoginUnixTime)
	}
}

func printTokenInfo(token *models.TokenRecord) {
	fmt.Printf("[CreatedAt] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(token.CreatedUnixTime), token.CreatedUnixTime)
	fmt.Printf("[ExpiredAt] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(token.ExpiredUnixTime), token.ExpiredUnixTime)
	fmt.Printf("[LastSeen] %s (%d)\n", utils.FormatUnixTimeToLongDateTimeInServerTimezone(token.LastSeenUnixTime), token.LastSeenUnixTime)
	fmt.Printf("[UserAgent] %s\n", token.UserAgent)
}

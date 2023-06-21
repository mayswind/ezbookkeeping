package cmd

import (
	"fmt"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"os"

	"github.com/urfave/cli/v2"

	clis "github.com/mayswind/ezbookkeeping/pkg/cli"
	"github.com/mayswind/ezbookkeeping/pkg/log"
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
			Action: addNewUser,
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
			Action: getUserInfo,
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
			Action: modifyUserPassword,
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
			Action: enableUser,
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
			Action: disableUser,
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
			Action: deleteUser,
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
			Action: disableUser2FA,
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
			Action: listUserTokens,
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
			Action: clearUserTokens,
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
			Action: checkUserTransactionAndAccount,
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
			Name:   "transaction-export",
			Usage:  "Export user all transactions to csv file",
			Action: exportUserTransaction,
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
			},
		},
	},
}

func addNewUser(c *cli.Context) error {
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
		log.BootErrorf("[user_data.addNewUser] error occurs when adding new user")
		return err
	}

	printUserInfo(user)

	return nil
}

func getUserInfo(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	user, err := clis.UserData.GetUserByUsername(c, username)

	if err != nil {
		log.BootErrorf("[user_data.getUserInfo] error occurs when getting user data")
		return err
	}

	printUserInfo(user)

	return nil
}

func modifyUserPassword(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	password := c.String("password")
	err = clis.UserData.ModifyUserPassword(c, username, password)

	if err != nil {
		log.BootErrorf("[user_data.modifyUserPassword] error occurs when modifying user password")
		return err
	}

	log.BootInfof("[user_data.modifyUserPassword] password of user \"%s\" has been changed", username)

	return nil
}

func enableUser(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.EnableUser(c, username)

	if err != nil {
		log.BootErrorf("[user_data.enableUser] error occurs when setting user enabled")
		return err
	}

	log.BootInfof("[user_data.enableUser] user \"%s\" has been set enabled", username)

	return nil
}

func disableUser(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.DisableUser(c, username)

	if err != nil {
		log.BootErrorf("[user_data.disableUser] error occurs when setting user disabled")
		return err
	}

	log.BootInfof("[user_data.disableUser] user \"%s\" has been set disabled", username)

	return nil
}

func deleteUser(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.DeleteUser(c, username)

	if err != nil {
		log.BootErrorf("[user_data.deleteUser] error occurs when deleting user")
		return err
	}

	log.BootInfof("[user_data.deleteUser] user \"%s\" has been deleted", username)

	return nil
}

func disableUser2FA(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.DisableUserTwoFactorAuthorization(c, username)

	if err != nil {
		log.BootErrorf("[user_data.disableUser2FA] error occurs when disabling user two factor authorization")
		return err
	}

	log.BootInfof("[user_data.disableUser2FA] two factor authorization of user \"%s\" has been disabled", username)

	return nil
}

func listUserTokens(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	tokens, err := clis.UserData.ListUserTokens(c, username)

	if err != nil {
		log.BootErrorf("[user_data.listUserTokens] error occurs when getting user tokens")
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

func clearUserTokens(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	err = clis.UserData.ClearUserTokens(c, username)

	if err != nil {
		log.BootErrorf("[user_data.clearUserTokens] error occurs when clearing user tokens")
		return err
	}

	log.BootInfof("[user_data.clearUserTokens] all tokens of user \"%s\" has been cleared", username)

	return nil
}

func checkUserTransactionAndAccount(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")

	log.BootInfof("[user_data.checkUserTransactionAndAccount] starting checking user \"%s\" data", username)

	_, err = clis.UserData.CheckTransactionAndAccount(c, username)

	if err != nil {
		log.BootErrorf("[user_data.checkUserTransactionAndAccount] error occurs when checking user data")
		return err
	}

	log.BootInfof("[user_data.checkUserTransactionAndAccount] user transactions and accounts data has been checked successfully, there is no problem with user data")

	return nil
}

func exportUserTransaction(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	username := c.String("username")
	filePath := c.String("file")

	if filePath == "" {
		log.BootErrorf("[user_data.exportUserTransaction] export file path is not specified")
		return os.ErrNotExist
	}

	fileExists, err := utils.IsExists(filePath)

	if fileExists {
		log.BootErrorf("[user_data.exportUserTransaction] specified file path already exists")
		return os.ErrExist
	}

	log.BootInfof("[user_data.exportUserTransaction] starting exporting user \"%s\" data", username)

	content, err := clis.UserData.ExportTransaction(c, username)

	if err != nil {
		log.BootErrorf("[user_data.exportUserTransaction] error occurs when exporting user data")
		return err
	}

	err = utils.WriteFile(filePath, content)

	if err != nil {
		log.BootErrorf("[user_data.exportUserTransaction] failed to write to %s", filePath)
		return err
	}

	log.BootInfof("[user_data.exportUserTransaction] user transactions have been exported to %s", filePath)

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
	fmt.Printf("[LongDateFormat] %s (%d)\n", user.LongDateFormat, user.LongDateFormat)
	fmt.Printf("[ShortDateFormat] %s (%d)\n", user.ShortDateFormat, user.ShortDateFormat)
	fmt.Printf("[LongTimeFormat] %s (%d)\n", user.LongTimeFormat, user.LongTimeFormat)
	fmt.Printf("[ShortTimeFormat] %s (%d)\n", user.ShortTimeFormat, user.ShortTimeFormat)
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
	fmt.Printf("[UserAgent] %s\n", token.UserAgent)
}

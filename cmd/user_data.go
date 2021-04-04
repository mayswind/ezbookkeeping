package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	clis "github.com/mayswind/lab/pkg/cli"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/utils"
)

// UserData represents the data command
var UserData = &cli.Command{
	Name:  "userdata",
	Usage: "lab user data maintenance",
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
			Name:   "user-token-clear",
			Usage:  "Clear user all tokens",
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

	utils.PrintObjectFields(user)

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

	utils.PrintObjectFields(user)

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

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
			Name:   "user-get",
			Usage:  "Get specified user info",
			Action: getUserInfo,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"n"},
					Usage:   "Specific user name",
				},
			},
		},
		{
			Name:   "user-modify-password",
			Usage:  "Modify user password",
			Action: modifyUserPassword,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"n"},
					Usage:   "Specific user name",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "User new password",
				},
			},
		},
		{
			Name:   "user-delete",
			Usage:  "Delete specified user",
			Action: deleteUser,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"n"},
					Usage:   "Specific user name",
				},
			},
		},
		{
			Name:   "transaction-check",
			Usage:  "Check whether user all transactions and accounts are correct",
			Action: checkUserTransactionAndAccount,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"n"},
					Usage:   "Specific user name",
				},
			},
		},
		{
			Name:   "transaction-export",
			Usage:  "Export user all transactions to csv file",
			Action: exportUserTransaction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"n"},
					Usage:   "Specific user name",
				},
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Specific exported file path (e.g. transaction.csv)",
				},
			},
		},
	},
}

func getUserInfo(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	userName := c.String("username")
	user, err := clis.UserData.GetUserByUsername(c, userName)

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

	userName := c.String("username")
	password := c.String("password")
	err = clis.UserData.ModifyUserPassword(c, userName, password)

	if err != nil {
		log.BootErrorf("[user_data.modifyUserPassword] error occurs when modifying user password")
		return err
	}

	log.BootInfof("[user_data.modifyUserPassword] password of user \"%s\" has been changed", userName)

	return nil
}

func deleteUser(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	userName := c.String("username")
	err = clis.UserData.DeleteUser(c, userName)

	if err != nil {
		log.BootErrorf("[user_data.deleteUser] error occurs when deleting user")
		return err
	}

	log.BootInfof("[user_data.deleteUser] user \"%s\" has been deleted", userName)

	return nil
}

func checkUserTransactionAndAccount(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	userName := c.String("username")
	uid, err := clis.UserData.GetUserIdByUsername(c, userName)

	if err != nil {
		log.BootErrorf("[user_data.checkUserTransactionAndAccount] error occurs when getting user id by user name")
		return err
	}

	log.BootInfof("[user_data.checkUserTransactionAndAccount] starting checking user \"%s\" data", userName)

	_, err = clis.UserData.CheckTransactionAndAccount(c, uid)

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

	userName := c.String("username")
	uid, err := clis.UserData.GetUserIdByUsername(c, userName)

	if err != nil {
		log.BootErrorf("[user_data.exportUserTransaction] error occurs when getting user id by user name")
		return err
	}

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

	log.BootInfof("[user_data.exportUserTransaction] starting exporting user \"%s\" data", userName)

	content, err := clis.UserData.ExportTransaction(c, uid)

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

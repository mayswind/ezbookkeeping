package cmd

import (
	"github.com/urfave/cli/v2"

	clis "github.com/mayswind/lab/pkg/cli"
	"github.com/mayswind/lab/pkg/log"
)

// UserData represents the data command
var UserData = &cli.Command{
	Name:  "userdata",
	Usage: "lab user data maintenance",
	Subcommands: []*cli.Command{
		{
			Name:   "check",
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
	},
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

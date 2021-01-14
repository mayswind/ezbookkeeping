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
			Usage:  "Check whether all user transactions and all user accounts are correct",
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
		log.BootErrorf("[data.checkAccountBalance] error occurs when getting user id by user name")
		return err
	}

	log.BootInfof("[data.checkAccountBalance] starting checking user \"%s\" data", userName)

	_, err = clis.UserData.CheckTransactionAndAccount(c, uid)

	if err != nil {
		log.BootErrorf("[data.checkAccountBalance] error occurs when checking user data")
		return err
	}

	log.BootInfof("[data.checkAccountBalance] user transactions and accounts data has been checked successfully, there is no problem with user data")

	return nil
}

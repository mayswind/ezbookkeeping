package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
)

var Database = &cli.Command{
	Name:  "database",
	Usage: "lab database maintenance",
	Subcommands: []*cli.Command{
		{
			Name:   "update",
			Usage:  "Update database structure",
			Action: updateDatabaseStructure,
		},
	},
}

func updateDatabaseStructure(c *cli.Context) error {
	_, err := initializeSystem(c)

	if err != nil {
		return err
	}

	log.BootInfof("[database.updateDatabaseStructure] starting maintaining")

	_ = datastore.Container.UserStore.SyncStructs(new(models.User), new(models.TwoFactor), new(models.TwoFactorRecoveryCode))
	_ = datastore.Container.TokenStore.SyncStructs(new(models.TokenRecord))
	_ = datastore.Container.UserDataStore.SyncStructs(new(models.Account))

	log.BootInfof("[database.updateDatabaseStructure] maintained successfully")

	return nil
}

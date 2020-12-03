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

	err = updateAllDatabaseTablesStructure()

	if err != nil {
		log.BootErrorf("[database.updateDatabaseStructure] update database table structure failed, because %s", err.Error())
		return err
	}

	log.BootInfof("[database.updateDatabaseStructure] all tables maintained successfully")
	return nil
}

func updateAllDatabaseTablesStructure() error {
	var err error

	err = datastore.Container.UserStore.SyncStructs(new(models.User))

	if err != nil {
		return err
	} else {
		log.BootInfof("[database.updateAllDatabaseTablesStructure] user table maintained successfully")
	}

	err = datastore.Container.UserStore.SyncStructs(new(models.TwoFactor))

	if err != nil {
		return err
	} else {
		log.BootInfof("[database.updateAllDatabaseTablesStructure] two factor table maintained successfully")
	}

	err = datastore.Container.UserStore.SyncStructs(new(models.TwoFactorRecoveryCode))

	if err != nil {
		return err
	} else {
		log.BootInfof("[database.updateAllDatabaseTablesStructure] two factor recovery code table maintained successfully")
	}

	err = datastore.Container.TokenStore.SyncStructs(new(models.TokenRecord))

	if err != nil {
		return err
	} else {
		log.BootInfof("[database.updateAllDatabaseTablesStructure] token record table maintained successfully")
	}

	err = datastore.Container.UserDataStore.SyncStructs(new(models.Account))

	if err != nil {
		return err
	} else {
		log.BootInfof("[database.updateAllDatabaseTablesStructure] account table maintained successfully")
	}

	err = datastore.Container.UserDataStore.SyncStructs(new(models.TransactionCategory))

	if err != nil {
		return err
	} else {
		log.BootInfof("[database.updateAllDatabaseTablesStructure] transaction category table maintained successfully")
	}

	return nil
}

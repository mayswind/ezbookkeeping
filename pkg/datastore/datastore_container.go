package datastore

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// DataStoreContainer contains all data storages
type DataStoreContainer struct {
	UserStore     *DataStore
	TokenStore    *DataStore
	UserDataStore *DataStore
}

// Initialize a data storage container singleton instance
var (
	Container = &DataStoreContainer{}
)

// InitializeDataStore initializes data storage according to the config
func InitializeDataStore(config *settings.Config) error {
	database, err := initializeDatabase(config.DatabaseConfig)

	if err != nil {
		return err
	}

	setDatabaseLogger(database, config)

	Container.UserStore, err = NewDataStore(database)

	if err != nil {
		return err
	}

	Container.TokenStore, err = NewDataStore(database)

	if err != nil {
		return err
	}

	Container.UserDataStore, err = NewDataStore(database)

	if err != nil {
		return err
	}

	return nil
}

func initializeDatabase(dbConfig *settings.DatabaseConfig) (*Database, error) {
	var connStr string
	var err error

	if dbConfig.DatabaseType == settings.Sqlite3DbType {
		if _, err = os.Stat(dbConfig.DatabasePath); err != nil {
			file, err := os.Create(dbConfig.DatabasePath)

			if err != nil {
				return nil, err
			}

			defer file.Close()
		}
	}

	if dbConfig.DatabaseType == settings.MySqlDbType {
		connStr, err = getMysqlConnectionString(dbConfig)
	} else if dbConfig.DatabaseType == settings.PostgresDbType {
		connStr, err = getPostgresConnectionString(dbConfig)
	} else if dbConfig.DatabaseType == settings.Sqlite3DbType {
		connStr, err = getSqlite3ConnectionString(dbConfig)
	} else {
		return nil, errs.ErrDatabaseTypeInvalid
	}

	if err != nil {
		return nil, err
	}

	connStrs := []string{
		connStr,
	}
	engineGroup, err := xorm.NewEngineGroup(dbConfig.DatabaseType, connStrs, xorm.RoundRobinPolicy())

	if err != nil {
		return nil, err
	}

	engineGroup.SetMaxIdleConns(int(dbConfig.MaxIdleConnection))
	engineGroup.SetMaxOpenConns(int(dbConfig.MaxOpenConnection))
	engineGroup.SetConnMaxLifetime(time.Duration(dbConfig.ConnectionMaxLifeTime) * time.Second)

	return &Database{
		databaseType: dbConfig.DatabaseType,
		engineGroup:  engineGroup,
	}, nil
}

func setDatabaseLogger(database *Database, config *settings.Config) {
	if config.EnableQueryLog {
		database.engineGroup.SetLogger(NewXOrmLoggerAdapter(config.EnableQueryLog, config.LogLevel))
		database.engineGroup.ShowSQL(true)
	}
}

func getMysqlConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	protocol := "tcp"

	if strings.HasPrefix(dbConfig.DatabaseHost, "/") { // unix socket path
		protocol = "unix"
	}

	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=true",
		dbConfig.DatabaseUser, dbConfig.DatabasePassword, protocol, dbConfig.DatabaseHost, dbConfig.DatabaseName), nil
}

func getPostgresConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	host, port, err := net.SplitHostPort(dbConfig.DatabaseHost)

	if err != nil {
		return "", errs.ErrDatabaseHostInvalid
	}

	if strings.HasPrefix(dbConfig.DatabaseHost, "/") { // unix socket path
		return fmt.Sprintf("postgres://%s:%s@:%s/%s?sslmode=%s&host=%s",
			url.QueryEscape(dbConfig.DatabaseUser), url.QueryEscape(dbConfig.DatabasePassword), port, dbConfig.DatabaseName, dbConfig.DatabaseSSLMode, host), nil
	} else {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			url.QueryEscape(dbConfig.DatabaseUser), url.QueryEscape(dbConfig.DatabasePassword), host, port, dbConfig.DatabaseName, dbConfig.DatabaseSSLMode), nil
	}
}

func getSqlite3ConnectionString(dbConfig *settings.DatabaseConfig) (string, error) {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbConfig.DatabasePath), nil
}

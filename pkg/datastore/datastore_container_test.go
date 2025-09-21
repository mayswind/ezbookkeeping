package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestGetMysqlConnectionString_TCP(t *testing.T) {
	expectedValue := "username:password@tcp(1.2.3.4:3306)/dbname?charset=utf8mb4&parseTime=true"
	actualValue, err := getMysqlConnectionString(&settings.DatabaseConfig{
		DatabaseType:     "mysql",
		DatabaseHost:     "1.2.3.4:3306",
		DatabaseName:     "dbname",
		DatabaseUser:     "username",
		DatabasePassword: "password",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetMysqlConnectionString_UnixSocket(t *testing.T) {
	expectedValue := "username:password@unix(/path/to/mysql.sock)/dbname?charset=utf8mb4&parseTime=true"
	actualValue, err := getMysqlConnectionString(&settings.DatabaseConfig{
		DatabaseType:     "mysql",
		DatabaseHost:     "/path/to/mysql.sock",
		DatabaseName:     "dbname",
		DatabaseUser:     "username",
		DatabasePassword: "password",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetPostgreSQLConnectionString_TCP(t *testing.T) {
	expectedValue := "postgres://username:password@1.2.3.4:5432/dbname?sslmode=disable"
	actualValue, err := getPostgresConnectionString(&settings.DatabaseConfig{
		DatabaseType:     "postgres",
		DatabaseHost:     "1.2.3.4:5432",
		DatabaseName:     "dbname",
		DatabaseUser:     "username",
		DatabasePassword: "password",
		DatabaseSSLMode:  "disable",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestGetPostgreSQLConnectionString_UnixSocket(t *testing.T) {
	expectedValue := "postgres:///dbname?sslmode=disable&host=/path/to/postgres.sock&user=username&password=password"
	actualValue, err := getPostgresConnectionString(&settings.DatabaseConfig{
		DatabaseType:     "postgres",
		DatabaseHost:     "/path/to/postgres.sock",
		DatabaseName:     "dbname",
		DatabaseUser:     "username",
		DatabasePassword: "password",
		DatabaseSSLMode:  "disable",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

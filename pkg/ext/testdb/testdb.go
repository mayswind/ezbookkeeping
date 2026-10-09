// Package testdb gives the ext tests a throw-away database. SQLite is the default; set EXT_TEST_DB to "postgres" or
// "mysql" (with EXT_TEST_DB_HOST, EXT_TEST_DB_USER and EXT_TEST_DB_PASSWD) to run the very same tests against a
// real server. Every test gets its own empty database, which is dropped afterwards.
package testdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// Kind returns the database type the tests run on
func Kind() string {
	if kind := os.Getenv("EXT_TEST_DB"); kind != "" {
		return kind
	}

	return settings.Sqlite3DbType
}

// Config returns a configuration pointing at a fresh empty database, plus the uuid settings the services need
func Config(t *testing.T) *settings.Config {
	t.Helper()

	config := &settings.Config{
		UuidGeneratorType: settings.InternalUuidGeneratorType,
		UuidServerId:      1,
	}

	host := os.Getenv("EXT_TEST_DB_HOST")
	user := os.Getenv("EXT_TEST_DB_USER")
	passwd := os.Getenv("EXT_TEST_DB_PASSWD")
	name := fmt.Sprintf("ext_test_%d_%d", time.Now().UnixNano()%1_000_000_000, rand.Intn(1_000_000))

	switch Kind() {
	case settings.PostgresDbType:
		admin, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/postgres?sslmode=disable", user, passwd, host))
		must(t, err)
		_, err = admin.Exec("CREATE DATABASE " + name)
		must(t, err)
		t.Cleanup(func() {
			_, _ = admin.Exec("DROP DATABASE IF EXISTS " + name + " WITH (FORCE)")
			_ = admin.Close()
		})

		config.DatabaseConfig = &settings.DatabaseConfig{DatabaseType: settings.PostgresDbType, DatabaseHost: host, DatabaseName: name, DatabaseUser: user, DatabasePassword: passwd, DatabaseSSLMode: "disable", MaxOpenConnection: 20}
	case settings.MySqlDbType:
		admin, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", user, passwd, host))
		must(t, err)
		_, err = admin.Exec("CREATE DATABASE " + name + " CHARACTER SET utf8mb4")
		must(t, err)
		t.Cleanup(func() {
			_, _ = admin.Exec("DROP DATABASE IF EXISTS " + name)
			_ = admin.Close()
		})

		config.DatabaseConfig = &settings.DatabaseConfig{DatabaseType: settings.MySqlDbType, DatabaseHost: host, DatabaseName: name, DatabaseUser: user, DatabasePassword: passwd, MaxOpenConnection: 20}
	default:
		config.DatabaseConfig = &settings.DatabaseConfig{DatabaseType: settings.Sqlite3DbType, DatabasePath: filepath.Join(t.TempDir(), "test.db")}
	}

	return config
}

func must(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("test database: %v", err)
	}
}

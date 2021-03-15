package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"github.com/mayswind/lab/pkg/errs"
)

const (
	labWorkDirEnvName     = "LAB_WORK_DIR"
	labEnvNamePrefix      = "LAB"
	defaultConfigPath     = "/conf/lab.ini"
	defaultStaticRootPath = "public"
)

// SystemMode represents running mode of system
type SystemMode string

// System running modes
const (
	MODE_DEVELOPMENT SystemMode = "development"
	MODE_PRODUCTION  SystemMode = "production"
)

// Scheme represents how the web backend service exposes
type Scheme string

// Scheme types
const (
	SCHEME_HTTP   Scheme = "http"
	SCHEME_HTTPS  Scheme = "https"
	SCHEME_SOCKET Scheme = "socket"
)

// Level represents log level
type Level string

// Log levels
const (
	LOGLEVEL_DEBUG Level = "debug"
	LOGLEVEL_INFO  Level = "info"
	LOGLEVEL_WARN  Level = "warn"
	LOGLEVEL_ERROR Level = "error"
)

// Database types
const (
	MySqlDbType    string = "mysql"
	PostgresDbType string = "postgres"
	Sqlite3DbType  string = "sqlite3"
)

// Uuid generator types
const (
	InternalUuidGeneratorType string = "internal"
)

// Exchange rates data source types
const (
	EuroCentralBankDataSource      string = "euro_central_bank"
	BankOfCanadaDataSource         string = "bank_of_canada"
	CzechNationalBankDataSource    string = "czech_national_bank"
	NationalBankOfPolandDataSource string = "national_bank_of_poland"
)

const (
	defaultAppName string = "lab"

	defaultHttpAddr string = "0.0.0.0"
	defaultHttpPort int    = 8080
	defaultDomain   string = "localhost"

	defaultDatabaseHost            string = "127.0.0.1:3306"
	defaultDatabaseName            string = "lab"
	defaultDatabaseMaxIdleConn     int    = 2
	defaultDatabaseMaxOpenConn     int    = 0
	defaultDatabaseConnMaxLifetime int    = 14400

	defaultLogMode  string = "console"
	defaultLoglevel Level  = LOGLEVEL_INFO

	defaultSecretKey                 string = "lab"
	defaultTokenExpiredTime          int    = 604800 // 7 days
	defaultTemporaryTokenExpiredTime int    = 300    // 5 minutes

	defaultExchangeRatesDataRequestTimeout int = 10000 // 10 seconds
)

// DatabaseConfig represents the database setting config
type DatabaseConfig struct {
	DatabaseType     string
	DatabaseHost     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string

	DatabaseSSLMode string

	DatabasePath string

	MaxIdleConnection     int
	MaxOpenConnection     int
	ConnectionMaxLifeTime int
}

// Config represents the global setting config
type Config struct {
	// Global
	AppName     string
	Mode        SystemMode
	WorkingPath string

	// Server
	Protocol Scheme
	HttpAddr string
	HttpPort int

	Domain  string
	RootUrl string

	CertFile    string
	CertKeyFile string

	UnixSocketPath string

	StaticRootPath string

	EnableGZip       bool
	EnableRequestLog bool

	// Database
	DatabaseConfig     *DatabaseConfig
	EnableQueryLog     bool
	AutoUpdateDatabase bool

	// Log
	LogModes         []string
	EnableConsoleLog bool
	EnableFileLog    bool

	LogLevel    Level
	FileLogPath string

	// Uuid
	UuidGeneratorType string
	UuidServerId      uint8

	// Secret
	SecretKey                         string
	EnableTwoFactor                   bool
	TokenExpiredTime                  int
	TokenExpiredTimeDuration          time.Duration
	TemporaryTokenExpiredTime         int
	TemporaryTokenExpiredTimeDuration time.Duration
	EnableRequestIdHeader             bool

	// User
	EnableUserRegister bool

	// Data
	EnableDataExport bool

	// Exchange Rates
	ExchangeRatesDataSource     string
	ExchangeRatesRequestTimeout int
}

// LoadConfiguration loads setting config from given config file path
func LoadConfiguration(configFilePath string) (*Config, error) {
	var err error

	cfgFile, err := ini.LoadSources(ini.LoadOptions{}, configFilePath)

	if err != nil {
		return nil, err
	}

	config := &Config{}
	config.WorkingPath, err = getWorkingPath()

	if err != nil {
		return nil, err
	}

	err = loadGlobalConfiguration(config, cfgFile, "global")

	if err != nil {
		return nil, err
	}

	err = loadServerConfiguration(config, cfgFile, "server")

	if err != nil {
		return nil, err
	}

	err = loadDatabaseConfiguration(config, cfgFile, "database")

	if err != nil {
		return nil, err
	}

	err = loadLogConfiguration(config, cfgFile, "log")

	if err != nil {
		return nil, err
	}

	err = loadUuidConfiguration(config, cfgFile, "uuid")

	if err != nil {
		return nil, err
	}

	err = loadSecurityConfiguration(config, cfgFile, "security")

	if err != nil {
		return nil, err
	}

	err = loadUserConfiguration(config, cfgFile, "user")

	if err != nil {
		return nil, err
	}

	err = loadDataConfiguration(config, cfgFile, "data")

	if err != nil {
		return nil, err
	}

	err = loadExchangeRatesConfiguration(config, cfgFile, "exchange_rates")

	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetDefaultConfigFilePath returns the defaule config file path
func GetDefaultConfigFilePath() (string, error) {
	workingPath, err := getWorkingPath()

	if err != nil {
		return "", err
	}

	cfgFilePath := filepath.Join(workingPath, defaultConfigPath)
	_, err = os.Stat(cfgFilePath)

	if err != nil {
		return "", err
	}

	return cfgFilePath, nil
}

func loadGlobalConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.AppName = getConfigItemStringValue(configFile, sectionName, "app_name", defaultAppName)
	config.Mode = MODE_PRODUCTION

	if getConfigItemStringValue(configFile, sectionName, "mode") == "development" {
		config.Mode = MODE_DEVELOPMENT
	}

	return nil
}

func loadServerConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "protocol") == "http" {
		config.Protocol = SCHEME_HTTP

		config.HttpAddr = getConfigItemStringValue(configFile, sectionName, "http_addr", defaultHttpAddr)
		config.HttpPort = getConfigItemIntValue(configFile, sectionName, "http_port", defaultHttpPort)
	} else if getConfigItemStringValue(configFile, sectionName, "protocol") == "https" {
		config.Protocol = SCHEME_HTTPS

		config.HttpAddr = getConfigItemStringValue(configFile, sectionName, "http_addr", defaultHttpAddr)
		config.HttpPort = getConfigItemIntValue(configFile, sectionName, "http_port", defaultHttpPort)

		config.CertFile = getConfigItemStringValue(configFile, sectionName, "cert_file")
		config.CertKeyFile = getConfigItemStringValue(configFile, sectionName, "cert_key_file")
	} else if getConfigItemStringValue(configFile, sectionName, "protocol") == "socket" {
		config.Protocol = SCHEME_SOCKET

		config.UnixSocketPath = getConfigItemStringValue(configFile, sectionName, "unix_socket")
	} else {
		return errs.ErrInvalidProtocol
	}

	config.Domain = getConfigItemStringValue(configFile, sectionName, "domain", defaultDomain)
	config.RootUrl = getConfigItemStringValue(configFile, sectionName, "root_url", fmt.Sprintf("%s://%s:%d/", string(config.Protocol), config.Domain, config.HttpPort))

	if config.RootUrl[len(config.RootUrl)-1] != '/' {
		config.RootUrl += "/"
	}

	staticRootPath := getConfigItemStringValue(configFile, sectionName, "static_root_path", defaultStaticRootPath)
	finalStaticRootPath, err := getFinalPath(config.WorkingPath, staticRootPath)

	if err != nil {
		return err
	}

	config.StaticRootPath = finalStaticRootPath

	config.EnableGZip = getConfigItemBoolValue(configFile, sectionName, "enable_gzip", false)
	config.EnableRequestLog = getConfigItemBoolValue(configFile, sectionName, "log_request", false)

	return nil
}

func loadDatabaseConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	dbConfig := &DatabaseConfig{}

	dbConfig.DatabaseType = getConfigItemStringValue(configFile, sectionName, "type", MySqlDbType)
	dbConfig.DatabaseHost = getConfigItemStringValue(configFile, sectionName, "host", defaultDatabaseHost)
	dbConfig.DatabaseName = getConfigItemStringValue(configFile, sectionName, "name", defaultDatabaseName)
	dbConfig.DatabaseUser = getConfigItemStringValue(configFile, sectionName, "user")
	dbConfig.DatabasePassword = getConfigItemStringValue(configFile, sectionName, "passwd")

	if dbConfig.DatabaseType == PostgresDbType {
		dbConfig.DatabaseSSLMode = getConfigItemStringValue(configFile, sectionName, "ssl_mode")
	}

	if dbConfig.DatabaseType == Sqlite3DbType {
		staticDBPath := getConfigItemStringValue(configFile, sectionName, "db_path")
		finalStaticDBPath, _ := getFinalPath(config.WorkingPath, staticDBPath)
		dbConfig.DatabasePath = finalStaticDBPath
	}

	dbConfig.MaxIdleConnection = getConfigItemIntValue(configFile, sectionName, "max_idle_conn", defaultDatabaseMaxIdleConn)
	dbConfig.MaxOpenConnection = getConfigItemIntValue(configFile, sectionName, "max_open_conn", defaultDatabaseMaxOpenConn)
	dbConfig.ConnectionMaxLifeTime = getConfigItemIntValue(configFile, sectionName, "conn_max_lifetime", defaultDatabaseConnMaxLifetime)

	config.DatabaseConfig = dbConfig
	config.EnableQueryLog = getConfigItemBoolValue(configFile, sectionName, "log_query", false)
	config.AutoUpdateDatabase = getConfigItemBoolValue(configFile, sectionName, "auto_update_database", true)

	return nil
}

func loadLogConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.LogModes = strings.Split(getConfigItemStringValue(configFile, sectionName, "mode", defaultLogMode), " ")

	for i := 0; i < len(config.LogModes); i++ {
		logMode := config.LogModes[i]

		if logMode == "console" {
			config.EnableConsoleLog = true
		} else if logMode == "file" {
			config.EnableFileLog = true
		} else {
			return errs.ErrInvalidLogMode
		}
	}

	config.LogLevel = getLogLevel(getConfigItemStringValue(configFile, sectionName, "level"), defaultLoglevel)

	if config.EnableFileLog {
		config.FileLogPath = getConfigItemStringValue(configFile, sectionName, "log_path")
	}

	return nil
}

func loadUuidConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "generator_type") == InternalUuidGeneratorType {
		config.UuidGeneratorType = InternalUuidGeneratorType
	} else {
		return errs.ErrInvalidUuidMode
	}

	config.UuidServerId = uint8(getConfigItemIntValue(configFile, sectionName, "server_id", 0))

	return nil
}

func loadSecurityConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.SecretKey = getConfigItemStringValue(configFile, sectionName, "secret_key", defaultSecretKey)
	config.EnableTwoFactor = getConfigItemBoolValue(configFile, sectionName, "enable_two_factor", true)

	config.TokenExpiredTime = getConfigItemIntValue(configFile, sectionName, "token_expired_time", defaultTokenExpiredTime)
	config.TokenExpiredTimeDuration = time.Duration(config.TokenExpiredTime) * time.Second

	config.TemporaryTokenExpiredTime = getConfigItemIntValue(configFile, sectionName, "temporary_token_expired_time", defaultTemporaryTokenExpiredTime)
	config.TemporaryTokenExpiredTimeDuration = time.Duration(config.TemporaryTokenExpiredTime) * time.Second

	config.EnableRequestIdHeader = getConfigItemBoolValue(configFile, sectionName, "request_id_header", true)

	return nil
}

func loadUserConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableUserRegister = getConfigItemBoolValue(configFile, sectionName, "enable_register", false)

	return nil
}

func loadDataConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableDataExport = getConfigItemBoolValue(configFile, sectionName, "enable_export", false)

	return nil
}

func loadExchangeRatesConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "data_source") == EuroCentralBankDataSource {
		config.ExchangeRatesDataSource = EuroCentralBankDataSource
	} else if getConfigItemStringValue(configFile, sectionName, "data_source") == BankOfCanadaDataSource {
		config.ExchangeRatesDataSource = BankOfCanadaDataSource
	} else if getConfigItemStringValue(configFile, sectionName, "data_source") == CzechNationalBankDataSource {
		config.ExchangeRatesDataSource = CzechNationalBankDataSource
	} else if getConfigItemStringValue(configFile, sectionName, "data_source") == NationalBankOfPolandDataSource {
		config.ExchangeRatesDataSource = NationalBankOfPolandDataSource
	} else {
		return errs.ErrInvalidExchangeRatesDataSource
	}

	config.ExchangeRatesRequestTimeout = getConfigItemIntValue(configFile, sectionName, "request_timeout", defaultExchangeRatesDataRequestTimeout)

	return nil
}

func getWorkingPath() (string, error) {
	workingPath := os.Getenv(labWorkDirEnvName)

	if workingPath != "" {
		return workingPath, nil
	}

	execFilePath, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return execFilePath, nil
}

func getFinalPath(workingPath, p string) (string, error) {
	var err error

	if !filepath.IsAbs(p) {
		p = filepath.Join(workingPath, p)
	}

	if _, err = os.Stat(p); err == nil {
		return p, nil
	}

	return p, err
}

func getConfigItemStringValue(configFile *ini.File, sectionName string, itemName string, defaultValue ...string) string {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		return environmentValue
	}

	section := configFile.Section(sectionName)

	if len(defaultValue) > 0 {
		return section.Key(itemName).MustString(defaultValue[0])
	} else {
		return section.Key(itemName).String()
	}
}

func getConfigItemIntValue(configFile *ini.File, sectionName string, itemName string, defaultValue ...int) int {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		value, err := strconv.ParseInt(environmentValue, 0, 64)

		if err == nil {
			return int(value)
		}
	}

	section := configFile.Section(sectionName)
	return section.Key(itemName).MustInt(defaultValue...)
}

func getConfigItemBoolValue(configFile *ini.File, sectionName string, itemName string, defaultValue ...bool) bool {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		value, err := strconv.ParseBool(environmentValue)

		if err == nil {
			return value
		}
	}

	section := configFile.Section(sectionName)
	return section.Key(itemName).MustBool(defaultValue...)
}

func getEnvironmentKey(sectionName string, itemName string) string {
	return fmt.Sprintf("%s_%s_%s", labEnvNamePrefix, strings.ToUpper(sectionName), strings.ToUpper(itemName))
}

func getLogLevel(logLevelStr string, defaultLogLevel Level) Level {
	if logLevelStr == "debug" {
		return LOGLEVEL_DEBUG
	} else if logLevelStr == "warn" {
		return LOGLEVEL_WARN
	} else if logLevelStr == "error" {
		return LOGLEVEL_ERROR
	}

	return defaultLogLevel
}

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
	LAB_WORK_DIR               = "LAB_WORK_DIR"
	LAB_ENVIRONMENT_KEY_PREFIX = "LAB"
	DEFAULT_CONFIG_PATH        = "/conf/lab.ini"
	DEFAULT_STATIC_ROOT_PATH   = "public"
)

type SystemMode string

const (
	MODE_DEVELOPMENT SystemMode = "development"
	MODE_PRODUCTION  SystemMode = "production"
)

type Scheme string

const (
	SCHEME_HTTP   Scheme = "http"
	SCHEME_HTTPS  Scheme = "https"
	SCHEME_SOCKET Scheme = "socket"
)

type Level string

const (
	LOGLEVEL_DEBUG Level = "debug"
	LOGLEVEL_INFO  Level = "info"
	LOGLEVEL_WARN  Level = "warn"
	LOGLEVEL_ERROR Level = "error"
)

const (
	DBTYPE_MYSQL    string = "mysql"
	DBTYPE_POSTGRES string = "postgres"
	DBTYPE_SQLITE3  string = "sqlite3"
)

const (
	UUID_GENERATOR_TYPE_INTERNAL string = "internal"
)

const (
	DEFAULT_APP_NAME string = "lab"

	DEFAULT_HTTP_ADDR string = "0.0.0.0"
	DEFAULT_HTTP_PORT int    = 8080
	DEFAULT_DOMAIN    string = "localhost"

	DEFAULT_DATABASE_HOST              string = "127.0.0.1:3306"
	DEFAULT_DATABASE_NAME              string = "lab"
	DEFAULT_DATABASE_MAX_IDLE_CONN     int    = 2
	DEFAULT_DATABASE_MAX_OPEN_CONN     int    = 0
	DEFAULT_DATABASE_CONN_MAX_LIFETIME int    = 14400

	DEFAULT_LOG_MODE  string = "console"
	DEFAULT_LOG_LEVEL Level  = LOGLEVEL_INFO

	DEFAULT_SECRET_KEY                   string = "lab"
	DEFAULT_TOKEN_EXPIRED_TIME           int    = 604800 // 7 days
	DEFAULT_TEMPORARY_TOKEN_EXPIRED_TIME int    = 300    // 5 minutes
)

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
	DatabaseConfig *DatabaseConfig
	EnableQueryLog bool

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
}

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

	return config, nil
}

func GetDefaultConfigFilePath() (string, error) {
	workingPath, err := getWorkingPath()

	if err != nil {
		return "", err
	}

	cfgFilePath := filepath.Join(workingPath, DEFAULT_CONFIG_PATH)
	_, err = os.Stat(cfgFilePath)

	if err != nil {
		return "", err
	}

	return cfgFilePath, nil
}

func loadGlobalConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.AppName = getConfigItemStringValue(configFile, sectionName, "app_name", DEFAULT_APP_NAME)
	config.Mode = MODE_PRODUCTION

	if getConfigItemStringValue(configFile, sectionName, "mode") == "development" {
		config.Mode = MODE_DEVELOPMENT
	}

	return nil
}

func loadServerConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "protocol") == "http" {
		config.Protocol = SCHEME_HTTP

		config.HttpAddr = getConfigItemStringValue(configFile, sectionName, "http_addr", DEFAULT_HTTP_ADDR)
		config.HttpPort = getConfigItemIntValue(configFile, sectionName, "http_port", DEFAULT_HTTP_PORT)
	} else if getConfigItemStringValue(configFile, sectionName, "protocol") == "https" {
		config.Protocol = SCHEME_HTTPS

		config.HttpAddr = getConfigItemStringValue(configFile, sectionName, "http_addr", DEFAULT_HTTP_ADDR)
		config.HttpPort = getConfigItemIntValue(configFile, sectionName, "http_port", DEFAULT_HTTP_PORT)

		config.CertFile = getConfigItemStringValue(configFile, sectionName, "cert_file")
		config.CertKeyFile = getConfigItemStringValue(configFile, sectionName, "cert_key_file")
	} else if getConfigItemStringValue(configFile, sectionName, "protocol") == "socket" {
		config.Protocol = SCHEME_SOCKET

		config.UnixSocketPath = getConfigItemStringValue(configFile, sectionName, "unix_socket")
	} else {
		return errs.ErrInvalidProtocol
	}

	config.Domain = getConfigItemStringValue(configFile, sectionName, "domain", DEFAULT_DOMAIN)
	config.RootUrl = getConfigItemStringValue(configFile, sectionName, "root_url", fmt.Sprintf("%s://%s:%d/", string(config.Protocol), config.Domain, config.HttpPort))

	if config.RootUrl[len(config.RootUrl)-1] != '/' {
		config.RootUrl += "/"
	}

	staticRootPath := getConfigItemStringValue(configFile, sectionName, "static_root_path", DEFAULT_STATIC_ROOT_PATH)
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

	dbConfig.DatabaseType = getConfigItemStringValue(configFile, sectionName, "type", DBTYPE_MYSQL)
	dbConfig.DatabaseHost = getConfigItemStringValue(configFile, sectionName, "host", DEFAULT_DATABASE_HOST)
	dbConfig.DatabaseName = getConfigItemStringValue(configFile, sectionName, "name", DEFAULT_DATABASE_NAME)
	dbConfig.DatabaseUser = getConfigItemStringValue(configFile, sectionName, "user")
	dbConfig.DatabasePassword = getConfigItemStringValue(configFile, sectionName, "passwd")

	if dbConfig.DatabaseType == DBTYPE_POSTGRES {
		dbConfig.DatabaseSSLMode = getConfigItemStringValue(configFile, sectionName, "ssl_mode")
	}

	if dbConfig.DatabaseType == DBTYPE_SQLITE3 {
		dbConfig.DatabasePath = getConfigItemStringValue(configFile, sectionName, "db_path")
	}

	dbConfig.MaxIdleConnection = getConfigItemIntValue(configFile, sectionName, "max_idle_conn", DEFAULT_DATABASE_MAX_IDLE_CONN)
	dbConfig.MaxOpenConnection = getConfigItemIntValue(configFile, sectionName, "max_open_conn", DEFAULT_DATABASE_MAX_OPEN_CONN)
	dbConfig.ConnectionMaxLifeTime = getConfigItemIntValue(configFile, sectionName, "conn_max_lifetime", DEFAULT_DATABASE_CONN_MAX_LIFETIME)

	config.DatabaseConfig = dbConfig
	config.EnableQueryLog = getConfigItemBoolValue(configFile, sectionName, "log_query", false)

	return nil
}

func loadLogConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.LogModes = strings.Split(getConfigItemStringValue(configFile, sectionName, "mode", DEFAULT_LOG_MODE), " ")

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

	config.LogLevel = getLogLevel(getConfigItemStringValue(configFile, sectionName, "level"), DEFAULT_LOG_LEVEL)

	if config.EnableFileLog {
		config.FileLogPath = getConfigItemStringValue(configFile, sectionName, "log_path")
	}

	return nil
}

func loadUuidConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "generator_type") == UUID_GENERATOR_TYPE_INTERNAL {
		config.UuidGeneratorType = UUID_GENERATOR_TYPE_INTERNAL
	} else {
		return errs.ErrInvalidUuidMode
	}

	config.UuidServerId = uint8(getConfigItemIntValue(configFile, sectionName, "server_id", 0))

	return nil
}

func loadSecurityConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.SecretKey = getConfigItemStringValue(configFile, sectionName, "secret_key", DEFAULT_SECRET_KEY)
	config.EnableTwoFactor = getConfigItemBoolValue(configFile, sectionName, "enable_two_factor", true)

	config.TokenExpiredTime = getConfigItemIntValue(configFile, sectionName, "token_expired_time", DEFAULT_TOKEN_EXPIRED_TIME)
	config.TokenExpiredTimeDuration = time.Duration(config.TokenExpiredTime) * time.Second

	config.TemporaryTokenExpiredTime = getConfigItemIntValue(configFile, sectionName, "temporary_token_expired_time", DEFAULT_TEMPORARY_TOKEN_EXPIRED_TIME)
	config.TemporaryTokenExpiredTimeDuration = time.Duration(config.TemporaryTokenExpiredTime) * time.Second

	config.EnableRequestIdHeader = getConfigItemBoolValue(configFile, sectionName, "request_id_header", true)

	return nil
}

func loadUserConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableUserRegister = getConfigItemBoolValue(configFile, sectionName, "enable_register", false)

	return nil
}

func getWorkingPath() (string, error) {
	workingPath := os.Getenv(LAB_WORK_DIR)

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

	return "", err
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
	return fmt.Sprintf("%s_%s_%s", LAB_ENVIRONMENT_KEY_PREFIX, strings.ToUpper(sectionName), strings.ToUpper(itemName))
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

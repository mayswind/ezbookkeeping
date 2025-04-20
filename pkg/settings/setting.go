package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/locales"
)

const (
	ebkWorkDirEnvName     = "EBK_WORK_DIR"
	ebkEnvNamePrefix      = "EBK"
	defaultConfigPath     = "/conf/ezbookkeeping.ini"
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

// Object Storage types
const (
	LocalFileSystemObjectStorageType string = "local_filesystem"
	MinIOStorageType                 string = "minio"
)

// Uuid generator types
const (
	InternalUuidGeneratorType string = "internal"
)

// Duplicate checker types
const (
	InMemoryDuplicateCheckerType string = "in_memory"
)

// Map provider types
const (
	OpenStreetMapProvider                  string = "openstreetmap"
	OpenStreetMapHumanitarianStyleProvider string = "openstreetmap_humanitarian"
	OpenTopoMapProvider                    string = "opentopomap"
	OPNVKarteMapProvider                   string = "opnvkarte"
	CyclOSMMapProvider                     string = "cyclosm"
	CartoDBMapProvider                     string = "cartodb"
	TomTomMapProvider                      string = "tomtom"
	TianDiTuProvider                       string = "tianditu"
	GoogleMapProvider                      string = "googlemap"
	BaiduMapProvider                       string = "baidumap"
	AmapProvider                           string = "amap"
	CustomProvider                         string = "custom"
)

// Amap security verification method
const (
	AmapSecurityVerificationInternalProxyMethod string = "internal_proxy"
	AmapSecurityVerificationExternalProxyMethod string = "external_proxy"
	AmapSecurityVerificationPlainTextMethod     string = "plain_text"
)

// Exchange rates data source types
const (
	ReserveBankOfAustraliaDataSource    string = "reserve_bank_of_australia"
	BankOfCanadaDataSource              string = "bank_of_canada"
	CzechNationalBankDataSource         string = "czech_national_bank"
	DanmarksNationalbankDataSource      string = "danmarks_national_bank"
	EuroCentralBankDataSource           string = "euro_central_bank"
	NationalBankOfGeorgiaDataSource     string = "national_bank_of_georgia"
	CentralBankOfHungaryDataSource      string = "central_bank_of_hungary"
	BankOfIsraelDataSource              string = "bank_of_israel"
	CentralBankOfMyanmarDataSource      string = "central_bank_of_myanmar"
	NorgesBankDataSource                string = "norges_bank"
	NationalBankOfPolandDataSource      string = "national_bank_of_poland"
	NationalBankOfRomaniaDataSource     string = "national_bank_of_romania"
	BankOfRussiaDataSource              string = "bank_of_russia"
	SwissNationalBankDataSource         string = "swiss_national_bank"
	NationalBankOfUkraineDataSource     string = "national_bank_of_ukraine"
	CentralBankOfUzbekistanDataSource   string = "central_bank_of_uzbekistan"
	InternationalMonetaryFundDataSource string = "international_monetary_fund"
)

const (
	defaultAppName string = "ezBookkeeping"

	defaultHttpAddr string = "0.0.0.0"
	defaultHttpPort uint16 = 8080
	defaultDomain   string = "localhost"

	defaultDatabaseHost            string = "127.0.0.1:3306"
	defaultDatabaseName            string = "ezbookkeeping"
	defaultDatabaseMaxIdleConn     uint16 = 2
	defaultDatabaseMaxOpenConn     uint16 = 0
	defaultDatabaseConnMaxLifetime uint32 = 14400

	defaultLogMode        string = "console"
	defaultLogFileMaxSize uint32 = 104857600 // 100 MB
	defaultLogFileMaxDays uint32 = 7         // days

	defaultInMemoryDuplicateCheckerCleanupInterval uint32 = 60  // 1 minutes
	defaultDuplicateSubmissionsInterval            uint32 = 300 // 5 minutes

	defaultSecretKey                     string = "ezbookkeeping"
	defaultTokenExpiredTime              uint32 = 2592000 // 30 days
	defaultTokenMinRefreshInterval       uint32 = 86400   // 1 day
	defaultTemporaryTokenExpiredTime     uint32 = 300     // 5 minutes
	defaultEmailVerifyTokenExpiredTime   uint32 = 3600    // 60 minutes
	defaultPasswordResetTokenExpiredTime uint32 = 3600    // 60 minutes
	defaultMaxFailuresPerIpPerMinute     uint32 = 5
	defaultMaxFailuresPerUserPerMinute   uint32 = 5

	defaultTransactionPictureFileMaxSize uint32 = 10485760 // 10MB
	defaultUserAvatarFileMaxSize         uint32 = 1048576  // 1MB

	defaultImportFileMaxSize uint32 = 10485760 // 10MB

	defaultExchangeRatesDataRequestTimeout uint32 = 10000 // 10 seconds
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

	MaxIdleConnection     uint16
	MaxOpenConnection     uint16
	ConnectionMaxLifeTime uint32
}

// SMTPConfig represents the SMTP setting config
type SMTPConfig struct {
	SMTPHost          string
	SMTPUser          string
	SMTPPasswd        string
	SMTPSkipTLSVerify bool
	FromAddress       string
}

// MinIOConfig represents the MinIO setting config
type MinIOConfig struct {
	Endpoint        string
	Location        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	SkipTLSVerify   bool
	Bucket          string
	RootPath        string
}

// TipConfig represents a tip setting config
type TipConfig struct {
	Enabled              bool
	DefaultContent       string
	MultiLanguageContent map[string]string
}

// NotificationConfig represents a notification setting config
type NotificationConfig struct {
	Enabled              bool
	DefaultContent       string
	MultiLanguageContent map[string]string
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
	HttpPort uint16

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

	// Mail
	EnableSMTP bool
	SMTPConfig *SMTPConfig

	// Log
	LogModes         []string
	EnableConsoleLog bool
	EnableFileLog    bool

	LogLevel           Level
	FileLogPath        string
	RequestFileLogPath string
	QueryFileLogPath   string
	LogFileRotate      bool
	LogFileMaxSize     uint32
	LogFileMaxDays     uint32

	// Storage
	StorageType         string
	LocalFileSystemPath string
	MinIOConfig         *MinIOConfig

	// Uuid
	UuidGeneratorType string
	UuidServerId      uint8

	// Duplicate Checker
	DuplicateCheckerType                            string
	InMemoryDuplicateCheckerCleanupInterval         uint32
	InMemoryDuplicateCheckerCleanupIntervalDuration time.Duration
	EnableDuplicateSubmissionsCheck                 bool
	DuplicateSubmissionsInterval                    uint32
	DuplicateSubmissionsIntervalDuration            time.Duration

	// Cron
	EnableRemoveExpiredTokens        bool
	EnableCreateScheduledTransaction bool

	// Secret
	SecretKeyNoSet                        bool
	SecretKey                             string
	EnableTwoFactor                       bool
	TokenExpiredTime                      uint32
	TokenExpiredTimeDuration              time.Duration
	TokenMinRefreshInterval               uint32
	TemporaryTokenExpiredTime             uint32
	TemporaryTokenExpiredTimeDuration     time.Duration
	EmailVerifyTokenExpiredTime           uint32
	EmailVerifyTokenExpiredTimeDuration   time.Duration
	PasswordResetTokenExpiredTime         uint32
	PasswordResetTokenExpiredTimeDuration time.Duration
	MaxFailuresPerIpPerMinute             uint32
	MaxFailuresPerUserPerMinute           uint32
	EnableRequestIdHeader                 bool

	// User
	EnableUserRegister               bool
	EnableUserVerifyEmail            bool
	EnableUserForceVerifyEmail       bool
	EnableUserForgetPassword         bool
	ForgetPasswordRequireVerifyEmail bool
	EnableTransactionPictures        bool
	MaxTransactionPictureFileSize    uint32
	EnableScheduledTransaction       bool
	AvatarProvider                   core.UserAvatarProviderType
	MaxAvatarFileSize                uint32
	DefaultFeatureRestrictions       core.UserFeatureRestrictions

	// Data
	EnableDataExport  bool
	EnableDataImport  bool
	MaxImportFileSize uint32

	// Tip
	LoginPageTips TipConfig

	// Notification
	AfterRegisterNotification NotificationConfig
	AfterLoginNotification    NotificationConfig
	AfterOpenNotification     NotificationConfig

	// Map
	MapProvider                           string
	EnableMapDataFetchProxy               bool
	MapProxy                              string
	TomTomMapAPIKey                       string
	TianDiTuAPIKey                        string
	GoogleMapAPIKey                       string
	BaiduMapAK                            string
	AmapApplicationKey                    string
	AmapSecurityVerificationMethod        string
	AmapApplicationSecret                 string
	AmapApiExternalProxyUrl               string
	CustomMapTileServerTileLayerUrl       string
	CustomMapTileServerAnnotationLayerUrl string
	CustomMapTileServerMinZoomLevel       uint8
	CustomMapTileServerMaxZoomLevel       uint8
	CustomMapTileServerDefaultZoomLevel   uint8

	// Exchange Rates
	ExchangeRatesDataSource                       string
	ExchangeRatesRequestTimeout                   uint32
	ExchangeRatesRequestTimeoutExceedDefaultValue bool
	ExchangeRatesProxy                            string
	ExchangeRatesSkipTLSVerify                    bool
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

	err = loadMailConfiguration(config, cfgFile, "mail")

	if err != nil {
		return nil, err
	}

	err = loadLogConfiguration(config, cfgFile, "log")

	if err != nil {
		return nil, err
	}

	err = loadStorageConfiguration(config, cfgFile, "storage")

	if err != nil {
		return nil, err
	}

	err = loadUuidConfiguration(config, cfgFile, "uuid")

	if err != nil {
		return nil, err
	}

	err = loadDuplicateCheckerConfiguration(config, cfgFile, "duplicate_checker")

	if err != nil {
		return nil, err
	}

	err = loadCronConfiguration(config, cfgFile, "cron")

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

	err = loadTipConfiguration(config, cfgFile, "tip")

	if err != nil {
		return nil, err
	}

	err = loadNotificationConfiguration(config, cfgFile, "notification")

	if err != nil {
		return nil, err
	}

	err = loadMapConfiguration(config, cfgFile, "map")

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

	if getConfigItemStringValue(configFile, sectionName, "mode") == "production" {
		config.Mode = MODE_PRODUCTION
	} else if getConfigItemStringValue(configFile, sectionName, "mode") == "development" {
		config.Mode = MODE_DEVELOPMENT
	} else {
		return errs.ErrInvalidServerMode
	}

	return nil
}

func loadServerConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "protocol") == "http" {
		config.Protocol = SCHEME_HTTP

		config.HttpAddr = getConfigItemStringValue(configFile, sectionName, "http_addr", defaultHttpAddr)
		config.HttpPort = getConfigItemUint16Value(configFile, sectionName, "http_port", defaultHttpPort)
	} else if getConfigItemStringValue(configFile, sectionName, "protocol") == "https" {
		config.Protocol = SCHEME_HTTPS

		config.HttpAddr = getConfigItemStringValue(configFile, sectionName, "http_addr", defaultHttpAddr)
		config.HttpPort = getConfigItemUint16Value(configFile, sectionName, "http_port", defaultHttpPort)

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

	if dbConfig.DatabaseType != MySqlDbType &&
		dbConfig.DatabaseType != PostgresDbType &&
		dbConfig.DatabaseType != Sqlite3DbType {
		return errs.ErrDatabaseTypeInvalid
	}

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

	dbConfig.MaxIdleConnection = getConfigItemUint16Value(configFile, sectionName, "max_idle_conn", defaultDatabaseMaxIdleConn)
	dbConfig.MaxOpenConnection = getConfigItemUint16Value(configFile, sectionName, "max_open_conn", defaultDatabaseMaxOpenConn)
	dbConfig.ConnectionMaxLifeTime = getConfigItemUint32Value(configFile, sectionName, "conn_max_lifetime", defaultDatabaseConnMaxLifetime)

	config.DatabaseConfig = dbConfig
	config.EnableQueryLog = getConfigItemBoolValue(configFile, sectionName, "log_query", false)
	config.AutoUpdateDatabase = getConfigItemBoolValue(configFile, sectionName, "auto_update_database", true)

	return nil
}

func loadMailConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableSMTP = getConfigItemBoolValue(configFile, sectionName, "enable_smtp", false)

	smtpConfig := &SMTPConfig{}
	smtpConfig.SMTPHost = getConfigItemStringValue(configFile, sectionName, "smtp_host")
	smtpConfig.SMTPUser = getConfigItemStringValue(configFile, sectionName, "smtp_user")
	smtpConfig.SMTPPasswd = getConfigItemStringValue(configFile, sectionName, "smtp_passwd")
	smtpConfig.SMTPSkipTLSVerify = getConfigItemBoolValue(configFile, sectionName, "smtp_skip_tls_verify", false)

	smtpConfig.FromAddress = getConfigItemStringValue(configFile, sectionName, "from_address")

	config.SMTPConfig = smtpConfig

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

	var err error
	config.LogLevel, err = getLogLevel(getConfigItemStringValue(configFile, sectionName, "level"))

	if err != nil {
		return err
	}

	if config.LogLevel != LOGLEVEL_DEBUG &&
		config.LogLevel != LOGLEVEL_INFO &&
		config.LogLevel != LOGLEVEL_WARN &&
		config.LogLevel != LOGLEVEL_ERROR {
		return errs.ErrInvalidLogLevel
	}

	if config.EnableFileLog {
		fileLogPath := getConfigItemStringValue(configFile, sectionName, "log_path")
		finalFileLogPath, _ := getFinalPath(config.WorkingPath, fileLogPath)
		config.FileLogPath = finalFileLogPath

		requestFileLogPath := getConfigItemStringValue(configFile, sectionName, "request_log_path")

		if requestFileLogPath != "" {
			finalRequestFileLogPath, _ := getFinalPath(config.WorkingPath, requestFileLogPath)
			config.RequestFileLogPath = finalRequestFileLogPath
		} else {
			config.RequestFileLogPath = ""
		}

		queryFileLogPath := getConfigItemStringValue(configFile, sectionName, "query_log_path")

		if queryFileLogPath != "" {
			finalQueryFileLogPath, _ := getFinalPath(config.WorkingPath, queryFileLogPath)
			config.QueryFileLogPath = finalQueryFileLogPath
		} else {
			config.QueryFileLogPath = ""
		}

		config.LogFileRotate = getConfigItemBoolValue(configFile, sectionName, "log_file_rotate", false)
		config.LogFileMaxSize = getConfigItemUint32Value(configFile, sectionName, "log_file_max_size", defaultLogFileMaxSize)
		config.LogFileMaxDays = getConfigItemUint32Value(configFile, sectionName, "log_file_max_days", defaultLogFileMaxDays)
	}

	return nil
}

func loadStorageConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "type") == LocalFileSystemObjectStorageType {
		config.StorageType = LocalFileSystemObjectStorageType
	} else if getConfigItemStringValue(configFile, sectionName, "type") == MinIOStorageType {
		config.StorageType = MinIOStorageType
	} else {
		return errs.ErrInvalidStorageType
	}

	localFileSystemRootPath := getConfigItemStringValue(configFile, sectionName, "local_filesystem_path")
	finalLocalFileSystemRootPath, err := getFinalPath(config.WorkingPath, localFileSystemRootPath)
	config.LocalFileSystemPath = finalLocalFileSystemRootPath

	if config.StorageType == LocalFileSystemObjectStorageType && err != nil {
		return errs.ErrInvalidLocalFileSystemStoragePath
	}

	minIOConfig := &MinIOConfig{}
	minIOConfig.Endpoint = getConfigItemStringValue(configFile, sectionName, "minio_endpoint")
	minIOConfig.Location = getConfigItemStringValue(configFile, sectionName, "minio_location")
	minIOConfig.AccessKeyID = getConfigItemStringValue(configFile, sectionName, "minio_access_key_id")
	minIOConfig.SecretAccessKey = getConfigItemStringValue(configFile, sectionName, "minio_secret_access_key")
	minIOConfig.UseSSL = getConfigItemBoolValue(configFile, sectionName, "minio_use_ssl", false)
	minIOConfig.SkipTLSVerify = getConfigItemBoolValue(configFile, sectionName, "minio_skip_tls_verify", false)
	minIOConfig.Bucket = getConfigItemStringValue(configFile, sectionName, "minio_bucket")
	minIOConfig.RootPath = getConfigItemStringValue(configFile, sectionName, "minio_root_path")

	config.MinIOConfig = minIOConfig

	return nil
}

func loadUuidConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "generator_type") == InternalUuidGeneratorType {
		config.UuidGeneratorType = InternalUuidGeneratorType
	} else {
		return errs.ErrInvalidUuidMode
	}

	config.UuidServerId = getConfigItemUint8Value(configFile, sectionName, "server_id", 0)

	return nil
}

func loadDuplicateCheckerConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	if getConfigItemStringValue(configFile, sectionName, "checker_type") == InMemoryDuplicateCheckerType {
		config.DuplicateCheckerType = InMemoryDuplicateCheckerType
	} else {
		return errs.ErrInvalidDuplicateCheckerType
	}

	config.InMemoryDuplicateCheckerCleanupInterval = getConfigItemUint32Value(configFile, sectionName, "cleanup_interval", defaultInMemoryDuplicateCheckerCleanupInterval)

	if config.InMemoryDuplicateCheckerCleanupInterval < 1 {
		return errs.ErrInvalidInMemoryDuplicateCheckerCleanupInterval
	}

	config.InMemoryDuplicateCheckerCleanupIntervalDuration = time.Duration(config.InMemoryDuplicateCheckerCleanupInterval) * time.Second

	duplicateSubmissionsInterval := getConfigItemUint32Value(configFile, sectionName, "duplicate_submissions_interval", defaultDuplicateSubmissionsInterval)

	config.EnableDuplicateSubmissionsCheck = duplicateSubmissionsInterval > 0

	if duplicateSubmissionsInterval < 1 {
		duplicateSubmissionsInterval = defaultDuplicateSubmissionsInterval
	}

	config.DuplicateSubmissionsInterval = duplicateSubmissionsInterval
	config.DuplicateSubmissionsIntervalDuration = time.Duration(config.DuplicateSubmissionsInterval) * time.Second

	return nil
}

func loadCronConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableRemoveExpiredTokens = getConfigItemBoolValue(configFile, sectionName, "enable_remove_expired_tokens", false)
	config.EnableCreateScheduledTransaction = getConfigItemBoolValue(configFile, sectionName, "enable_create_scheduled_transaction", false)

	return nil
}

func loadSecurityConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.SecretKeyNoSet = !getConfigItemIsSet(configFile, sectionName, "secret_key")
	config.SecretKey = getConfigItemStringValue(configFile, sectionName, "secret_key", defaultSecretKey)
	config.EnableTwoFactor = getConfigItemBoolValue(configFile, sectionName, "enable_two_factor", true)

	config.TokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "token_expired_time", defaultTokenExpiredTime)

	if config.TokenExpiredTime < 60 {
		return errs.ErrInvalidTokenExpiredTime
	}

	config.TokenExpiredTimeDuration = time.Duration(config.TokenExpiredTime) * time.Second

	config.TokenMinRefreshInterval = getConfigItemUint32Value(configFile, sectionName, "token_min_refresh_interval", defaultTokenMinRefreshInterval)

	if config.TokenMinRefreshInterval >= config.TokenExpiredTime {
		return errs.ErrInvalidTokenMinRefreshInterval
	}

	config.TemporaryTokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "temporary_token_expired_time", defaultTemporaryTokenExpiredTime)

	if config.TemporaryTokenExpiredTime < 60 {
		return errs.ErrInvalidTemporaryTokenExpiredTime
	}

	config.TemporaryTokenExpiredTimeDuration = time.Duration(config.TemporaryTokenExpiredTime) * time.Second

	config.EmailVerifyTokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "email_verify_token_expired_time", defaultEmailVerifyTokenExpiredTime)

	if config.EmailVerifyTokenExpiredTime < 60 {
		return errs.ErrInvalidEmailVerifyTokenExpiredTime
	}

	config.EmailVerifyTokenExpiredTimeDuration = time.Duration(config.EmailVerifyTokenExpiredTime) * time.Second

	config.PasswordResetTokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "password_reset_token_expired_time", defaultPasswordResetTokenExpiredTime)

	if config.PasswordResetTokenExpiredTime < 60 {
		return errs.ErrInvalidPasswordResetTokenExpiredTime
	}

	config.PasswordResetTokenExpiredTimeDuration = time.Duration(config.PasswordResetTokenExpiredTime) * time.Second

	config.MaxFailuresPerIpPerMinute = getConfigItemUint32Value(configFile, sectionName, "max_failures_per_ip_per_minute", defaultMaxFailuresPerIpPerMinute)
	config.MaxFailuresPerUserPerMinute = getConfigItemUint32Value(configFile, sectionName, "max_failures_per_user_per_minute", defaultMaxFailuresPerUserPerMinute)

	config.EnableRequestIdHeader = getConfigItemBoolValue(configFile, sectionName, "request_id_header", true)

	return nil
}

func loadUserConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableUserRegister = getConfigItemBoolValue(configFile, sectionName, "enable_register", false)
	config.EnableUserVerifyEmail = getConfigItemBoolValue(configFile, sectionName, "enable_email_verify", false)
	config.EnableUserForceVerifyEmail = getConfigItemBoolValue(configFile, sectionName, "enable_force_email_verify", false)
	config.EnableUserForgetPassword = getConfigItemBoolValue(configFile, sectionName, "enable_forget_password", false)
	config.ForgetPasswordRequireVerifyEmail = getConfigItemBoolValue(configFile, sectionName, "forget_password_require_email_verify", false)
	config.EnableTransactionPictures = getConfigItemBoolValue(configFile, sectionName, "enable_transaction_picture", false)
	config.MaxTransactionPictureFileSize = getConfigItemUint32Value(configFile, sectionName, "max_transaction_picture_size", defaultTransactionPictureFileMaxSize)
	config.EnableScheduledTransaction = getConfigItemBoolValue(configFile, sectionName, "enable_scheduled_transaction", false)

	if getConfigItemStringValue(configFile, sectionName, "avatar_provider") == string(core.USER_AVATAR_PROVIDER_INTERNAL) {
		config.AvatarProvider = core.USER_AVATAR_PROVIDER_INTERNAL
	} else if getConfigItemStringValue(configFile, sectionName, "avatar_provider") == string(core.USER_AVATAR_PROVIDER_GRAVATAR) {
		config.AvatarProvider = core.USER_AVATAR_PROVIDER_GRAVATAR
	} else if getConfigItemStringValue(configFile, sectionName, "avatar_provider") == "" {
		config.AvatarProvider = ""
	} else {
		return errs.ErrInvalidAvatarProvider
	}

	config.MaxAvatarFileSize = getConfigItemUint32Value(configFile, sectionName, "max_user_avatar_size", defaultUserAvatarFileMaxSize)
	config.DefaultFeatureRestrictions = core.ParseUserFeatureRestrictions(getConfigItemStringValue(configFile, sectionName, "default_feature_restrictions", ""))

	return nil
}

func loadDataConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableDataExport = getConfigItemBoolValue(configFile, sectionName, "enable_export", false)
	config.EnableDataImport = getConfigItemBoolValue(configFile, sectionName, "enable_import", false)
	config.MaxImportFileSize = getConfigItemUint32Value(configFile, sectionName, "max_import_file_size", defaultImportFileMaxSize)

	return nil
}

func loadTipConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.LoginPageTips = getTipConfiguration(configFile, sectionName, "enable_tips_in_login_page", "login_page_tips_content")

	return nil
}

func loadNotificationConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.AfterRegisterNotification = getNotificationConfiguration(configFile, sectionName, "enable_notification_after_register", "after_register_notification_content")
	config.AfterLoginNotification = getNotificationConfiguration(configFile, sectionName, "enable_notification_after_login", "after_login_notification_content")
	config.AfterOpenNotification = getNotificationConfiguration(configFile, sectionName, "enable_notification_after_open", "after_open_notification_content")

	return nil
}

func loadMapConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	mapProvider := getConfigItemStringValue(configFile, sectionName, "map_provider")

	if mapProvider == "" {
		config.MapProvider = ""
	} else if mapProvider == OpenStreetMapProvider {
		config.MapProvider = OpenStreetMapProvider
	} else if mapProvider == OpenStreetMapHumanitarianStyleProvider {
		config.MapProvider = OpenStreetMapHumanitarianStyleProvider
	} else if mapProvider == OpenTopoMapProvider {
		config.MapProvider = OpenTopoMapProvider
	} else if mapProvider == OPNVKarteMapProvider {
		config.MapProvider = OPNVKarteMapProvider
	} else if mapProvider == CyclOSMMapProvider {
		config.MapProvider = CyclOSMMapProvider
	} else if mapProvider == CartoDBMapProvider {
		config.MapProvider = CartoDBMapProvider
	} else if mapProvider == TomTomMapProvider {
		config.MapProvider = TomTomMapProvider
	} else if mapProvider == TianDiTuProvider {
		config.MapProvider = TianDiTuProvider
	} else if mapProvider == GoogleMapProvider {
		config.MapProvider = GoogleMapProvider
	} else if mapProvider == BaiduMapProvider {
		config.MapProvider = BaiduMapProvider
	} else if mapProvider == AmapProvider {
		config.MapProvider = AmapProvider
	} else if mapProvider == CustomProvider {
		config.MapProvider = CustomProvider
	} else {
		return errs.ErrInvalidMapProvider
	}

	config.EnableMapDataFetchProxy = getConfigItemBoolValue(configFile, sectionName, "map_data_fetch_proxy", false)
	config.MapProxy = getConfigItemStringValue(configFile, sectionName, "proxy", "system")
	config.TomTomMapAPIKey = getConfigItemStringValue(configFile, sectionName, "tomtom_map_api_key")
	config.TianDiTuAPIKey = getConfigItemStringValue(configFile, sectionName, "tianditu_map_app_key")
	config.GoogleMapAPIKey = getConfigItemStringValue(configFile, sectionName, "google_map_api_key")
	config.BaiduMapAK = getConfigItemStringValue(configFile, sectionName, "baidu_map_ak")
	config.AmapApplicationKey = getConfigItemStringValue(configFile, sectionName, "amap_application_key")

	amapSecurityVerificationMethod := getConfigItemStringValue(configFile, sectionName, "amap_security_verification_method")

	if amapSecurityVerificationMethod == AmapSecurityVerificationInternalProxyMethod {
		config.AmapSecurityVerificationMethod = AmapSecurityVerificationInternalProxyMethod
	} else if amapSecurityVerificationMethod == AmapSecurityVerificationExternalProxyMethod {
		config.AmapSecurityVerificationMethod = AmapSecurityVerificationExternalProxyMethod
	} else if amapSecurityVerificationMethod == AmapSecurityVerificationPlainTextMethod {
		config.AmapSecurityVerificationMethod = AmapSecurityVerificationPlainTextMethod
	} else {
		return errs.ErrInvalidAmapSecurityVerificationMethod
	}

	config.AmapApplicationSecret = getConfigItemStringValue(configFile, sectionName, "amap_application_secret")
	config.AmapApiExternalProxyUrl = getConfigItemStringValue(configFile, sectionName, "amap_api_external_proxy_url")

	config.CustomMapTileServerTileLayerUrl = getConfigItemStringValue(configFile, sectionName, "custom_map_tile_server_url")
	config.CustomMapTileServerAnnotationLayerUrl = getConfigItemStringValue(configFile, sectionName, "custom_map_tile_server_annotation_url")
	config.CustomMapTileServerMinZoomLevel = getConfigItemUint8Value(configFile, sectionName, "custom_map_tile_server_min_zoom_level", 1)
	config.CustomMapTileServerMaxZoomLevel = getConfigItemUint8Value(configFile, sectionName, "custom_map_tile_server_max_zoom_level", 18)
	config.CustomMapTileServerDefaultZoomLevel = getConfigItemUint8Value(configFile, sectionName, "custom_map_tile_server_default_zoom_level", 14)

	return nil
}
func loadExchangeRatesConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	dataSource := getConfigItemStringValue(configFile, sectionName, "data_source")

	if dataSource == ReserveBankOfAustraliaDataSource ||
		dataSource == BankOfCanadaDataSource ||
		dataSource == CzechNationalBankDataSource ||
		dataSource == DanmarksNationalbankDataSource ||
		dataSource == EuroCentralBankDataSource ||
		dataSource == NationalBankOfGeorgiaDataSource ||
		dataSource == CentralBankOfHungaryDataSource ||
		dataSource == BankOfIsraelDataSource ||
		dataSource == CentralBankOfMyanmarDataSource ||
		dataSource == NorgesBankDataSource ||
		dataSource == NationalBankOfPolandDataSource ||
		dataSource == NationalBankOfRomaniaDataSource ||
		dataSource == BankOfRussiaDataSource ||
		dataSource == SwissNationalBankDataSource ||
		dataSource == NationalBankOfUkraineDataSource ||
		dataSource == CentralBankOfUzbekistanDataSource ||
		dataSource == InternationalMonetaryFundDataSource {
		config.ExchangeRatesDataSource = dataSource
	} else {
		return errs.ErrInvalidExchangeRatesDataSource
	}

	config.ExchangeRatesProxy = getConfigItemStringValue(configFile, sectionName, "proxy", "system")
	config.ExchangeRatesRequestTimeout = getConfigItemUint32Value(configFile, sectionName, "request_timeout", defaultExchangeRatesDataRequestTimeout)

	if config.ExchangeRatesRequestTimeout > defaultExchangeRatesDataRequestTimeout {
		config.ExchangeRatesRequestTimeoutExceedDefaultValue = true
	}

	config.ExchangeRatesSkipTLSVerify = getConfigItemBoolValue(configFile, sectionName, "skip_tls_verify", false)

	return nil
}

func getWorkingPath() (string, error) {
	workingPath := os.Getenv(ebkWorkDirEnvName)

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

func getTipConfiguration(configFile *ini.File, sectionName string, enableKey string, contentKey string) TipConfig {
	config := TipConfig{
		Enabled:              getConfigItemBoolValue(configFile, sectionName, enableKey, false),
		DefaultContent:       getConfigItemStringValue(configFile, sectionName, contentKey, ""),
		MultiLanguageContent: make(map[string]string),
	}

	for languageTag := range locales.AllLanguages {
		multiLanguageContentKey := strings.ToLower(languageTag)
		multiLanguageContentKey = strings.Replace(multiLanguageContentKey, "-", "_", -1)
		multiLanguageContentKey = contentKey + "_" + multiLanguageContentKey
		content := getConfigItemStringValue(configFile, sectionName, multiLanguageContentKey, "")

		if content != "" {
			config.MultiLanguageContent[languageTag] = content
		}
	}

	return config
}

func getNotificationConfiguration(configFile *ini.File, sectionName string, enableKey string, contentKey string) NotificationConfig {
	config := NotificationConfig{
		Enabled:              getConfigItemBoolValue(configFile, sectionName, enableKey, false),
		DefaultContent:       getConfigItemStringValue(configFile, sectionName, contentKey, ""),
		MultiLanguageContent: make(map[string]string),
	}

	for languageTag := range locales.AllLanguages {
		multiLanguageContentKey := strings.ToLower(languageTag)
		multiLanguageContentKey = strings.Replace(multiLanguageContentKey, "-", "_", -1)
		multiLanguageContentKey = contentKey + "_" + multiLanguageContentKey
		content := getConfigItemStringValue(configFile, sectionName, multiLanguageContentKey, "")

		if content != "" {
			config.MultiLanguageContent[languageTag] = content
		}
	}

	return config
}

func getConfigItemIsSet(configFile *ini.File, sectionName string, itemName string) bool {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		return true
	}

	section := configFile.Section(sectionName)

	if !section.HasKey(itemName) {
		return false
	}

	return section.Key(itemName).String() != ""
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

func getConfigItemUint8Value(configFile *ini.File, sectionName string, itemName string, defaultValue uint8) uint8 {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		value, err := strconv.ParseUint(environmentValue, 10, 8)

		if err == nil {
			return uint8(value)
		}
	}

	section := configFile.Section(sectionName)
	value, err := strconv.ParseUint(section.Key(itemName).String(), 10, 8)

	if err == nil {
		return uint8(value)
	}

	return defaultValue
}

func getConfigItemUint16Value(configFile *ini.File, sectionName string, itemName string, defaultValue uint16) uint16 {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		value, err := strconv.ParseUint(environmentValue, 10, 16)

		if err == nil {
			return uint16(value)
		}
	}

	section := configFile.Section(sectionName)
	value, err := strconv.ParseUint(section.Key(itemName).String(), 10, 16)

	if err == nil {
		return uint16(value)
	}

	return defaultValue
}

func getConfigItemUint32Value(configFile *ini.File, sectionName string, itemName string, defaultValue uint32) uint32 {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		value, err := strconv.ParseUint(environmentValue, 10, 32)

		if err == nil {
			return uint32(value)
		}
	}

	section := configFile.Section(sectionName)
	value, err := strconv.ParseUint(section.Key(itemName).String(), 10, 32)

	if err == nil {
		return uint32(value)
	}

	return defaultValue
}

func getConfigItemBoolValue(configFile *ini.File, sectionName string, itemName string, defaultValue bool) bool {
	environmentKey := getEnvironmentKey(sectionName, itemName)
	environmentValue := os.Getenv(environmentKey)

	if len(environmentValue) > 0 {
		value, err := strconv.ParseBool(environmentValue)

		if err == nil {
			return value
		}
	}

	section := configFile.Section(sectionName)
	return section.Key(itemName).MustBool(defaultValue)
}

func getEnvironmentKey(sectionName string, itemName string) string {
	return fmt.Sprintf("%s_%s_%s", ebkEnvNamePrefix, strings.ToUpper(sectionName), strings.ToUpper(itemName))
}

func getLogLevel(logLevelStr string) (Level, error) {
	if logLevelStr == "debug" {
		return LOGLEVEL_DEBUG, nil
	} else if logLevelStr == "info" {
		return LOGLEVEL_INFO, nil
	} else if logLevelStr == "warn" {
		return LOGLEVEL_WARN, nil
	} else if logLevelStr == "error" {
		return LOGLEVEL_ERROR, nil
	}

	return "", errs.ErrInvalidLogLevel
}

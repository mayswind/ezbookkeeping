package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
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

// Uuid generator types
const (
	InternalUuidGeneratorType string = "internal"
)

// User avatar provider types
const (
	GravatarProvider string = "gravatar"
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
	EuroCentralBankDataSource              string = "euro_central_bank"
	BankOfCanadaDataSource                 string = "bank_of_canada"
	ReserveBankOfAustraliaDataSource       string = "reserve_bank_of_australia"
	CzechNationalBankDataSource            string = "czech_national_bank"
	NationalBankOfPolandDataSource         string = "national_bank_of_poland"
	MonetaryAuthorityOfSingaporeDataSource string = "monetary_authority_of_singapore"
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

	defaultLogMode  string = "console"
	defaultLoglevel Level  = LOGLEVEL_INFO

	defaultSecretKey                     string = "ezbookkeeping"
	defaultTokenExpiredTime              uint32 = 604800 // 7 days
	defaultTemporaryTokenExpiredTime     uint32 = 300    // 5 minutes
	defaultEmailVerifyTokenExpiredTime   uint32 = 3600   // 60 minutes
	defaultPasswordResetTokenExpiredTime uint32 = 3600   // 60 minutes

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

	LogLevel    Level
	FileLogPath string

	// Uuid
	UuidGeneratorType string
	UuidServerId      uint8

	// Secret
	SecretKeyNoSet                        bool
	SecretKey                             string
	EnableTwoFactor                       bool
	TokenExpiredTime                      uint32
	TokenExpiredTimeDuration              time.Duration
	TemporaryTokenExpiredTime             uint32
	TemporaryTokenExpiredTimeDuration     time.Duration
	EmailVerifyTokenExpiredTime           uint32
	EmailVerifyTokenExpiredTimeDuration   time.Duration
	PasswordResetTokenExpiredTime         uint32
	PasswordResetTokenExpiredTimeDuration time.Duration
	EnableRequestIdHeader                 bool

	// User
	EnableUserRegister               bool
	EnableUserVerifyEmail            bool
	EnableUserForceVerifyEmail       bool
	EnableUserForgetPassword         bool
	ForgetPasswordRequireVerifyEmail bool
	AvatarProvider                   string

	// Data
	EnableDataExport bool

	// Map
	MapProvider                         string
	EnableMapDataFetchProxy             bool
	MapProxy                            string
	TomTomMapAPIKey                     string
	GoogleMapAPIKey                     string
	BaiduMapAK                          string
	AmapApplicationKey                  string
	AmapSecurityVerificationMethod      string
	AmapApplicationSecret               string
	AmapApiExternalProxyUrl             string
	CustomMapTileServerUrl              string
	CustomMapTileServerMinZoomLevel     uint8
	CustomMapTileServerMaxZoomLevel     uint8
	CustomMapTileServerDefaultZoomLevel uint8

	// Exchange Rates
	ExchangeRatesDataSource     string
	ExchangeRatesRequestTimeout uint32
	ExchangeRatesProxy          string
	ExchangeRatesSkipTLSVerify  bool
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

	config.LogLevel = getLogLevel(getConfigItemStringValue(configFile, sectionName, "level"), defaultLoglevel)

	if config.EnableFileLog {
		fileLogPath := getConfigItemStringValue(configFile, sectionName, "log_path")
		finalFileLogPath, _ := getFinalPath(config.WorkingPath, fileLogPath)
		config.FileLogPath = finalFileLogPath
	}

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

func loadSecurityConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.SecretKeyNoSet = !getConfigItemIsSet(configFile, sectionName, "secret_key")
	config.SecretKey = getConfigItemStringValue(configFile, sectionName, "secret_key", defaultSecretKey)
	config.EnableTwoFactor = getConfigItemBoolValue(configFile, sectionName, "enable_two_factor", true)

	config.TokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "token_expired_time", defaultTokenExpiredTime)
	config.TokenExpiredTimeDuration = time.Duration(config.TokenExpiredTime) * time.Second

	config.TemporaryTokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "temporary_token_expired_time", defaultTemporaryTokenExpiredTime)
	config.TemporaryTokenExpiredTimeDuration = time.Duration(config.TemporaryTokenExpiredTime) * time.Second

	config.EmailVerifyTokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "email_verify_token_expired_time", defaultEmailVerifyTokenExpiredTime)
	config.EmailVerifyTokenExpiredTimeDuration = time.Duration(config.EmailVerifyTokenExpiredTime) * time.Second

	config.PasswordResetTokenExpiredTime = getConfigItemUint32Value(configFile, sectionName, "password_reset_token_expired_time", defaultPasswordResetTokenExpiredTime)
	config.PasswordResetTokenExpiredTimeDuration = time.Duration(config.PasswordResetTokenExpiredTime) * time.Second

	config.EnableRequestIdHeader = getConfigItemBoolValue(configFile, sectionName, "request_id_header", true)

	return nil
}

func loadUserConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableUserRegister = getConfigItemBoolValue(configFile, sectionName, "enable_register", false)
	config.EnableUserVerifyEmail = getConfigItemBoolValue(configFile, sectionName, "enable_email_verify", false)
	config.EnableUserForceVerifyEmail = getConfigItemBoolValue(configFile, sectionName, "enable_force_email_verify", false)
	config.EnableUserForgetPassword = getConfigItemBoolValue(configFile, sectionName, "enable_forget_password", false)
	config.ForgetPasswordRequireVerifyEmail = getConfigItemBoolValue(configFile, sectionName, "forget_password_require_email_verify", false)

	if getConfigItemStringValue(configFile, sectionName, "avatar_provider") == "" {
		config.AvatarProvider = ""
	} else if getConfigItemStringValue(configFile, sectionName, "avatar_provider") == GravatarProvider {
		config.AvatarProvider = GravatarProvider
	}

	return nil
}

func loadDataConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	config.EnableDataExport = getConfigItemBoolValue(configFile, sectionName, "enable_export", false)

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

	config.CustomMapTileServerUrl = getConfigItemStringValue(configFile, sectionName, "custom_map_tile_server_url")
	config.CustomMapTileServerMinZoomLevel = getConfigItemUint8Value(configFile, sectionName, "custom_map_tile_server_min_zoom_level", 1)
	config.CustomMapTileServerMaxZoomLevel = getConfigItemUint8Value(configFile, sectionName, "custom_map_tile_server_max_zoom_level", 18)
	config.CustomMapTileServerDefaultZoomLevel = getConfigItemUint8Value(configFile, sectionName, "custom_map_tile_server_default_zoom_level", 14)

	return nil
}
func loadExchangeRatesConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	dataSource := getConfigItemStringValue(configFile, sectionName, "data_source")

	if dataSource == EuroCentralBankDataSource {
		config.ExchangeRatesDataSource = EuroCentralBankDataSource
	} else if dataSource == BankOfCanadaDataSource {
		config.ExchangeRatesDataSource = BankOfCanadaDataSource
	} else if dataSource == ReserveBankOfAustraliaDataSource {
		config.ExchangeRatesDataSource = ReserveBankOfAustraliaDataSource
	} else if dataSource == CzechNationalBankDataSource {
		config.ExchangeRatesDataSource = CzechNationalBankDataSource
	} else if dataSource == NationalBankOfPolandDataSource {
		config.ExchangeRatesDataSource = NationalBankOfPolandDataSource
	} else if dataSource == MonetaryAuthorityOfSingaporeDataSource {
		config.ExchangeRatesDataSource = MonetaryAuthorityOfSingaporeDataSource
	} else {
		return errs.ErrInvalidExchangeRatesDataSource
	}

	config.ExchangeRatesProxy = getConfigItemStringValue(configFile, sectionName, "proxy", "system")
	config.ExchangeRatesRequestTimeout = getConfigItemUint32Value(configFile, sectionName, "request_timeout", defaultExchangeRatesDataRequestTimeout)
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

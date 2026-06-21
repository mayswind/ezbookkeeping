package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v3"

	"github.com/mayswind/ezbookkeeping/pkg/api"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cron"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/mcp"
	"github.com/mayswind/ezbookkeeping/pkg/middlewares"
	"github.com/mayswind/ezbookkeeping/pkg/requestid"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// WebServer represents the server command
var WebServer = &cli.Command{
	Name:  "server",
	Usage: "ezBookkeeping web server operation",
	Commands: []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run ezBookkeeping web server",
			Action: bindAction(startWebServer),
		},
	},
}

func startWebServer(c *core.CliContext) error {
	config, err := initializeSystem(c)

	if err != nil {
		return err
	}

	log.BootInfof(c, "[webserver.startWebServer] static root path is %s", config.StaticRootPath)

	if config.AutoUpdateDatabase {
		err = updateAllDatabaseTablesStructure(c)

		if err != nil {
			log.BootErrorf(c, "[webserver.startWebServer] update database table structure failed, because %s", err.Error())
			return err
		}
	}

	err = requestid.InitializeRequestIdGenerator(c, config)

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] initializes requestid generator failed, because %s", err.Error())
		return err
	}

	err = mcp.InitializeMCPHandlers(config)

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] initializes mcp handlers failed, because %s", err.Error())
		return err
	}

	err = oauth2.InitializeOAuth2Provider(config)

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] initializes oauth 2.0 provider failed, because %s", err.Error())
		return err
	}

	err = cron.InitializeCronJobSchedulerContainer(c, config, true)

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] initializes cron job scheduler failed, because %s", err.Error())
		return err
	}

	serverInfo := fmt.Sprintf("current server id is %d, current instance id is %d", requestid.Container.GetCurrentServerUniqId(), requestid.Container.GetCurrentInstanceUniqId())
	uuidServerInfo := ""
	if config.UuidGeneratorType == settings.InternalUuidGeneratorType {
		uuidServerInfo = fmt.Sprintf(", current uuid server id is %d", config.UuidServerId)
	}

	log.BootInfof(c, "[webserver.startWebServer] %s%s", serverInfo, uuidServerInfo)

	if config.Mode == settings.MODE_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	workboxFileNames := utils.ListFileNamesWithPrefixAndSuffix(config.StaticRootPath, "workbox-", ".js")

	router := gin.New()
	router.Use(bindMiddleware(middlewares.Recovery, config))

	err = router.SetTrustedProxies(config.TrustedProxyTextualIPs)

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] set trusted proxy failed, because %s", err.Error())
		return err
	}

	if config.EnableGZip {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("notBlank", validators.NotBlank)
		_ = v.RegisterValidation("validUsername", validators.ValidUsername)
		_ = v.RegisterValidation("validEmail", validators.ValidEmail)
		_ = v.RegisterValidation("validNickname", validators.ValidNickname)
		_ = v.RegisterValidation("validCurrency", validators.ValidCurrency)
		_ = v.RegisterValidation("validHexRGBColor", validators.ValidHexRGBColor)
		_ = v.RegisterValidation("validAmountFilter", validators.ValidAmountFilter)
		_ = v.RegisterValidation("validTagFilter", validators.ValidTagFilter)
		_ = v.RegisterValidation("validFiscalYearStart", validators.ValidateFiscalYearStart)
	}

	router.NoRoute(bindApi(api.Default.ApiNotFound, config))
	router.NoMethod(bindApi(api.Default.MethodNotAllowed, config))

	serverSettingsCacheStore := persistence.NewInMemoryStore(time.Minute)

	router.StaticFile("/", filepath.Join(config.StaticRootPath, "index.html"))
	router.Static("/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/fonts", filepath.Join(config.StaticRootPath, "fonts"))

	router.StaticFile("robots.txt", filepath.Join(config.StaticRootPath, "robots.txt"))
	router.StaticFile("favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))
	router.StaticFile("sw.js", filepath.Join(config.StaticRootPath, "sw.js"))
	router.GET("/server_settings.js", bindCachedJs(api.ServerSettings.ServerSettingsJavascriptHandler, config, serverSettingsCacheStore))

	for i := 0; i < len(workboxFileNames); i++ {
		router.StaticFile("/"+workboxFileNames[i], filepath.Join(config.StaticRootPath, workboxFileNames[i]))
	}

	router.StaticFile("/mobile", filepath.Join(config.StaticRootPath, "mobile.html"))
	router.Match([]string{http.MethodHead, http.MethodGet}, "/mobile#/*fragment", bindLocalFile(filepath.Join(config.StaticRootPath, "mobile.html")))  // add compatibility for browsers that send the full URL with the fragment to the server
	router.Match([]string{http.MethodHead, http.MethodGet}, "/mobile#!/*fragment", bindLocalFile(filepath.Join(config.StaticRootPath, "mobile.html"))) // add compatibility for browsers that send the full URL with the fragment to the server
	router.Static("/mobile/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/mobile/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/mobile/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/mobile/fonts", filepath.Join(config.StaticRootPath, "fonts"))
	router.StaticFile("/mobile/favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("/mobile/favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("/mobile/touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("/mobile/manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))
	router.StaticFile("/mobile/sw.js", filepath.Join(config.StaticRootPath, "sw.js"))
	router.GET("/mobile/server_settings.js", bindCachedJs(api.ServerSettings.ServerSettingsJavascriptHandler, config, serverSettingsCacheStore))

	for i := 0; i < len(workboxFileNames); i++ {
		router.StaticFile("/mobile/"+workboxFileNames[i], filepath.Join(config.StaticRootPath, workboxFileNames[i]))
	}

	router.StaticFile("/desktop", filepath.Join(config.StaticRootPath, "desktop.html"))
	router.Match([]string{http.MethodHead, http.MethodGet}, "/desktop#/*fragment", bindLocalFile(filepath.Join(config.StaticRootPath, "desktop.html"))) // add compatibility for browsers that send the full URL with the fragment to the server
	router.Static("/desktop/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/desktop/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/desktop/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/desktop/fonts", filepath.Join(config.StaticRootPath, "fonts"))
	router.StaticFile("/desktop/favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("/desktop/favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("/desktop/touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("/desktop/manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))
	router.StaticFile("/desktop/sw.js", filepath.Join(config.StaticRootPath, "sw.js"))
	router.GET("/desktop/server_settings.js", bindCachedJs(api.ServerSettings.ServerSettingsJavascriptHandler, config, serverSettingsCacheStore))

	for i := 0; i < len(workboxFileNames); i++ {
		router.StaticFile("/desktop/"+workboxFileNames[i], filepath.Join(config.StaticRootPath, workboxFileNames[i]))
	}

	if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
		avatarRoute := router.Group("/avatar")
		avatarRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString(config), config))
		{
			avatarRoute.GET("/:fileName", bindImage(api.Users.UserGetAvatarHandler, config))
		}
	}

	if config.EnableTransactionPictures {
		pictureRoute := router.Group("/pictures")
		pictureRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString(config), config))
		{
			pictureRoute.GET("/:fileName", bindImage(api.TransactionPictures.TransactionPictureGetHandler, config))
		}
	}

	router.GET("/healthz.json", bindApi(api.Healths.HealthStatusHandler, config))

	proxyRoute := router.Group("/proxy")
	proxyRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString(config), config))
	{
		if config.EnableMapDataFetchProxy {
			if config.MapProvider == settings.OpenStreetMapProvider ||
				config.MapProvider == settings.OpenStreetMapHumanitarianStyleProvider ||
				config.MapProvider == settings.OpenTopoMapProvider ||
				config.MapProvider == settings.OPNVKarteMapProvider ||
				config.MapProvider == settings.CyclOSMMapProvider ||
				config.MapProvider == settings.CartoDBMapProvider ||
				config.MapProvider == settings.TomTomMapProvider ||
				config.MapProvider == settings.TianDiTuProvider ||
				config.MapProvider == settings.CustomProvider {
				proxyRoute.GET("/map/tile/:zoomLevel/:coordinateX/:fileName", bindProxy(api.MapImages.MapTileImageProxyHandler, config))
			}

			if config.MapProvider == settings.TianDiTuProvider ||
				(config.MapProvider == settings.CustomProvider && config.CustomMapTileServerAnnotationLayerUrl != "") {
				proxyRoute.GET("/map/annotation/:zoomLevel/:coordinateX/:fileName", bindProxy(api.MapImages.MapAnnotationImageProxyHandler, config))
			}
		}
	}

	if config.MapProvider == settings.AmapProvider && config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationInternalProxyMethod {
		amapApiProxyRoute := router.Group("/_AMapService")
		amapApiProxyRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByCookie(config), config))
		{
			amapApiProxyRoute.GET("/*action", bindProxy(api.AmapApis.AmapApiProxyHandler, config))
		}
	}

	qrCodeRoute := router.Group("/qrcode")
	qrCodeRoute.Use(bindMiddleware(middlewares.RequestId(config), config))
	{
		qrCodeCacheStore := persistence.NewInMemoryStore(time.Minute)
		qrCodeRoute.GET("/mobile_url.png", bindCachedImage(api.QrCodes.MobileUrlQrCodeHandler, config, qrCodeCacheStore))
	}

	if config.EnableMCPServer {
		mcpRoute := router.Group("/mcp")
		mcpRoute.Use(bindMiddleware(middlewares.RequestId(config), config))
		mcpRoute.Use(bindMiddleware(middlewares.RequestLog, config))
		mcpRoute.Use(bindMiddleware(middlewares.MCPServerIpLimit(config), config))
		mcpRoute.Use(bindMiddleware(middlewares.JWTMCPAuthorization(config), config))
		{
			mcpRoute.POST("", bindJSONRPCApi(map[string]core.JSONRPCApiHandlerFunc{
				"initialize":     api.ModelContextProtocols.InitializeHandler,
				"resources/list": api.ModelContextProtocols.ListResourcesHandler,
				"resources/read": api.ModelContextProtocols.ReadResourceHandler,
				"tools/list":     api.ModelContextProtocols.ListToolsHandler,
				"tools/call":     api.ModelContextProtocols.CallToolHandler,
				"ping":           api.ModelContextProtocols.PingHandler,
			}, map[string]int{
				"notifications/initialized": http.StatusAccepted,
			}, config))
			mcpRoute.GET("", bindApi(api.Default.MethodNotAllowed, config))
		}
	}

	if config.EnableOAuth2Login {
		oauth2Route := router.Group("/oauth2")
		oauth2Route.Use(bindMiddleware(middlewares.RequestId(config), config))
		oauth2Route.Use(bindMiddleware(middlewares.RequestLog, config))
		{
			oauth2Route.GET("/login", bindRedirect(api.OAuth2Authentications.LoginHandler, config))
			oauth2Route.GET("/callback", bindRedirect(api.OAuth2Authentications.CallbackHandler, config))
		}
	}

	apiRoute := router.Group("/api")

	apiRoute.Use(bindMiddleware(middlewares.RequestId(config), config))
	apiRoute.Use(bindMiddleware(middlewares.RequestLog, config))
	{
		if config.EnableInternalAuth {
			apiRoute.POST("/authorize.json", bindApiWithTokenUpdate(api.Authorizations.AuthorizeHandler, config))
		}

		if config.EnableInternalAuth && config.EnableTwoFactor {
			twoFactorRoute := apiRoute.Group("/2fa")
			twoFactorRoute.Use(bindMiddleware(middlewares.JWTTwoFactorAuthorization(config), config))
			{
				twoFactorRoute.POST("/authorize.json", bindApiWithTokenUpdate(api.Authorizations.TwoFactorAuthorizeHandler, config))
				twoFactorRoute.POST("/recovery.json", bindApiWithTokenUpdate(api.Authorizations.TwoFactorAuthorizeByRecoveryCodeHandler, config))
			}
		}

		if config.EnableOAuth2Login {
			oauth2Route := apiRoute.Group("/oauth2")
			oauth2Route.Use(bindMiddleware(middlewares.JWTOAuth2CallbackAuthorization(config), config))
			{
				oauth2Route.POST("/authorize.json", bindApiWithTokenUpdate(api.Authorizations.OAuth2CallbackAuthorizeHandler, config))
			}
		}

		if config.EnableInternalAuth && config.EnableUserRegister {
			apiRoute.POST("/register.json", bindApiWithTokenUpdate(api.Users.UserRegisterHandler, config))
		}

		if config.EnableUserVerifyEmail {
			apiRoute.POST("/verify_email/resend.json", bindApi(api.Users.UserSendVerifyEmailByUnloginUserHandler, config))

			emailVerifyRoute := apiRoute.Group("/verify_email")
			emailVerifyRoute.Use(bindMiddleware(middlewares.JWTEmailVerifyAuthorization(config), config))
			{
				emailVerifyRoute.POST("/by_token.json", bindApi(api.Users.UserEmailVerifyHandler, config))
			}
		}

		if config.EnableInternalAuth && config.EnableUserForgetPassword {
			apiRoute.POST("/forget_password/request.json", bindApi(api.ForgetPasswords.UserForgetPasswordRequestHandler, config))

			resetPasswordRoute := apiRoute.Group("/forget_password/reset")
			resetPasswordRoute.Use(bindMiddleware(middlewares.JWTResetPasswordAuthorization(config), config))
			{
				resetPasswordRoute.POST("/by_token.json", bindApi(api.ForgetPasswords.UserResetPasswordHandler, config))
			}
		}

		apiRoute.GET("/logout.json", bindApiWithTokenUpdate(api.Tokens.TokenRevokeCurrentHandler, config))

		apiV1Route := apiRoute.Group("/v1")
		apiV1Route.Use(bindMiddleware(middlewares.JWTAuthorization(config), config))
		apiV1Route.Use(bindMiddleware(middlewares.APITokenIpLimit(config), config))
		{
			// Tokens
			apiV1Route.GET("/tokens/list.json", bindApi(api.Tokens.TokenListHandler, config))
			apiV1Route.POST("/tokens/generate/api.json", bindApi(api.Tokens.TokenGenerateAPIHandler, config))
			apiV1Route.POST("/tokens/generate/mcp.json", bindApi(api.Tokens.TokenGenerateMCPHandler, config))
			apiV1Route.POST("/tokens/revoke.json", bindApi(api.Tokens.TokenRevokeHandler, config))
			apiV1Route.POST("/tokens/revoke_all.json", bindApi(api.Tokens.TokenRevokeAllHandler, config))
			apiV1Route.POST("/tokens/refresh.json", bindApiWithTokenUpdate(api.Tokens.TokenRefreshHandler, config))

			// Users
			apiV1Route.GET("/users/profile/get.json", bindApi(api.Users.UserProfileHandler, config))
			apiV1Route.POST("/users/profile/update.json", bindApiWithTokenUpdate(api.Users.UserUpdateProfileHandler, config))

			if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
				apiV1Route.POST("/users/avatar/update.json", bindApi(api.Users.UserUpdateAvatarHandler, config))
				apiV1Route.POST("/users/avatar/remove.json", bindApi(api.Users.UserRemoveAvatarHandler, config))
			}

			if config.EnableUserVerifyEmail {
				apiV1Route.POST("/users/verify_email/resend.json", bindApi(api.Users.UserSendVerifyEmailByLoginedUserHandler, config))
			}

			// External Authentications
			if config.EnableOAuth2Login {
				apiV1Route.GET("/users/external_auth/list.json", bindApi(api.UserExternalAuths.ExternalAuthListHandler, config))
				apiV1Route.POST("/users/external_auth/unlink.json", bindApi(api.UserExternalAuths.UnlinkExternalAuthHandler, config))
			}

			// Application Cloud Settings
			apiV1Route.GET("/users/settings/cloud/get.json", bindApi(api.UserApplicationCloudSettings.ApplicationSettingsGetHandler, config))
			apiV1Route.POST("/users/settings/cloud/update.json", bindApi(api.UserApplicationCloudSettings.ApplicationSettingsUpdateHandler, config))
			apiV1Route.POST("/users/settings/cloud/disable.json", bindApi(api.UserApplicationCloudSettings.ApplicationSettingsDisableHandler, config))

			// Two-Factor Authorization
			if config.EnableTwoFactor {
				apiV1Route.GET("/users/2fa/status.json", bindApi(api.TwoFactorAuthorizations.TwoFactorStatusHandler, config))
				apiV1Route.POST("/users/2fa/enable/request.json", bindApi(api.TwoFactorAuthorizations.TwoFactorEnableRequestHandler, config))
				apiV1Route.POST("/users/2fa/enable/confirm.json", bindApiWithTokenUpdate(api.TwoFactorAuthorizations.TwoFactorEnableConfirmHandler, config))
				apiV1Route.POST("/users/2fa/disable.json", bindApi(api.TwoFactorAuthorizations.TwoFactorDisableHandler, config))
				apiV1Route.POST("/users/2fa/recovery/regenerate.json", bindApi(api.TwoFactorAuthorizations.TwoFactorRecoveryCodeRegenerateHandler, config))
			}

			// Data
			apiV1Route.GET("/data/statistics.json", bindApi(api.DataManagements.DataStatisticsHandler, config))
			apiV1Route.POST("/data/clear/all.json", bindApi(api.DataManagements.ClearAllDataHandler, config))
			apiV1Route.POST("/data/clear/transactions.json", bindApi(api.DataManagements.ClearAllTransactionsHandler, config))
			apiV1Route.POST("/data/clear/transactions/by_account.json", bindApi(api.DataManagements.ClearAllTransactionsByAccountHandler, config))

			if config.EnableDataExport {
				apiV1Route.GET("/data/export.csv", bindCsv(api.DataManagements.ExportDataToEzbookkeepingCSVHandler, config))
				apiV1Route.GET("/data/export.tsv", bindTsv(api.DataManagements.ExportDataToEzbookkeepingTSVHandler, config))
			}

			// Accounts
			apiV1Route.GET("/accounts/list.json", bindApi(api.Accounts.AccountListHandler, config))
			apiV1Route.GET("/accounts/get.json", bindApi(api.Accounts.AccountGetHandler, config))
			apiV1Route.POST("/accounts/add.json", bindApi(api.Accounts.AccountCreateHandler, config))
			apiV1Route.POST("/accounts/modify.json", bindApi(api.Accounts.AccountModifyHandler, config))
			apiV1Route.POST("/accounts/update/last_reconciled_time.json", bindApi(api.Accounts.AccountUpdateLastReconciledTimeHandler, config))
			apiV1Route.POST("/accounts/hide.json", bindApi(api.Accounts.AccountHideHandler, config))
			apiV1Route.POST("/accounts/move.json", bindApi(api.Accounts.AccountMoveHandler, config))
			apiV1Route.POST("/accounts/delete.json", bindApi(api.Accounts.AccountDeleteHandler, config))
			apiV1Route.POST("/accounts/sub_account/delete.json", bindApi(api.Accounts.SubAccountDeleteHandler, config))

			// Transactions
			apiV1Route.GET("/transactions/count.json", bindApi(api.Transactions.TransactionCountHandler, config))
			apiV1Route.GET("/transactions/list.json", bindApi(api.Transactions.TransactionListHandler, config))
			apiV1Route.GET("/transactions/list/by_month.json", bindApi(api.Transactions.TransactionMonthListHandler, config))
			apiV1Route.GET("/transactions/list/all.json", bindApi(api.Transactions.TransactionListAllHandler, config))
			apiV1Route.GET("/transactions/reconciliation_statements.json", bindApi(api.Transactions.TransactionReconciliationStatementHandler, config))
			apiV1Route.GET("/transactions/statistics.json", bindApi(api.Transactions.TransactionStatisticsHandler, config))
			apiV1Route.GET("/transactions/statistics/trends.json", bindApi(api.Transactions.TransactionStatisticsTrendsHandler, config))
			apiV1Route.GET("/transactions/statistics/asset_trends.json", bindApi(api.Transactions.TransactionStatisticsAssetTrendsHandler, config))
			apiV1Route.GET("/transactions/amounts.json", bindApi(api.Transactions.TransactionAmountsHandler, config))
			apiV1Route.GET("/transactions/get.json", bindApi(api.Transactions.TransactionGetHandler, config))
			apiV1Route.POST("/transactions/add.json", bindApi(api.Transactions.TransactionCreateHandler, config))
			apiV1Route.POST("/transactions/modify.json", bindApi(api.Transactions.TransactionModifyHandler, config))
			apiV1Route.POST("/transactions/batch_update/category.json", bindApi(api.Transactions.TransactionBatchUpdateCategoriesHandler, config))
			apiV1Route.POST("/transactions/batch_update/account.json", bindApi(api.Transactions.TransactionBatchUpdateAccountsHandler, config))
			apiV1Route.POST("/transactions/batch_update/tag/add.json", bindApi(api.Transactions.TransactionBatchAddTagsHandler, config))
			apiV1Route.POST("/transactions/batch_update/tag/remove.json", bindApi(api.Transactions.TransactionBatchRemoveTagsHandler, config))
			apiV1Route.POST("/transactions/batch_update/tag/clear.json", bindApi(api.Transactions.TransactionBatchClearTagsHandler, config))
			apiV1Route.POST("/transactions/move/all.json", bindApi(api.Transactions.TransactionMoveAllBetweenAccountsHandler, config))
			apiV1Route.POST("/transactions/delete.json", bindApi(api.Transactions.TransactionDeleteHandler, config))
			apiV1Route.POST("/transactions/batch_delete.json", bindApi(api.Transactions.TransactionBatchDeleteHandler, config))

			if config.EnableDataImport {
				apiV1Route.POST("/transactions/parse_custom_file.json", bindApi(api.Transactions.TransactionParseImportCustomFileDataHandler, config))
				apiV1Route.POST("/transactions/parse_import.json", bindApi(api.Transactions.TransactionParseImportFileHandler, config))
				apiV1Route.POST("/transactions/import.json", bindApi(api.Transactions.TransactionImportHandler, config))
				apiV1Route.GET("/transactions/import/process.json", bindApi(api.Transactions.TransactionImportProcessHandler, config))
			}

			// Transaction Pictures
			if config.EnableTransactionPictures {
				apiV1Route.POST("/transaction/pictures/upload.json", bindApi(api.TransactionPictures.TransactionPictureUploadHandler, config))
				apiV1Route.POST("/transaction/pictures/remove_unused.json", bindApi(api.TransactionPictures.TransactionPictureRemoveUnusedHandler, config))
			}

			// Transaction Categories
			apiV1Route.GET("/transaction/categories/list.json", bindApi(api.TransactionCategories.CategoryListHandler, config))
			apiV1Route.GET("/transaction/categories/get.json", bindApi(api.TransactionCategories.CategoryGetHandler, config))
			apiV1Route.POST("/transaction/categories/add.json", bindApi(api.TransactionCategories.CategoryCreateHandler, config))
			apiV1Route.POST("/transaction/categories/add_batch.json", bindApi(api.TransactionCategories.CategoryCreateBatchHandler, config))
			apiV1Route.POST("/transaction/categories/modify.json", bindApi(api.TransactionCategories.CategoryModifyHandler, config))
			apiV1Route.POST("/transaction/categories/hide.json", bindApi(api.TransactionCategories.CategoryHideHandler, config))
			apiV1Route.POST("/transaction/categories/move.json", bindApi(api.TransactionCategories.CategoryMoveHandler, config))
			apiV1Route.POST("/transaction/categories/delete.json", bindApi(api.TransactionCategories.CategoryDeleteHandler, config))

			// Transaction Tag Groups
			apiV1Route.GET("/transaction/tags/groups/list.json", bindApi(api.TransactionTagGroups.TagGroupListHandler, config))
			apiV1Route.GET("/transaction/tags/groups/get.json", bindApi(api.TransactionTagGroups.TagGroupGetHandler, config))
			apiV1Route.POST("/transaction/tags/groups/add.json", bindApi(api.TransactionTagGroups.TagGroupCreateHandler, config))
			apiV1Route.POST("/transaction/tags/groups/modify.json", bindApi(api.TransactionTagGroups.TagGroupModifyHandler, config))
			apiV1Route.POST("/transaction/tags/groups/move.json", bindApi(api.TransactionTagGroups.TagGroupMoveHandler, config))
			apiV1Route.POST("/transaction/tags/groups/delete.json", bindApi(api.TransactionTagGroups.TagGroupDeleteHandler, config))

			// Transaction Tags
			apiV1Route.GET("/transaction/tags/list.json", bindApi(api.TransactionTags.TagListHandler, config))
			apiV1Route.GET("/transaction/tags/get.json", bindApi(api.TransactionTags.TagGetHandler, config))
			apiV1Route.POST("/transaction/tags/add.json", bindApi(api.TransactionTags.TagCreateHandler, config))
			apiV1Route.POST("/transaction/tags/add_batch.json", bindApi(api.TransactionTags.TagCreateBatchHandler, config))
			apiV1Route.POST("/transaction/tags/modify.json", bindApi(api.TransactionTags.TagModifyHandler, config))
			apiV1Route.POST("/transaction/tags/hide.json", bindApi(api.TransactionTags.TagHideHandler, config))
			apiV1Route.POST("/transaction/tags/move.json", bindApi(api.TransactionTags.TagMoveHandler, config))
			apiV1Route.POST("/transaction/tags/delete.json", bindApi(api.TransactionTags.TagDeleteHandler, config))

			// Transaction Templates
			apiV1Route.GET("/transaction/templates/list.json", bindApi(api.TransactionTemplates.TemplateListHandler, config))
			apiV1Route.GET("/transaction/templates/get.json", bindApi(api.TransactionTemplates.TemplateGetHandler, config))
			apiV1Route.POST("/transaction/templates/add.json", bindApi(api.TransactionTemplates.TemplateCreateHandler, config))
			apiV1Route.POST("/transaction/templates/modify.json", bindApi(api.TransactionTemplates.TemplateModifyHandler, config))
			apiV1Route.POST("/transaction/templates/hide.json", bindApi(api.TransactionTemplates.TemplateHideHandler, config))
			apiV1Route.POST("/transaction/templates/move.json", bindApi(api.TransactionTemplates.TemplateMoveHandler, config))
			apiV1Route.POST("/transaction/templates/delete.json", bindApi(api.TransactionTemplates.TemplateDeleteHandler, config))

			// Insights Explorers
			apiV1Route.GET("/insights/explorers/list.json", bindApi(api.InsightsExplorers.InsightsExplorerListHandler, config))
			apiV1Route.GET("/insights/explorers/get.json", bindApi(api.InsightsExplorers.InsightsExplorerGetHandler, config))
			apiV1Route.POST("/insights/explorers/add.json", bindApi(api.InsightsExplorers.InsightsExplorerCreateHandler, config))
			apiV1Route.POST("/insights/explorers/modify.json", bindApi(api.InsightsExplorers.InsightsExplorerModifyHandler, config))
			apiV1Route.POST("/insights/explorers/hide.json", bindApi(api.InsightsExplorers.InsightsExplorerHideHandler, config))
			apiV1Route.POST("/insights/explorers/move.json", bindApi(api.InsightsExplorers.InsightsExplorerMoveHandler, config))
			apiV1Route.POST("/insights/explorers/delete.json", bindApi(api.InsightsExplorers.InsightsExplorerDeleteHandler, config))

			// Large Language Models
			if config.TextRecognitionLLMConfig != nil && config.TextRecognitionLLMConfig.LLMProvider != "" {
				if config.TransactionFromAITextRecognition {
					apiV1Route.POST("/llm/transactions/recognize_text.json", bindApi(api.LargeLanguageModels.RecognizeTransactionTextHandler, config))
				}
			}

			if config.ReceiptImageRecognitionLLMConfig != nil && config.ReceiptImageRecognitionLLMConfig.LLMProvider != "" {
				if config.TransactionFromAIImageRecognition {
					apiV1Route.POST("/llm/transactions/recognize_receipt_image.json", bindApi(api.LargeLanguageModels.RecognizeReceiptImageHandler, config))
				}
			}

			// Exchange Rates
			apiV1Route.GET("/exchange_rates/latest.json", bindApi(api.ExchangeRates.LatestExchangeRateHandler, config))
			apiV1Route.POST("/exchange_rates/user_custom/update.json", bindApi(api.ExchangeRates.UserCustomExchangeRateUpdateHandler, config))
			apiV1Route.POST("/exchange_rates/user_custom/delete.json", bindApi(api.ExchangeRates.UserCustomExchangeRateDeleteHandler, config))

			// System
			apiV1Route.GET("/systems/version.json", bindApi(api.Systems.VersionHandler, config))
		}
	}

	listenAddr := fmt.Sprintf("%s:%d", config.HttpAddr, config.HttpPort)

	if config.Protocol == settings.SCHEME_SOCKET {
		log.BootInfof(c, "[webserver.startWebServer] will run at socks:%s", config.UnixSocketPath)
		err = router.RunUnix(config.UnixSocketPath)
	} else if config.Protocol == settings.SCHEME_HTTP {
		log.BootInfof(c, "[webserver.startWebServer] will run at http://%s", listenAddr)
		err = router.Run(listenAddr)
	} else if config.Protocol == settings.SCHEME_HTTPS {
		log.BootInfof(c, "[webserver.startWebServer] will run at https://%s", listenAddr)
		err = router.RunTLS(listenAddr, config.CertFile, config.CertKeyFile)
	} else {
		err = errs.ErrInvalidProtocol
	}

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] cannot start, because %s", err)
		return err
	}

	return nil
}

func bindMiddleware(fn core.MiddlewareHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(core.WrapWebContext(c, config.TrustedProxyIPs))
	}
}

func bindLocalFile(filePath string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.File(filePath)
	}
}

func bindRedirect(fn core.RedirectHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		url, err := fn(c)

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
		} else {
			c.Redirect(http.StatusFound, url)
		}
	}
}

func bindApi(fn core.ApiHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, err := fn(c)

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
		} else {
			utils.PrintJsonSuccessResult(c, result)
		}
	}
}

func bindApiWithTokenUpdate(fn core.ApiHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, err := fn(c)

		if err == nil && config.MapProvider == settings.AmapProvider && config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationInternalProxyMethod {
			middlewares.AmapApiProxyAuthCookie(c, config)
		}

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
		} else {
			utils.PrintJsonSuccessResult(c, result)
		}
	}
}

func bindJSONRPCApi(fns map[string]core.JSONRPCApiHandlerFunc, skipMethods map[string]int, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)

		var jsonRPCRequest core.JSONRPCRequest
		reqErr := c.ShouldBindBodyWithJSON(&jsonRPCRequest)

		if reqErr != nil {
			utils.PrintJSONRPCErrorResult(c, nil, errs.NewIncompleteOrIncorrectSubmissionError(reqErr))
			return
		}

		if skipMethods != nil {
			httpStatusCode, exists := skipMethods[jsonRPCRequest.Method]

			if exists {
				c.AbortWithStatus(httpStatusCode)
				return
			}
		}

		fn, exists := fns[jsonRPCRequest.Method]

		if !exists {
			utils.PrintJSONRPCErrorResult(c, &jsonRPCRequest, errs.ErrApiNotFound)
			return
		}

		result, err := fn(c, &jsonRPCRequest)

		if err != nil {
			utils.PrintJSONRPCErrorResult(c, &jsonRPCRequest, err)
		} else {
			utils.PrintJSONRPCSuccessResult(c, &jsonRPCRequest, result)
		}
	}
}

func bindEventStreamApi(fn core.EventStreamApiHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		utils.SetEventStreamHeader(c)
		err := fn(c)

		if err != nil {
			utils.WriteEventStreamJsonErrorResult(c, err)
		}
	}
}

func bindCachedJs(fn core.DataHandlerFunc, config *settings.Config, store persistence.CacheStore) gin.HandlerFunc {
	return cache.CachePage(store, time.Minute, func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, _, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/javascript", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/javascript; charset=utf-8", "", result)
		}
	})
}

func bindCsv(fn core.DataHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, fileName, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/csv; charset=utf-8", fileName, result)
		}
	}
}

func bindTsv(fn core.DataHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, fileName, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/tab-separated-values; charset=utf-8", fileName, result)
		}
	}
}

func bindImage(fn core.ImageHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, contentType, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, contentType, "", result)
		}
	}
}

func bindCachedImage(fn core.ImageHandlerFunc, config *settings.Config, store persistence.CacheStore) gin.HandlerFunc {
	return cache.CachePage(store, time.Minute, func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		result, contentType, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, contentType, "", result)
		}
	})
}

func bindProxy(fn core.ProxyHandlerFunc, config *settings.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx, config.TrustedProxyIPs)
		proxy, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}
}

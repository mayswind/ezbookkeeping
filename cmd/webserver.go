package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"

	"github.com/mayswind/ezbookkeeping/pkg/api"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cron"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
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
	Subcommands: []*cli.Command{
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

	err = cron.InitializeCronJobSchedulerContainer(c, config, true)

	if err != nil {
		log.BootErrorf(c, "[webserver.startWebServer] initializes cron job scheduler failed, because %s", err.Error())
		return err
	}

	serverInfo := fmt.Sprintf("current server id is %d, current instance id is %d", requestid.Container.Current.GetCurrentServerUniqId(), requestid.Container.Current.GetCurrentInstanceUniqId())
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
	router.Use(bindMiddleware(middlewares.Recovery))

	if config.EnableGZip {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("notBlank", validators.NotBlank)
		_ = v.RegisterValidation("validUsername", validators.ValidUsername)
		_ = v.RegisterValidation("validEmail", validators.ValidEmail)
		_ = v.RegisterValidation("validCurrency", validators.ValidCurrency)
		_ = v.RegisterValidation("validHexRGBColor", validators.ValidHexRGBColor)
		_ = v.RegisterValidation("validAmountFilter", validators.ValidAmountFilter)
		_ = v.RegisterValidation("validFiscalYearStart", validators.ValidateFiscalYearStart)
	}

	router.NoRoute(bindApi(api.Default.ApiNotFound))
	router.NoMethod(bindApi(api.Default.MethodNotAllowed))

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
	router.GET("/server_settings.js", bindCachedJs(api.ServerSettings.ServerSettingsJavascriptHandler, serverSettingsCacheStore))

	for i := 0; i < len(workboxFileNames); i++ {
		router.StaticFile("/"+workboxFileNames[i], filepath.Join(config.StaticRootPath, workboxFileNames[i]))
	}

	router.StaticFile("/mobile", filepath.Join(config.StaticRootPath, "mobile.html"))
	router.Static("/mobile/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/mobile/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/mobile/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/mobile/fonts", filepath.Join(config.StaticRootPath, "fonts"))
	router.StaticFile("/mobile/favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("/mobile/favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("/mobile/touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("/mobile/manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))
	router.StaticFile("/mobile/sw.js", filepath.Join(config.StaticRootPath, "sw.js"))
	router.GET("/mobile/server_settings.js", bindCachedJs(api.ServerSettings.ServerSettingsJavascriptHandler, serverSettingsCacheStore))

	for i := 0; i < len(workboxFileNames); i++ {
		router.StaticFile("/mobile/"+workboxFileNames[i], filepath.Join(config.StaticRootPath, workboxFileNames[i]))
	}

	router.StaticFile("/desktop", filepath.Join(config.StaticRootPath, "desktop.html"))
	router.Static("/desktop/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/desktop/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/desktop/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/desktop/fonts", filepath.Join(config.StaticRootPath, "fonts"))
	router.StaticFile("/desktop/favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("/desktop/favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("/desktop/touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("/desktop/manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))
	router.StaticFile("/desktop/sw.js", filepath.Join(config.StaticRootPath, "sw.js"))
	router.GET("/desktop/server_settings.js", bindCachedJs(api.ServerSettings.ServerSettingsJavascriptHandler, serverSettingsCacheStore))

	for i := 0; i < len(workboxFileNames); i++ {
		router.StaticFile("/desktop/"+workboxFileNames[i], filepath.Join(config.StaticRootPath, workboxFileNames[i]))
	}

	if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
		avatarRoute := router.Group("/avatar")
		avatarRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString))
		{
			avatarRoute.GET("/:fileName", bindImage(api.Users.UserGetAvatarHandler))
		}
	}

	if config.EnableTransactionPictures {
		pictureRoute := router.Group("/pictures")
		pictureRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString))
		{
			pictureRoute.GET("/:fileName", bindImage(api.TransactionPictures.TransactionPictureGetHandler))
		}
	}

	router.GET("/healthz.json", bindApi(api.Healths.HealthStatusHandler))

	proxyRoute := router.Group("/proxy")
	proxyRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString))
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
				proxyRoute.GET("/map/tile/:zoomLevel/:coordinateX/:fileName", bindProxy(api.MapImages.MapTileImageProxyHandler))
			}

			if config.MapProvider == settings.TianDiTuProvider ||
				(config.MapProvider == settings.CustomProvider && config.CustomMapTileServerAnnotationLayerUrl != "") {
				proxyRoute.GET("/map/annotation/:zoomLevel/:coordinateX/:fileName", bindProxy(api.MapImages.MapAnnotationImageProxyHandler))
			}
		}
	}

	if config.MapProvider == settings.AmapProvider && config.AmapSecurityVerificationMethod == settings.AmapSecurityVerificationInternalProxyMethod {
		amapApiProxyRoute := router.Group("/_AMapService")
		amapApiProxyRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByCookie))
		{
			amapApiProxyRoute.GET("/*action", bindProxy(api.AmapApis.AmapApiProxyHandler))
		}
	}

	qrCodeRoute := router.Group("/qrcode")
	qrCodeRoute.Use(bindMiddleware(middlewares.RequestId(config)))
	{
		qrCodeCacheStore := persistence.NewInMemoryStore(time.Minute)
		qrCodeRoute.GET("/mobile_url.png", bindCachedImage(api.QrCodes.MobileUrlQrCodeHandler, qrCodeCacheStore))
	}

	apiRoute := router.Group("/api")

	apiRoute.Use(bindMiddleware(middlewares.RequestId(config)))
	apiRoute.Use(bindMiddleware(middlewares.RequestLog))
	{
		apiRoute.POST("/authorize.json", bindApiWithTokenUpdate(api.Authorizations.AuthorizeHandler, config))

		if config.EnableTwoFactor {
			twoFactorRoute := apiRoute.Group("/2fa")
			twoFactorRoute.Use(bindMiddleware(middlewares.JWTTwoFactorAuthorization))
			{
				twoFactorRoute.POST("/authorize.json", bindApiWithTokenUpdate(api.Authorizations.TwoFactorAuthorizeHandler, config))
				twoFactorRoute.POST("/recovery.json", bindApiWithTokenUpdate(api.Authorizations.TwoFactorAuthorizeByRecoveryCodeHandler, config))
			}
		}

		if config.EnableUserRegister {
			apiRoute.POST("/register.json", bindApiWithTokenUpdate(api.Users.UserRegisterHandler, config))
		}

		if config.EnableUserVerifyEmail {
			apiRoute.POST("/verify_email/resend.json", bindApi(api.Users.UserSendVerifyEmailByUnloginUserHandler))

			emailVerifyRoute := apiRoute.Group("/verify_email")
			emailVerifyRoute.Use(bindMiddleware(middlewares.JWTEmailVerifyAuthorization))
			{
				emailVerifyRoute.POST("/by_token.json", bindApi(api.Users.UserEmailVerifyHandler))
			}
		}

		if config.EnableUserForgetPassword {
			apiRoute.POST("/forget_password/request.json", bindApi(api.ForgetPasswords.UserForgetPasswordRequestHandler))

			resetPasswordRoute := apiRoute.Group("/forget_password/reset")
			resetPasswordRoute.Use(bindMiddleware(middlewares.JWTResetPasswordAuthorization))
			{
				resetPasswordRoute.POST("/by_token.json", bindApi(api.ForgetPasswords.UserResetPasswordHandler))
			}
		}

		apiRoute.GET("/logout.json", bindApiWithTokenUpdate(api.Tokens.TokenRevokeCurrentHandler, config))

		apiV1Route := apiRoute.Group("/v1")
		apiV1Route.Use(bindMiddleware(middlewares.JWTAuthorization))
		{
			// Tokens
			apiV1Route.GET("/tokens/list.json", bindApi(api.Tokens.TokenListHandler))
			apiV1Route.POST("/tokens/revoke.json", bindApi(api.Tokens.TokenRevokeHandler))
			apiV1Route.POST("/tokens/revoke_all.json", bindApi(api.Tokens.TokenRevokeAllHandler))
			apiV1Route.POST("/tokens/refresh.json", bindApiWithTokenUpdate(api.Tokens.TokenRefreshHandler, config))

			// Users
			apiV1Route.GET("/users/profile/get.json", bindApi(api.Users.UserProfileHandler))
			apiV1Route.POST("/users/profile/update.json", bindApiWithTokenUpdate(api.Users.UserUpdateProfileHandler, config))

			if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
				apiV1Route.POST("/users/avatar/update.json", bindApi(api.Users.UserUpdateAvatarHandler))
				apiV1Route.POST("/users/avatar/remove.json", bindApi(api.Users.UserRemoveAvatarHandler))
			}

			if config.EnableUserVerifyEmail {
				apiV1Route.POST("/users/verify_email/resend.json", bindApi(api.Users.UserSendVerifyEmailByLoginedUserHandler))
			}

			// Two-Factor Authorization
			if config.EnableTwoFactor {
				apiV1Route.GET("/users/2fa/status.json", bindApi(api.TwoFactorAuthorizations.TwoFactorStatusHandler))
				apiV1Route.POST("/users/2fa/enable/request.json", bindApi(api.TwoFactorAuthorizations.TwoFactorEnableRequestHandler))
				apiV1Route.POST("/users/2fa/enable/confirm.json", bindApiWithTokenUpdate(api.TwoFactorAuthorizations.TwoFactorEnableConfirmHandler, config))
				apiV1Route.POST("/users/2fa/disable.json", bindApi(api.TwoFactorAuthorizations.TwoFactorDisableHandler))
				apiV1Route.POST("/users/2fa/recovery/regenerate.json", bindApi(api.TwoFactorAuthorizations.TwoFactorRecoveryCodeRegenerateHandler))
			}

			// Data
			apiV1Route.GET("/data/statistics.json", bindApi(api.DataManagements.DataStatisticsHandler))
			apiV1Route.POST("/data/clear.json", bindApi(api.DataManagements.ClearDataHandler))

			if config.EnableDataExport {
				apiV1Route.GET("/data/export.csv", bindCsv(api.DataManagements.ExportDataToEzbookkeepingCSVHandler))
				apiV1Route.GET("/data/export.tsv", bindTsv(api.DataManagements.ExportDataToEzbookkeepingTSVHandler))
			}

			// Accounts
			apiV1Route.GET("/accounts/list.json", bindApi(api.Accounts.AccountListHandler))
			apiV1Route.GET("/accounts/get.json", bindApi(api.Accounts.AccountGetHandler))
			apiV1Route.POST("/accounts/add.json", bindApi(api.Accounts.AccountCreateHandler))
			apiV1Route.POST("/accounts/modify.json", bindApi(api.Accounts.AccountModifyHandler))
			apiV1Route.POST("/accounts/hide.json", bindApi(api.Accounts.AccountHideHandler))
			apiV1Route.POST("/accounts/move.json", bindApi(api.Accounts.AccountMoveHandler))
			apiV1Route.POST("/accounts/delete.json", bindApi(api.Accounts.AccountDeleteHandler))

			// Transactions
			apiV1Route.GET("/transactions/count.json", bindApi(api.Transactions.TransactionCountHandler))
			apiV1Route.GET("/transactions/list.json", bindApi(api.Transactions.TransactionListHandler))
			apiV1Route.GET("/transactions/list/by_month.json", bindApi(api.Transactions.TransactionMonthListHandler))
			apiV1Route.GET("/transactions/statistics.json", bindApi(api.Transactions.TransactionStatisticsHandler))
			apiV1Route.GET("/transactions/statistics/trends.json", bindApi(api.Transactions.TransactionStatisticsTrendsHandler))
			apiV1Route.GET("/transactions/amounts.json", bindApi(api.Transactions.TransactionAmountsHandler))
			apiV1Route.GET("/transactions/get.json", bindApi(api.Transactions.TransactionGetHandler))
			apiV1Route.POST("/transactions/add.json", bindApi(api.Transactions.TransactionCreateHandler))
			apiV1Route.POST("/transactions/modify.json", bindApi(api.Transactions.TransactionModifyHandler))
			apiV1Route.POST("/transactions/delete.json", bindApi(api.Transactions.TransactionDeleteHandler))

			if config.EnableDataImport {
				apiV1Route.POST("/transactions/parse_dsv_file.json", bindApi(api.Transactions.TransactionParseImportDsvFileDataHandler))
				apiV1Route.POST("/transactions/parse_import.json", bindApi(api.Transactions.TransactionParseImportFileHandler))
				apiV1Route.POST("/transactions/import.json", bindApi(api.Transactions.TransactionImportHandler))
			}

			// Transaction Pictures
			if config.EnableTransactionPictures {
				apiV1Route.POST("/transaction/pictures/upload.json", bindApi(api.TransactionPictures.TransactionPictureUploadHandler))
				apiV1Route.POST("/transaction/pictures/remove_unused.json", bindApi(api.TransactionPictures.TransactionPictureRemoveUnusedHandler))
			}

			// Transaction Categories
			apiV1Route.GET("/transaction/categories/list.json", bindApi(api.TransactionCategories.CategoryListHandler))
			apiV1Route.GET("/transaction/categories/get.json", bindApi(api.TransactionCategories.CategoryGetHandler))
			apiV1Route.POST("/transaction/categories/add.json", bindApi(api.TransactionCategories.CategoryCreateHandler))
			apiV1Route.POST("/transaction/categories/add_batch.json", bindApi(api.TransactionCategories.CategoryCreateBatchHandler))
			apiV1Route.POST("/transaction/categories/modify.json", bindApi(api.TransactionCategories.CategoryModifyHandler))
			apiV1Route.POST("/transaction/categories/hide.json", bindApi(api.TransactionCategories.CategoryHideHandler))
			apiV1Route.POST("/transaction/categories/move.json", bindApi(api.TransactionCategories.CategoryMoveHandler))
			apiV1Route.POST("/transaction/categories/delete.json", bindApi(api.TransactionCategories.CategoryDeleteHandler))

			// Transaction Tags
			apiV1Route.GET("/transaction/tags/list.json", bindApi(api.TransactionTags.TagListHandler))
			apiV1Route.GET("/transaction/tags/get.json", bindApi(api.TransactionTags.TagGetHandler))
			apiV1Route.POST("/transaction/tags/add.json", bindApi(api.TransactionTags.TagCreateHandler))
			apiV1Route.POST("/transaction/tags/add_batch.json", bindApi(api.TransactionTags.TagCreateBatchHandler))
			apiV1Route.POST("/transaction/tags/modify.json", bindApi(api.TransactionTags.TagModifyHandler))
			apiV1Route.POST("/transaction/tags/hide.json", bindApi(api.TransactionTags.TagHideHandler))
			apiV1Route.POST("/transaction/tags/move.json", bindApi(api.TransactionTags.TagMoveHandler))
			apiV1Route.POST("/transaction/tags/delete.json", bindApi(api.TransactionTags.TagDeleteHandler))

			// Transaction Templates
			apiV1Route.GET("/transaction/templates/list.json", bindApi(api.TransactionTemplates.TemplateListHandler))
			apiV1Route.GET("/transaction/templates/get.json", bindApi(api.TransactionTemplates.TemplateGetHandler))
			apiV1Route.POST("/transaction/templates/add.json", bindApi(api.TransactionTemplates.TemplateCreateHandler))
			apiV1Route.POST("/transaction/templates/modify.json", bindApi(api.TransactionTemplates.TemplateModifyHandler))
			apiV1Route.POST("/transaction/templates/hide.json", bindApi(api.TransactionTemplates.TemplateHideHandler))
			apiV1Route.POST("/transaction/templates/move.json", bindApi(api.TransactionTemplates.TemplateMoveHandler))
			apiV1Route.POST("/transaction/templates/delete.json", bindApi(api.TransactionTemplates.TemplateDeleteHandler))

			// Exchange Rates
			apiV1Route.GET("/exchange_rates/latest.json", bindApi(api.ExchangeRates.LatestExchangeRateHandler))
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

func bindMiddleware(fn core.MiddlewareHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(core.WrapWebContext(c))
	}
}

func bindApi(fn core.ApiHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
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
		c := core.WrapWebContext(ginCtx)
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

func bindCachedJs(fn core.DataHandlerFunc, store persistence.CacheStore) gin.HandlerFunc {
	return cache.CachePage(store, time.Minute, func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
		result, _, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/javascript", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/javascript", "", result)
		}
	})
}

func bindCsv(fn core.DataHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
		result, fileName, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/csv", fileName, result)
		}
	}
}

func bindTsv(fn core.DataHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
		result, fileName, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/tab-separated-values", fileName, result)
		}
	}
}

func bindImage(fn core.ImageHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
		result, contentType, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, contentType, "", result)
		}
	}
}

func bindCachedImage(fn core.ImageHandlerFunc, store persistence.CacheStore) gin.HandlerFunc {
	return cache.CachePage(store, time.Minute, func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
		result, contentType, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, contentType, "", result)
		}
	})
}

func bindProxy(fn core.ProxyHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapWebContext(ginCtx)
		proxy, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}
}

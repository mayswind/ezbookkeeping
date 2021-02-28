package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mssola/user_agent"
	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/pkg/api"
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/middlewares"
	"github.com/mayswind/lab/pkg/requestid"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/utils"
	"github.com/mayswind/lab/pkg/validators"
)

// WebServer represents the server command
var WebServer = &cli.Command{
	Name:  "server",
	Usage: "lab web server operation",
	Subcommands: []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run lab web server",
			Action: startWebServer,
		},
	},
}

func startWebServer(c *cli.Context) error {
	config, err := initializeSystem(c)

	if err != nil {
		return err
	}

	log.BootInfof("[server.startWebServer] static root path is %s", config.StaticRootPath)

	if config.AutoUpdateDatabase {
		err = updateAllDatabaseTablesStructure()

		if err != nil {
			log.BootErrorf("[server.startWebServer] update database table structure failed, because %s", err.Error())
			return err
		}
	}

	err = requestid.InitializeRequestIdGenerator(config)

	if err != nil {
		log.BootErrorf("[server.startWebServer] initializes requestid generator failed, because %s", err.Error())
		return err
	}

	serverInfo := fmt.Sprintf("current server id is %d, current instance id is %d", requestid.Container.Current.GetCurrentServerUniqId(), requestid.Container.Current.GetCurrentInstanceUniqId())
	uuidServerInfo := ""
	if config.UuidGeneratorType == settings.InternalUuidGeneratorType {
		uuidServerInfo = fmt.Sprintf(", current uuid server id is %d", config.UuidServerId)
	}

	log.BootInfof("[server.startWebServer] %s%s", serverInfo, uuidServerInfo)

	if config.Mode == settings.MODE_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

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
	}

	router.NoRoute(bindApi(api.Default.ApiNotFound))
	router.NoMethod(bindApi(api.Default.MethodNotAllowed))

	router.GET("/", func(c *gin.Context) {
		ua := user_agent.New(c.GetHeader("User-Agent"))

		if ua.Mobile() {
			c.Redirect(http.StatusMovedPermanently, "/mobile/")
		} else {
			c.Redirect(http.StatusMovedPermanently, "/desktop/")
		}
	})

	router.StaticFile("robots.txt", filepath.Join(config.StaticRootPath, "robots.txt"))
	router.StaticFile("favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))

	mobileEntryRoute := router.Group("/mobile")
	mobileEntryRoute.Use(bindMiddleware(middlewares.ServerSettingsCookie(config)))
	{
		mobileEntryRoute.StaticFile("/", filepath.Join(config.StaticRootPath, "mobile.html"))
	}
	router.Static("/mobile/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/mobile/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/mobile/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/mobile/fonts", filepath.Join(config.StaticRootPath, "fonts"))
	router.Static("/mobile/sw", filepath.Join(config.StaticRootPath, "sw"))
	router.StaticFile("/mobile/favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("/mobile/favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("/mobile/touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("/mobile/manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))
	router.StaticFile("/mobile/sw.js", filepath.Join(config.StaticRootPath, "sw.js"))

	desktopEntryRoute := router.Group("/desktop")
	desktopEntryRoute.Use(bindMiddleware(middlewares.ServerSettingsCookie(config)))
	{
		desktopEntryRoute.StaticFile("/", filepath.Join(config.StaticRootPath, "desktop.html"))
	}
	router.Static("/desktop/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/desktop/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/desktop/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/desktop/fonts", filepath.Join(config.StaticRootPath, "fonts"))
	router.StaticFile("/desktop/favicon.ico", filepath.Join(config.StaticRootPath, "favicon.ico"))
	router.StaticFile("/desktop/favicon.png", filepath.Join(config.StaticRootPath, "favicon.png"))
	router.StaticFile("/desktop/touchicon.png", filepath.Join(config.StaticRootPath, "touchicon.png"))
	router.StaticFile("/desktop/manifest.json", filepath.Join(config.StaticRootPath, "manifest.json"))

	apiRoute := router.Group("/api")

	apiRoute.Use(bindMiddleware(middlewares.RequestId(config)))
	apiRoute.Use(bindMiddleware(middlewares.RequestLog))
	{
		apiRoute.POST("/authorize.json", bindApi(api.Authorizations.AuthorizeHandler))

		if config.EnableTwoFactor {
			twoFactorRoute := apiRoute.Group("/2fa")
			twoFactorRoute.Use(bindMiddleware(middlewares.JWTTwoFactorAuthorization))
			{
				twoFactorRoute.POST("/authorize.json", bindApi(api.Authorizations.TwoFactorAuthorizeHandler))
				twoFactorRoute.POST("/recovery.json", bindApi(api.Authorizations.TwoFactorAuthorizeByRecoveryCodeHandler))
			}
		}

		if config.EnableUserRegister {
			apiRoute.POST("/register.json", bindApi(api.Users.UserRegisterHandler))
		}

		if config.EnableDataExport {
			dataRoute := apiRoute.Group("/data")
			dataRoute.Use(bindMiddleware(middlewares.JWTAuthorizationByQueryString))
			{
				dataRoute.GET("/export.csv", bindCsv(api.DataManagements.ExportDataHandler))
			}
		}

		apiRoute.GET("/logout.json", bindApi(api.Tokens.TokenRevokeCurrentHandler))

		apiV1Route := apiRoute.Group("/v1")
		apiV1Route.Use(bindMiddleware(middlewares.JWTAuthorization))
		{
			// Tokens
			apiV1Route.GET("/tokens/list.json", bindApi(api.Tokens.TokenListHandler))
			apiV1Route.POST("/tokens/revoke.json", bindApi(api.Tokens.TokenRevokeHandler))
			apiV1Route.POST("/tokens/revoke_all.json", bindApi(api.Tokens.TokenRevokeAllHandler))
			apiV1Route.POST("/tokens/refresh.json", bindApi(api.Tokens.TokenRefreshHandler))

			// Users
			apiV1Route.GET("/users/profile/get.json", bindApi(api.Users.UserProfileHandler))
			apiV1Route.POST("/users/profile/update.json", bindApi(api.Users.UserUpdateProfileHandler))

			// Two Factor Authorization
			if config.EnableTwoFactor {
				apiV1Route.GET("/users/2fa/status.json", bindApi(api.TwoFactorAuthorizations.TwoFactorStatusHandler))
				apiV1Route.POST("/users/2fa/enable/request.json", bindApi(api.TwoFactorAuthorizations.TwoFactorEnableRequestHandler))
				apiV1Route.POST("/users/2fa/enable/confirm.json", bindApi(api.TwoFactorAuthorizations.TwoFactorEnableConfirmHandler))
				apiV1Route.POST("/users/2fa/disable.json", bindApi(api.TwoFactorAuthorizations.TwoFactorDisableHandler))
				apiV1Route.POST("/users/2fa/recovery/regenerate.json", bindApi(api.TwoFactorAuthorizations.TwoFactorRecoveryCodeRegenerateHandler))
			}

			// Data
			apiV1Route.POST("/data/clear.json", bindApi(api.DataManagements.ClearDataHandler))

			// Overview
			apiV1Route.GET("/overviews/transaction.json", bindApi(api.Overviews.TransactionOverviewHandler))

			// Accounts
			apiV1Route.GET("/accounts/list.json", bindApi(api.Accounts.AccountListHandler))
			apiV1Route.GET("/accounts/get.json", bindApi(api.Accounts.AccountGetHandler))
			apiV1Route.POST("/accounts/add.json", bindApi(api.Accounts.AccountCreateHandler))
			apiV1Route.POST("/accounts/modify.json", bindApi(api.Accounts.AccountModifyHandler))
			apiV1Route.POST("/accounts/hide.json", bindApi(api.Accounts.AccountHideHandler))
			apiV1Route.POST("/accounts/move.json", bindApi(api.Accounts.AccountMoveHandler))
			apiV1Route.POST("/accounts/delete.json", bindApi(api.Accounts.AccountDeleteHandler))

			// Transactions
			apiV1Route.GET("/transactions/list.json", bindApi(api.Transactions.TransactionListHandler))
			apiV1Route.GET("/transactions/list/by_month.json", bindApi(api.Transactions.TransactionListHandler))
			apiV1Route.GET("/transactions/get.json", bindApi(api.Transactions.TransactionGetHandler))
			apiV1Route.POST("/transactions/add.json", bindApi(api.Transactions.TransactionCreateHandler))
			apiV1Route.POST("/transactions/modify.json", bindApi(api.Transactions.TransactionModifyHandler))
			apiV1Route.POST("/transactions/delete.json", bindApi(api.Transactions.TransactionDeleteHandler))

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
			apiV1Route.POST("/transaction/tags/modify.json", bindApi(api.TransactionTags.TagModifyHandler))
			apiV1Route.POST("/transaction/tags/hide.json", bindApi(api.TransactionTags.TagHideHandler))
			apiV1Route.POST("/transaction/tags/move.json", bindApi(api.TransactionTags.TagMoveHandler))
			apiV1Route.POST("/transaction/tags/delete.json", bindApi(api.TransactionTags.TagDeleteHandler))

			// Statistics
			apiV1Route.GET("/statistics/transaction.json", bindApi(api.Statistics.TransactionStatisticsHandler))

			// Exchange Rates
			apiV1Route.GET("/exchange_rates/latest.json", bindApi(api.ExchangeRates.LatestExchangeRateHandler))
		}
	}

	listenAddr := fmt.Sprintf("%s:%d", config.HttpAddr, config.HttpPort)

	if config.Protocol == settings.SCHEME_SOCKET {
		log.BootInfof("[server.startWebServer] will run at socks:%s", config.UnixSocketPath)
		err = router.RunUnix(config.UnixSocketPath)
	} else if config.Protocol == settings.SCHEME_HTTP {
		log.BootInfof("[server.startWebServer] will run at http://%s", listenAddr)
		err = router.Run(listenAddr)
	} else if config.Protocol == settings.SCHEME_HTTPS {
		log.BootInfof("[server.startWebServer] will run at https://%s", listenAddr)
		err = router.RunTLS(listenAddr, config.CertFile, config.CertKeyFile)
	} else {
		err = errs.ErrInvalidProtocol
	}

	if err != nil {
		log.BootErrorf("[server.startWebServer] cannot start, because %s", err)
		return err
	}

	return nil
}

func bindMiddleware(fn core.MiddlewareHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(core.WrapContext(c))
	}
}

func bindApi(fn core.ApiHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapContext(ginCtx)
		result, err := fn(c)

		if err != nil {
			utils.PrintJsonErrorResult(c, err)
		} else {
			utils.PrintJsonSuccessResult(c, result)
		}
	}
}

func bindCsv(fn core.DataHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := core.WrapContext(ginCtx)
		result, fileName, err := fn(c)

		if err != nil {
			utils.PrintDataErrorResult(c, "text/text", err)
		} else {
			utils.PrintDataSuccessResult(c, "text/csv", fileName, result)
		}
	}
}

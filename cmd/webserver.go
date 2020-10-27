package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli"

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

var WebServer = cli.Command{
	Name:  "server",
	Usage: "lab web server operation",
	Subcommands: []cli.Command{
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

	err = requestid.InitializeRequestIdGenerator(config)

	if err != nil {
		log.BootErrorf("[server.startWebServer] initializes requestid generator failed, because %s", err.Error())
		return err
	}

	serverInfo := fmt.Sprintf("current server id is %d, current instance id is %d", requestid.Container.Current.GetCurrentServerUniqId(), requestid.Container.Current.GetCurrentInstanceUniqId())
	uuidServerInfo := ""
	if config.UuidGeneratorType == settings.UUID_GENERATOR_TYPE_INTERNAL {
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
	}

	router.NoRoute(bindApi(api.Default.ApiNotFound))
	router.NoMethod(bindApi(api.Default.MethodNotAllowed))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/mobile/")
	})

	router.StaticFile("robots.txt", filepath.Join(config.StaticRootPath, "robots.txt"))
	router.Static("/js", filepath.Join(config.StaticRootPath, "js"))
	router.Static("/css", filepath.Join(config.StaticRootPath, "css"))
	router.Static("/img", filepath.Join(config.StaticRootPath, "img"))
	router.Static("/fonts", filepath.Join(config.StaticRootPath, "fonts"))

	mobileEntryRoute := router.Group("/mobile")
	mobileEntryRoute.Use(bindMiddleware(middlewares.ServerSettingsCookie(config)))
	{
		mobileEntryRoute.StaticFile("/", filepath.Join(config.StaticRootPath, "mobile.html"))
	}

	desktopEntryRoute := router.Group("/desktop")
	desktopEntryRoute.Use(bindMiddleware(middlewares.ServerSettingsCookie(config)))
	{
		desktopEntryRoute.StaticFile("/", filepath.Join(config.StaticRootPath, "desktop.html"))
	}

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

		apiRoute.GET("/logout.json", bindApi(api.Tokens.TokenRevokeCurrentHandler))

		apiV1Route := apiRoute.Group("/v1")
		apiV1Route.Use(bindMiddleware(middlewares.JWTAuthorization))
		{
			// Tokens
			apiV1Route.GET("/tokens/list.json", bindApi(api.Tokens.TokenListHandler))
			apiV1Route.POST("/tokens/revoke.json", bindApi(api.Tokens.TokenRevokeHandler))
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
	} else  {
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
			utils.PrintErrorResult(c, err)
		} else {
			utils.PrintSuccessResult(c, result)
		}
	}
}

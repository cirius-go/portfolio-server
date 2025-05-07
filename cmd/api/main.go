// @title						portfolio-server API
// @version					1.0
// @termsOfService				N/A
// @schemes					http https
// @contact.name				API Support
// @contact.email				hieu.tran21198@gmail.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:4000
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @basePath					/
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	echoswag "github.com/swaggo/echo-swagger"

	"github.com/cirius-go/portfolio-server/docs/swagger"
	_ "github.com/cirius-go/portfolio-server/docs/swagger"
	"github.com/cirius-go/portfolio-server/internal/api/apicms"
	"github.com/cirius-go/portfolio-server/internal/config"
	"github.com/cirius-go/portfolio-server/internal/service/servicecms"
	"github.com/cirius-go/portfolio-server/internal/uow"
	"github.com/cirius-go/portfolio-server/pkg/db"
	"github.com/cirius-go/portfolio-server/pkg/errors"
	"github.com/cirius-go/portfolio-server/pkg/server"
)

var (
	cfgFile    = flag.String("cfg", ".env", "the path to the config file")
	exampleCfg = flag.Bool("example", false, "print the example config")
)

func main() {
	flag.Parse()

	// load configuration
	cfg, err := config.Load(*cfgFile)
	panicIf(err)

	// connect to the database
	pg, err := db.NewPostgres(cfg.PGDB)
	panicIf(err)
	defer pg.Conn.Close()

	// create unit of work
	unitOfWork := uow.New(pg.DB)

	// init rbac enforcer
	// enf := casbin.NewEnforcer(
	// 	util.NewRBACModel(),
	// 	config.IsLocal(), // debug if local
	// )

	// aws config
	// awsCfg, err := awscfg.LoadDefaultConfig(context.Background())
	// panicIf(err)
	//
	// if config.IsLocal() {
	// 	awsCfg, err = config.GetAWSExecutionLocalConfig()
	// 	panicIf(err)
	// }

	// create services
	var (
		//+codegen=DefineCmsServices
		userSvc    = servicecms.NewUser(unitOfWork)
		projectSvc = servicecms.NewProject(unitOfWork)
		articleSvc = servicecms.NewArticle(unitOfWork)
	)

	// new http server with config
	srvCfg := server.C().
		SetAddress(cfg.HTTPServer.Host, cfg.HTTPServer.Port).
		SetCustomErrorHandler(errors.CreateEchoErrorHandler())
	srv := server.NewHTTPWithConfig(srvCfg)
	router := srv.Echo
	router.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"Status": "OK",
		})
	})

	{
		swagger.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
		router.GET("/swagger/*", echoswag.WrapHandler)
	}
	// bind services to the http server
	cmsRouter := router.Group("/cms")
	for _, registrar := range []HTTPRegistrar{
		//+codegen=DefineCmsAPIs
		apicms.NewUser(userSvc),
		apicms.NewProject(projectSvc),
		apicms.NewArticle(articleSvc),
	} {
		registrar.RegisterHTTP(cmsRouter)
	}

	if config.IsInAWSLambda() {
		startLambda(router)
	} else {
		err = server.StartHTTP(srv)
		panicIf(err)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// HTTPRegistrar is a handler that registers HTTP handlers.
type HTTPRegistrar interface {
	RegisterHTTP(g *echo.Group)
}

func startLambda(router *echo.Echo) {
	router.HideBanner = true
	echoLambda := echoadapter.NewV2(router)
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return echoLambda.ProxyWithContext(ctx, req)
	})
}

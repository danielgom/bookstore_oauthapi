package app

import (
	"github.com/danielgom/bookstore_oauthapi/src/datasource/clients/cassandra"
	"github.com/danielgom/bookstore_oauthapi/src/http"
	"github.com/danielgom/bookstore_oauthapi/src/repository/db"
	"github.com/danielgom/bookstore_oauthapi/src/repository/usersdb"
	"github.com/danielgom/bookstore_oauthapi/src/services/accesstoken"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	router    = echo.New()
	atHandler = http.NewHandler(accesstoken.NewService(db.NewRepository(), usersdb.NewRepository()))
)

func StartApplication() {

	cassandra.Init()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          nil,
		Format:           "[ECHO] ${time_rfc3339} | ${status} |   ${latency_human} | ${method}  \"${uri}\" ${protocol}\n",
		CustomTimeFormat: "",
		Output:           nil,
	}))

	mapUrls()

	// Run with https or http 2
	//router.Logger.Fatal(router.StartTLS(":8080",
	//	"/Users/danielg/cert.pem", "/Users/danielg/key.pem"))

	// Normal run
	router.Logger.Fatal(router.Start(":8080"))

}

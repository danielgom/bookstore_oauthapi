package http

import (
	atDomain "github.com/danielgom/bookstore_oauthapi/src/domain/accesstoken"
	"github.com/danielgom/bookstore_oauthapi/src/services/accesstoken"
	"github.com/danielgom/bookstore_utils-go/errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewHandler(service accesstoken.Service) AccessTokenHandler {
	return &accessTokenHandler{service}
}

type AccessTokenHandler interface {
	GetById(echo.Context) error
	Create(echo.Context) error
	Health(echo.Context) error
}

type accessTokenHandler struct {
	service accesstoken.Service
}

func (handler *accessTokenHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"Status": "Ready to rumble"})
}

func (handler *accessTokenHandler) GetById(c echo.Context) error {

	aT, err := handler.service.GetByID(c.Param("atId"))
	if err != nil {
		return echo.NewHTTPError(err.Status(), err)
	}

	return c.JSON(http.StatusOK, aT)
}

func (handler *accessTokenHandler) Create(c echo.Context) error {
	request := new(atDomain.AtRequest)

	if err := c.Bind(request); err != nil {
		restErr := errors.NewBadRequestError("Invalid json body")
		return echo.NewHTTPError(restErr.Status(), restErr)
	}

	at, err := handler.service.Create(request)
	if err != nil {
		return echo.NewHTTPError(err.Status(), err)
	}

	return c.JSON(http.StatusCreated, at)
}

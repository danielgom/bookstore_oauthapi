package db

import (
	"fmt"
	"github.com/danielgom/bookstore_oauthapi/src/domain/accesstoken"
	"github.com/danielgom/bookstore_utils-go/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = `SELECT accesstoken, clientid, expires, userid FROM access_tokens WHERE accesstoken=?;`
	queryCreateAccessToken = `INSERT INTO access_tokens(accesstoken, clientid, expires, userid) VALUES (?, ?, ?, ?);`
	queryUpdateExpires     = `UPDATE access_tokens SET expires=? WHERE accesstoken=?;`
)

type CQLSession interface {
	Query(string, ...interface{}) *gocql.Query
}

var Session CQLSession

func NewRepository() DRepository {
	return &repository{}
}

type DRepository interface {
	GetByID(string) (*accesstoken.AccessToken, errors.RestErr)
	Create(*accesstoken.AccessToken) errors.RestErr
	UpdateExpirationTime(*accesstoken.AccessToken) errors.RestErr
}

type repository struct {
}

func (r *repository) GetByID(id string) (*accesstoken.AccessToken, errors.RestErr) {

	tk := new(accesstoken.AccessToken)
	if err := Session.Query(queryGetAccessToken, id).Scan(&tk.AccessToken, &tk.ClientId, &tk.Expires, &tk.UserId); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("No access token found with given id")
		}
		return nil, errors.NewInternalServerError(fmt.Sprintf("error retrieving user with id %s", id), err)
	}

	return tk, nil
}

func (r *repository) Create(at *accesstoken.AccessToken) errors.RestErr {

	if err := Session.Query(queryCreateAccessToken, at.AccessToken, at.ClientId, at.Expires, at.UserId).Exec(); err != nil {
		return errors.NewInternalServerError(" error creating access token", err)
	}

	return nil
}

func (r *repository) UpdateExpirationTime(at *accesstoken.AccessToken) errors.RestErr {

	if err := Session.Query(queryUpdateExpires, at.Expires, at.AccessToken).Exec(); err != nil {
		return errors.NewInternalServerError("error updating access token", err)
	}

	return nil
}

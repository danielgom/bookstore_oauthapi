package db

import (
	"github.com/danielgom/bookstore_oauthapi/src/datasource/clients/cassandra"
	"github.com/danielgom/bookstore_oauthapi/src/domain/accesstoken"
	"github.com/danielgom/bookstore_utils-go/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = `SELECT accesstoken, clientid, expires, userid FROM access_tokens WHERE accesstoken=?;`
	queryCreateAccessToken = `INSERT INTO access_tokens(accesstoken, clientid, expires, userid) VALUES (?, ?, ?, ?);`
	queryUpdateExpires     = `UPDATE access_tokens SET expires=? WHERE accesstoken=?;`
)

func NewRepository() DRepository {
	return &repository{}
}

type DRepository interface {
	GetByID(string) (*accesstoken.AccessToken, *errors.RestErr)
	Create(*accesstoken.AccessToken) *errors.RestErr
	UpdateExpirationTime(*accesstoken.AccessToken) *errors.RestErr
}

type repository struct {
}

func (r *repository) GetByID(id string) (*accesstoken.AccessToken, *errors.RestErr) {

	tk := new(accesstoken.AccessToken)
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(&tk.AccessToken, &tk.ClientId, &tk.Expires, &tk.UserId); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("No access token found with given id")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}

	return tk, nil
}

func (r *repository) Create(at *accesstoken.AccessToken) *errors.RestErr {

	if err := cassandra.GetSession().Query(queryCreateAccessToken, at.AccessToken, at.ClientId, at.Expires, at.UserId).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *repository) UpdateExpirationTime(at *accesstoken.AccessToken) *errors.RestErr {

	if err := cassandra.GetSession().Query(queryUpdateExpires, at.Expires, at.AccessToken).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

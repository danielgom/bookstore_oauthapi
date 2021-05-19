package accesstoken

import (
	"github.com/danielgom/bookstore_oauthapi/src/domain/accesstoken"
	"github.com/danielgom/bookstore_oauthapi/src/repository/db"
	"github.com/danielgom/bookstore_oauthapi/src/repository/usersdb"
	"github.com/danielgom/bookstore_utils-go/errors"
	"strings"
)

func NewService(dbRepo db.DRepository, usersRepo usersdb.UsersRepository) Service {
	return &service{dbRepo, usersRepo}
}

type Service interface {
	GetByID(string) (*accesstoken.AccessToken, errors.RestErr)
	Create(request *accesstoken.AtRequest) (*accesstoken.AccessToken, errors.RestErr)
	UpdateExpirationTime(*accesstoken.AccessToken) errors.RestErr
}

type service struct {
	DbRepository    db.DRepository
	usersRepository usersdb.UsersRepository
}

func (s *service) GetByID(id string) (*accesstoken.AccessToken, errors.RestErr) {

	atId := strings.TrimSpace(id)

	if len(atId) == 0 {
		return nil, errors.NewBadRequestError("Invalid access token id")
	}

	at, err := s.DbRepository.GetByID(atId)
	if err != nil {
		return nil, err
	}

	return at, nil
}

func (s *service) Create(request *accesstoken.AtRequest) (*accesstoken.AccessToken, errors.RestErr) {

	if err := request.Validate(); err != nil {
		return nil, err
	}

	//TODO: support both grant types

	user, err := s.usersRepository.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	at := accesstoken.GetNewAccessToken(user.Id)

	if err = s.DbRepository.Create(at); err != nil {
		return nil, err
	}

	return at, nil
}

func (s *service) UpdateExpirationTime(at *accesstoken.AccessToken) errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.DbRepository.UpdateExpirationTime(at)
}

package accesstoken

import (
	errors2 "errors"
	"github.com/danielgom/bookstore_oauthapi/src/domain/accesstoken"
	"github.com/danielgom/bookstore_oauthapi/src/domain/users"
	"github.com/danielgom/bookstore_utils-go/errors"
	"testing"
)

type MockGetById func(id string) (*accesstoken.AccessToken, errors.RestErr)
type MockCreate func(at *accesstoken.AccessToken) errors.RestErr
type MockUpdateExpirationTime func(at *accesstoken.AccessToken) errors.RestErr

type MockLoginUser func(email, password string) (*users.User, errors.RestErr)

type MockDbRepository struct {
	MockGetById              MockGetById
	MockCreate               MockCreate
	MockUpdateExpirationTime MockUpdateExpirationTime
}

type MockUsersRepository struct {
	MockLoginUser MockLoginUser
}

func (m *MockDbRepository) GetByID(id string) (*accesstoken.AccessToken, errors.RestErr) {
	return m.MockGetById(id)
}

func (m *MockDbRepository) Create(token *accesstoken.AccessToken) errors.RestErr {
	return m.MockCreate(token)
}

func (m *MockDbRepository) UpdateExpirationTime(token *accesstoken.AccessToken) errors.RestErr {
	return m.MockUpdateExpirationTime(token)
}

func (m *MockUsersRepository) LoginUser(email string, password string) (*users.User, errors.RestErr) {
	return m.MockLoginUser(email, password)
}

func TestNewService(t *testing.T) {
	dbRepository := &MockDbRepository{}
	userRepository := &MockUsersRepository{}
	newService := NewService(dbRepository, userRepository)

	if newService == nil {
		t.Error("service should not be nil")
	}
}

func TestServiceGetByID(t *testing.T) {

	t.Run("Should throw error when id is 0", func(t *testing.T) {

		dbRepository := &MockDbRepository{}
		MockService := service{
			DbRepository: dbRepository,
		}
		atString := ""
		aT, err := MockService.GetByID(atString)

		if aT != nil {
			t.Error("access token should be nil")
		}

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should return db error", func(t *testing.T) {

		dbRepository := &MockDbRepository{MockGetById: func(id string) (*accesstoken.AccessToken, errors.RestErr) {
			return nil, errors.NewNotFoundError("No access token found with given id")
		}}
		MockService := service{
			DbRepository: dbRepository,
		}
		atString := "22"
		aT, err := MockService.GetByID(atString)

		if aT != nil {
			t.Error("access token should be nil")
		}

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should return accessToken", func(t *testing.T) {

		dbRepository := &MockDbRepository{MockGetById: func(id string) (*accesstoken.AccessToken, errors.RestErr) {
			return &accesstoken.AccessToken{
				AccessToken: "123456",
				UserId:      123,
				ClientId:    456,
				Expires:     365,
			}, nil
		}}
		MockService := service{
			DbRepository: dbRepository,
		}
		atString := "123456"
		aT, err := MockService.GetByID(atString)

		if aT == nil {
			t.Error("access token should be nil")
		}

		if aT != nil && aT.AccessToken != "123456" {
			t.Errorf("access token should be %s but %s received", "123456", aT.AccessToken)
		}

		if aT != nil && aT.Expires != 365 {
			t.Errorf("Expires should be %s but %d received", "123456", aT.Expires)
		}

		if aT != nil && aT.UserId != 123 {
			t.Errorf("UserId should be %s but %d received", "123456", aT.UserId)
		}

		if aT != nil && aT.ClientId != 456 {
			t.Errorf("ClientId should be %s but %d received", "123456", aT.ClientId)
		}

		if err != nil {
			t.Error("error should not be nil")
		}
	})
}

func TestServiceCreate(t *testing.T) {

	t.Run("Should return error on validation", func(t *testing.T) {

		MockService := service{}

		atRequest := &accesstoken.AtRequest{
			GrantType: "none",
		}

		aT, err := MockService.Create(atRequest)

		if aT != nil {
			t.Error("access token should be nil")
		}

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should return error on login user", func(t *testing.T) {

		userRepository := &MockUsersRepository{
			MockLoginUser: func(email, password string) (*users.User, errors.RestErr) {
				return nil, errors.NewInternalServerError("Invalid response from user API while trying to login",
					errors2.New("no response"))
			}}

		MockService := service{
			usersRepository: userRepository,
		}

		atRequest := &accesstoken.AtRequest{
			GrantType: "password",
		}

		aT, err := MockService.Create(atRequest)

		if aT != nil {
			t.Error("access token should be nil")
		}

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should return error on create", func(t *testing.T) {

		userRepository := &MockUsersRepository{
			MockLoginUser: func(email, password string) (*users.User, errors.RestErr) {
				return &users.User{
					Id:        123456,
					FirstName: "daniel",
					LastName:  "gomez",
					Email:     "test@hotmail.com",
				}, nil
			}}

		dbRepository := &MockDbRepository{MockCreate: func(at *accesstoken.AccessToken) errors.RestErr {
			return errors.NewInternalServerError(" error creating access token", errors2.New("db error"))
		}}

		MockService := service{
			DbRepository:    dbRepository,
			usersRepository: userRepository,
		}

		atRequest := &accesstoken.AtRequest{
			GrantType: "password",
		}

		aT, err := MockService.Create(atRequest)

		if aT != nil {
			t.Error("access token should be nil")
		}

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should return created access token", func(t *testing.T) {

		userRepository := &MockUsersRepository{
			MockLoginUser: func(email, password string) (*users.User, errors.RestErr) {
				return &users.User{
					Id:        123456,
					FirstName: "daniel",
					LastName:  "gomez",
					Email:     "test@hotmail.com",
				}, nil
			}}

		dbRepository := &MockDbRepository{MockCreate: func(at *accesstoken.AccessToken) errors.RestErr {
			return nil
		}}

		MockService := service{
			DbRepository:    dbRepository,
			usersRepository: userRepository,
		}

		atRequest := &accesstoken.AtRequest{
			GrantType: "password",
		}

		aT, err := MockService.Create(atRequest)

		if err != nil {
			t.Error("error should be nil")
		}

		if aT == nil {
			t.Error("access token should not be nil")
		}
	})
}

func TestServiceUpdateExpirationTime(t *testing.T) {

	t.Run("Should return error on validation", func(t *testing.T) {
		MockService := service{}
		err := MockService.UpdateExpirationTime(nil)

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should return no error", func(t *testing.T) {

		dbRepository := &MockDbRepository{MockUpdateExpirationTime: func(at *accesstoken.AccessToken) errors.RestErr {
			return nil
		}}

		MockService := service{
			DbRepository: dbRepository,
		}
		err := MockService.UpdateExpirationTime(&accesstoken.AccessToken{
			AccessToken: "1234",
			UserId:      1,
			ClientId:    1,
			Expires:     365,
		})

		if err != nil {
			t.Error("error should not be nil")
		}
	})
}

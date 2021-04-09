package usersdb

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/danielgom/bookstore_oauthapi/src/domain/users"
	"github.com/danielgom/bookstore_utils-go/errors"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func NewRepository() UsersRepository {
	return &usersRepository{}
}

type UsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct {
}

func (u *usersRepository) LoginUser(email, password string) (*users.User, *errors.RestErr) {

	request := users.LoginRequest{
		Email:    email,
		Password: password,
	}

	b, _ := json.Marshal(request)
	postBody := bytes.NewBuffer(b)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8081/users/login", postBody)

	resp, err := Client.Do(r)

	if err != nil {
		return nil, errors.NewInternalServerError("Invalid response from user API while trying to login")
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
		}
	}()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode > 399 {
		restErr := new(errors.RestErr)
		err = json.Unmarshal(respBody, restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("Invalid error interface when trying to login the user")
		}
		return nil, restErr
	}

	user := new(users.User)
	if err = json.Unmarshal(respBody, user); err != nil {
		return nil, errors.NewInternalServerError("Error when trying to unmarshal user response")
	}

	return user, nil
}

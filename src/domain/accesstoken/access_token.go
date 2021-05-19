package accesstoken

import (
	"fmt"
	"github.com/danielgom/bookstore_oauthapi/src/utils/cryptoutils"
	"github.com/danielgom/bookstore_utils-go/errors"
	"strings"
	"time"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "clientCredentials"
)

type AtRequest struct {
	GrantType string `json:"grantType"`
	Scope     string `json:"scope"`

	// Used for password gran type

	Username string `json:"username"`
	Password string `json:"password"`

	// Used for clientCredentials grant type

	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret"`
}

func (request *AtRequest) Validate() errors.RestErr {
	switch request.GrantType {
	case grantTypePassword:

	case grantTypeClientCredentials:

	default:
		return errors.NewBadRequestError("Invalid grantType parameter")

	}
	// Todo: Validate parameters for each grantType
	return nil
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
	UserId      int64  `json:"userId"`
	ClientId    int64  `json:"clientId,omitempty"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() errors.RestErr {

	if at == nil {
		return errors.NewBadRequestError("Access token pointer is nil")
	}

	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("Invalid access token")
	}

	if at.UserId <= 0 {
		return errors.NewBadRequestError("Invalid user id")
	}

	if at.ClientId <= 0 {
		return errors.NewBadRequestError("Invalid client id")
	}

	if at.Expires <= 0 {
		return errors.NewBadRequestError("Invalid expiration time")
	}

	return nil
}

func GetNewAccessToken(userId int64) *AccessToken {
	at := &AccessToken{
		UserId:  userId,
		Expires: time.Now().Add(time.Hour * expirationTime).Unix(),
	}
	at.AccessToken = cryptoutils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
	return at
}

func (at *AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now())
}

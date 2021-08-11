package accesstoken

import (
	"testing"
	"time"
)

func TestAccessTokenConstants(t *testing.T) {

	if expirationTime != 24 {
		t.Error("Expiration time should be 24 hours")
	}
}

func TestAtRequestValidate(t *testing.T) {
	t.Parallel()

	t.Run("Should throw error on invalid validation", func(t *testing.T) {
		t.Parallel()
		atR := &AtRequest{
			GrantType: "none",
		}

		err := atR.Validate()

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should pass the validation with password grant_type", func(t *testing.T) {
		t.Parallel()
		atR := &AtRequest{
			GrantType: "password",
		}

		err := atR.Validate()

		if err != nil {
			t.Error("error should be nil")
		}
	})

	t.Run("Should pass the validation with credentials grant_type", func(t *testing.T) {
		t.Parallel()
		atR := &AtRequest{
			GrantType: "clientCredentials",
		}

		err := atR.Validate()

		if err != nil {
			t.Error("error should be nil")
		}
	})
}

func TestAccessTokenValidate(t *testing.T) {
	t.Parallel()

	t.Run("Should pass the access token validation", func(t *testing.T) {
		t.Parallel()
		aT := &AccessToken{
			AccessToken: "12345",
			UserId:      1,
			ClientId:    2,
			Expires:     3456,
		}
		err := aT.Validate()

		if err != nil {
			t.Error("error should be nil")
		}
	})

	t.Run("Should throw error when access token is nil", func(t *testing.T) {
		t.Parallel()
		var aT AccessToken
		err := aT.Validate()

		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("Should throw error when access token string is empty or nil", func(t *testing.T) {
		t.Parallel()
		aT := &AccessToken{
			AccessToken: "",
		}
		err := aT.Validate()

		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("Should throw error when the user id is 0", func(t *testing.T) {
		t.Parallel()
		aT := &AccessToken{
			AccessToken: "12345",
			UserId:      0,
		}
		err := aT.Validate()

		if err == nil {
			t.Error("error should not be nil")
		}
	})
	t.Run("Should throw error when expires is 0", func(t *testing.T) {
		t.Parallel()
		aT := &AccessToken{
			AccessToken: "12345",
			UserId:      1,
			ClientId:    1,
			Expires:     0,
		}
		err := aT.Validate()

		if err == nil {
			t.Error("error should not be nil")
		}
	})

	t.Run("Should throw error when the client id is 0", func(t *testing.T) {
		t.Parallel()
		aT := &AccessToken{
			AccessToken: "12345",
			UserId:      1,
			ClientId:    0,
		}
		err := aT.Validate()

		if err == nil {
			t.Error("error should not be nil")
		}
	})

}

func TestGetNewAccessToken(t *testing.T) {
	t.Parallel()
	at := GetNewAccessToken(0)
	if at.IsExpired() {
		t.Error("Brand access token should not be nil")
	}

	if at.AccessToken == "" {
		t.Error("New access token should not have defined access token id")
	}

	if at.UserId != 0 {
		t.Error("New access token should not have an associated user id")
	}
}

func TestIsExpired(t *testing.T) {
	t.Parallel()
	at := AccessToken{}
	if !at.IsExpired() {
		t.Error("Empty access token should be expired by default")
	}

	at.Expires = time.Now().Add(time.Hour * 3).Unix()

	if at.IsExpired() {
		t.Error("Access token created for 3 hours should NOT be expired")
	}
}

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

func TestGetNewAccessToken(t *testing.T) {
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
	at := AccessToken{}
	if !at.IsExpired() {
		t.Error("Empty access token should be expired by default")
	}

	at.Expires = time.Now().Add(time.Hour * 3).Unix()

	if at.IsExpired() {
		t.Error("Access token created for 3 hours should NOT be expired")
	}
}

package usersdb

import (
	"bytes"
	"encoding/json"
	"github.com/danielgom/bookstore_oauthapi/src/domain/users"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestLoginUserTimeout(t *testing.T) {

	repository := usersRepository{}
	user, restErr := repository.LoginUser("test@gmail.com", "the-password")

	expectedString := "Invalid response from user API while trying to login"

	if user != nil {
		t.Error("User should be a nil value")
	}
	if restErr == nil {
		t.Error("Error should not be a nil value")
	}
	if restErr != nil && restErr.Message() != expectedString {
		t.Errorf("\n Expected: %s, \n Received: %s", expectedString, restErr.Message())
	}
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {

	response := `{"message" : "Email not found", "status": 404,}`

	r := io.NopCloser(bytes.NewReader([]byte(response)))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil
		},
	}

	repository := usersRepository{}
	user, restErr := repository.LoginUser("test@gmail.com", "the_password")

	expectedString := "Invalid error interface when trying to login the user"

	if user != nil {
		t.Error("User should be a nil value")
	}
	if restErr == nil {
		t.Error("Error should not be a nil value")
	}
	if restErr != nil && restErr.Status() != 500 {
		t.Error("Status returned should be 500")
	}
	if restErr != nil && restErr.Message() != expectedString {
		t.Errorf("\n Expected: %s, \n Received: %s", expectedString, restErr.Message())
	}
}

func TestLoginUserInvalidCredentials(t *testing.T) {

	response := `{"message" : "Invalid email or password", "status": 400, "error": "Bad request"}`

	r := io.NopCloser(bytes.NewReader([]byte(response)))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       r,
			}, nil
		},
	}

	repository := usersRepository{}
	user, restErr := repository.LoginUser("test@gmail.com", "the_password")

	expectedString := "Invalid email or password"

	if user != nil {
		t.Error("User should be a nil value")
	}
	if restErr == nil {
		t.Error("Error should not be a nil value")
	}
	if restErr != nil && restErr.Status() != 400 {
		t.Error("Status returned should be 400")
	}
	if restErr != nil && restErr.Message() != expectedString {
		t.Errorf("\n Expected: %s, \n Received: %s", expectedString, restErr.Message())
	}
}

func TestLoginUserEmailNotFound(t *testing.T) {

	response := `{"message" : "Username with email test@gmail.com not found", "status": 404, "error": "Bad request"}`

	r := io.NopCloser(bytes.NewReader([]byte(response)))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       r,
			}, nil
		},
	}

	repository := usersRepository{}
	user, restErr := repository.LoginUser("test@gmail.com", "the_password")

	expectedString := "Username with email test@gmail.com not found"

	if user != nil {
		t.Error("User should be a nil value")
	}
	if restErr == nil {
		t.Error("Error should not be a nil value")
	}
	if restErr != nil && restErr.Status() != 404 {
		t.Error("Status returned should be 404")
	}
	if restErr != nil && restErr.Message() != expectedString {
		t.Errorf("\n Expected: %s, \n Received: %s", expectedString, restErr.Message())
	}
}

func TestLoginUserInvalidJsonResponse(t *testing.T) {

	wrongResponse := `{"id": "1", "firstName": "testing-name", "lastName": "testing-lasst",}`

	r := io.NopCloser(bytes.NewReader([]byte(wrongResponse)))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	repository := usersRepository{}
	user, restErr := repository.LoginUser("test@gmail.com", "the_password")

	expectedString := "Error when trying to unmarshal user response"

	if user != nil {
		t.Error("User should be a nil value")
	}
	if restErr == nil {
		t.Error("Error should not be a nil value")
	}
	if restErr != nil && restErr.Status() != 500 {
		t.Error("Status returned should be 500")
	}
	if restErr != nil && restErr.Message() != expectedString {
		t.Errorf("\n Expected: %s, \n Received: %s", expectedString, restErr.Message())
	}
}

func TestLoginUserNoError(t *testing.T) {

	expectedUser := &users.User{
		Id:        1,
		FirstName: "Daniel",
		LastName:  "Gomez",
		Email:     "daniel@gmail.com",
	}

	response, _ := json.Marshal(expectedUser)

	r := io.NopCloser(bytes.NewReader(response))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	repository := usersRepository{}
	actualUser, restErr := repository.LoginUser("daniel@gmail.com", "the_password")

	if actualUser == nil {
		t.Error("User should not be a nil value")
	}
	if restErr != nil {
		t.Error("Error should be a nil value")
	}
	if actualUser != nil && actualUser.Id != 1 {
		t.Errorf("\nUser first name expected: %d, \n Received: %d", expectedUser.Id, actualUser.Id)
	}
	if actualUser != nil && actualUser.FirstName != "Daniel" {
		t.Errorf("\nUser first name expected: %s, \n Received: %s", expectedUser.FirstName, actualUser.FirstName)
	}
	if actualUser != nil && actualUser.LastName != "Gomez" {
		t.Errorf("\nUser first name expected: %s, \n Received: %s", expectedUser.LastName, actualUser.LastName)
	}
	if actualUser != nil && actualUser.Email != "daniel@gmail.com" {
		t.Errorf("\nUser first name expected: %s, \n Received: %s", expectedUser.Email, actualUser.Email)
	}
}

func TestNewRepository(t *testing.T) {
	expected := usersRepository{}
	repository := NewRepository()

	if reflect.DeepEqual(expected, repository) {
		t.Error("Values should be equal")
	}
}

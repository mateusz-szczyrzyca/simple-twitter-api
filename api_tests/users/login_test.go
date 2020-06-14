package login

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"twitter/common"
)

var Endpoint string

func TestMain(m *testing.M) {
	flag.StringVar(&Endpoint, "endpoint", "http://localhost:58123", "usersEndpoint")
	flag.Parse()

	m.Run()
}

const (
	InvalidCredentials = 401
	SuccessCode        = 200
	InvalidToken       = 403
)

type TestStruct struct {
	Name                 string
	requestBody          string
	expectedStatusCode   int
	expectedResponseBody string
}

func TestUserInvalidLogins(t *testing.T) {
	invalidLoginBody := `{"message":"Invalid login or password.","token":""}`
	tests := []TestStruct{
		{"Empty request", ``, InvalidCredentials, invalidLoginBody},
		{"Empty JSON", `{}`, InvalidCredentials, invalidLoginBody},
		{"Empty username", `{"Username":""}`, InvalidCredentials, invalidLoginBody},
		{"Empty password", `{"Password":""}`, InvalidCredentials, invalidLoginBody},
		{"Invalid username", `{"Username":"_invalid_"}`, InvalidCredentials, invalidLoginBody},
		{"Invalid both", `{"Username":"_invalid_","Password" : "123456"}`, InvalidCredentials, invalidLoginBody},
	}

	for _, testCase := range tests {
		t.Log(testCase.Name)
		var reader io.Reader
		reader = strings.NewReader(testCase.requestBody)
		request, err := http.NewRequest("POST", Endpoint+"/users/login", reader)
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Error(err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		assert.Equal(t, testCase.expectedResponseBody, strings.TrimRight(string(body), "\r\n"))
	}
}

func TestUserLoginsAndLogouts(t *testing.T) {
	validLoginMessage := "Congratulations, you've provided correct credentials!"
	invalidLogoutMessage := "Invalid token."
	validLogoutMessage := "You have successfully logged out."

	tests := []TestStruct{
		{"Testing credentials admin1", `{"Username":"admin1","Password":"zecret111"}`, SuccessCode, validLoginMessage},
		{"Testing credentials admin2", `{"Username":"admin2","Password":"zecret222"}`, SuccessCode, validLoginMessage},
		{"Testing credentials admin3", `{"Username":"admin3","Password":"zecret333"}`, SuccessCode, validLoginMessage},
	}

	for _, testCase := range tests {
		t.Log(testCase.Name)
		var reader io.Reader
		var Response common.Response
		reader = strings.NewReader(testCase.requestBody) //Convert string to reader
		requestLogin, err := http.NewRequest("POST", Endpoint+"/users/login", reader)
		res, err := http.DefaultClient.Do(requestLogin)
		if err != nil {
			t.Error(err)
		}

		// Valid login first
		decoder := json.NewDecoder(res.Body)
		assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		err = decoder.Decode(&Response)
		if err != nil {
			t.Error(err)
		}

		validToken := Response.Token
		assert.Equal(t, validLoginMessage, Response.Message)

		// First invalid scenario (invalid token)
		_ = &Response
		logoutRequestBodyInvalid := fmt.Sprintf(`{"Token":"%sx"}`, validToken)
		reader = strings.NewReader(logoutRequestBodyInvalid) //Convert string to reader
		requestLogoutBad, err := http.NewRequest("POST", Endpoint+"/users/logout", reader)
		res, err = http.DefaultClient.Do(requestLogoutBad)
		if err != nil {
			t.Error(err)
		}
		decoder = json.NewDecoder(res.Body)
		assert.Equal(t, InvalidToken, res.StatusCode)
		err = decoder.Decode(&Response)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, invalidLogoutMessage, Response.Message)
		// Now success scenario (token is valid)
		_ = &Response
		logoutRequestBodyValid := fmt.Sprintf(`{"Token":"%s"}`, validToken)
		reader = strings.NewReader(logoutRequestBodyValid) //Convert string to reader
		requestLogoutGood, err := http.NewRequest("POST", Endpoint+"/users/logout", reader)
		res, err = http.DefaultClient.Do(requestLogoutGood)
		if err != nil {
			t.Error(err)
		}
		decoder = json.NewDecoder(res.Body)
		assert.Equal(t, SuccessCode, res.StatusCode)
		err = decoder.Decode(&Response)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, validLogoutMessage, Response.Message)
	}
}

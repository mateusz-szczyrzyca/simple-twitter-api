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

type TestStructGetMessages struct {
	Name                 string
	Message              string
	expectedStatusCode   int
	expectedTags         []string
	expectedResponseBody string
}

func TestGettingMessages(t *testing.T) {

	loginJSON := `{"Username":"admin2","Password":"zecret222"}`
	var reader io.Reader
	var Response common.Response
	reader = strings.NewReader(loginJSON) //Convert string to reader
	requestLogin, err := http.NewRequest("POST", Endpoint+"/users/login", reader)
	res, err := http.DefaultClient.Do(requestLogin)
	if err != nil {
		t.Error(err)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&Response)
	if err != nil {
		t.Error(err)
	}

	validToken := Response.Token

	loginJSON2 := `{"Username":"admin1","Password":"zecret111"}`
	reader = strings.NewReader(loginJSON2) //Convert string to reader
	requestLogin2, err := http.NewRequest("POST", Endpoint+"/users/login", reader)
	res, err = http.DefaultClient.Do(requestLogin2)
	if err != nil {
		t.Error(err)
	}

	// Valid login first
	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&Response)
	if err != nil {
		t.Error(err)
	}

	cannotFilterTimeFrameToken := Response.Token

	// Now requests for listing messages.
	tests := []TestStructGetMessages{
		{
			"Bad JSON.",
			fmt.Sprintf(`
{{};'
Tags["ag10
}]'`),
			422,
			nil,
			`{"message":"Invalid JSON request.","token":""}`,
		},
		//{
		//	"No tag, token and timeframe - testing limit 10",
		//	fmt.Sprintf(`{}`),
		//	200,
		//	nil,
		//	`[{"Datetime":"2015-01-26T10:10:10.555555Z","Tags":["tag100"],"Message":"Content of message 15"},{"Datetime":"2018-03-26T10:10:10.555555Z","Tags":["tag2018"],"Message":"Content of message 19"},{"Datetime":"2015-12-12T10:10:10.555555Z","Tags":["tag1,tag5"],"Message":"Content of message 14"},{"Datetime":"2015-01-25T10:10:10.555555Z","Tags":["tag1,tag2"],"Message":"Content of message 1"},{"Datetime":"2017-01-15T10:10:10.555555Z","Tags":["tag3,tag4,tag5"],"Message":"Content of message 4"},{"Datetime":"2017-07-26T10:10:10.555555Z","Tags":["tag1,tag2"],"Message":"Content of message 3"},{"Datetime":"2016-06-06T10:10:10.555555Z","Tags":["tag5"],"Message":"Content of message 8"},{"Datetime":"2018-01-01T10:10:10.555555Z","Tags":["tag2018"],"Message":"Content of message 18"},{"Datetime":"2019-03-26T10:10:10.555555Z","Tags":["tag2019"],"Message":"Content of message 9"},{"Datetime":"2016-01-26T10:10:10.555555Z","Tags":["tag2"],"Message":"Content of message 12"}]`,
		//},
		{
			"Single tag, no token",
			fmt.Sprintf(`
{
"Tags":["tag100"]
}`),
			200,
			[]string{"tag100"},
			`[{"Datetime":"2015-01-26T10:10:10.555555Z","Tags":["tag100"],"Message":"Content of message 15"}]`,
		},
		{
			"Token with tag",
			fmt.Sprintf(`
{
"Token":"%s",
"Tags":["tag100"]
}`,
				validToken),
			200,
			[]string{"tag100"},
			`[{"Datetime":"2015-01-26T10:10:10.555555Z","Tags":["tag100"],"Message":"Content of message 15"}]`,
		},
		{
			"Many tags, no token.",
			fmt.Sprintf(`
{
"Tags":["tag1","tag2","tag3"]
}`),
			200,
			[]string{"tag1", "tag2", "tag3"},
			`[{"Datetime":"2016-12-24T10:10:10.555555Z","Tags":["tag1,tag2,tag3,tag4,tag5"],"Message":"Content of message 17"}]`,
		},
		{
			"Many tags, with token.",
			fmt.Sprintf(`
{
"Token":"%s",
"Tags":["tag1","tag2","tag3"]
}`, validToken),
			200,
			[]string{"tag1", "tag2", "tag3"},
			`[{"Datetime":"2016-12-24T10:10:10.555555Z","Tags":["tag1,tag2,tag3,tag4,tag5"],"Message":"Content of message 17"}]`,
		},
		{
			"Time frame, incorrect token",
			fmt.Sprintf(`
{
	"Token":"%s_",
	"Tags":["tag5"],
	"TimeFrom":"2015-09-01",
	"TimeTo":"2016-01-28"
}`, validToken),
			403,
			nil,
			`{"message":"Your are not allowed to filter by timeline.","token":""}`,
		},
		{
			"Time frame, token with no permissions",
			fmt.Sprintf(`
{
	"Token":"%s",
	"Tags":["tag5"],
	"TimeFrom":"2015-09-01",
	"TimeTo":"2016-01-28"
}`, cannotFilterTimeFrameToken),
			403,
			nil,
			`{"message":"Your are not allowed to filter by timeline.","token":""}`,
		},
		{
			"Time frame with correct token, no tags",
			fmt.Sprintf(`
{
	"Token":"%s",
	"TimeFrom":"2015-09-01",
	"TimeTo":"2016-01-28"
}`, validToken),
			200,
			[]string{"tag1", "tag2", "tag3"},
			`[{"Datetime":"2016-01-26T10:10:10.555555Z","Tags":["tag2"],"Message":"Content of message 12"},{"Datetime":"2016-01-26T10:10:10.555555Z","Tags":["tag1,tag2"],"Message":"Content of message 2"},{"Datetime":"2015-12-12T10:10:10.555555Z","Tags":["tag1,tag5"],"Message":"Content of message 14"}]`,
		},
		{
			"Time frame with correct token, with tag",
			fmt.Sprintf(`
{
	"Token":"%s",
	"Tags":["tag5"],
	"TimeFrom":"2015-09-01",
	"TimeTo":"2016-01-28"
}`, validToken),
			200,
			[]string{"tag5"},
			`[{"Datetime":"2015-12-12T10:10:10.555555Z","Tags":["tag1,tag5"],"Message":"Content of message 14"}]`,
		},
	}

	for _, testCase := range tests {
		t.Log(testCase.Name)
		var reader io.Reader
		var Resp common.Response
		reader = strings.NewReader(testCase.Message) //Convert string to reader
		requestGETMessages, err := http.NewRequest("GET", Endpoint+"/messages", reader)
		res, err := http.DefaultClient.Do(requestGETMessages)
		if err != nil {
			t.Error(err)
		}

		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		jsonResponse := string(response)
		jsonResponse = strings.Trim(jsonResponse, "\n")

		if testCase.expectedStatusCode != 200 {
			decoder := json.NewDecoder(strings.NewReader(jsonResponse))
			err = decoder.Decode(&Resp)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
			assert.Equal(t, testCase.expectedResponseBody, jsonResponse)
			continue
		}

		var Messages []common.Message
		decoder := json.NewDecoder(strings.NewReader(jsonResponse))
		err = decoder.Decode(&Messages)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		assert.Equal(t, testCase.expectedResponseBody, jsonResponse)
	}
}

func TestAddingMessagesLoggedUser(t *testing.T) {
	loginJSON := `{"Username":"admin1","Password":"zecret111"}`
	var reader io.Reader
	var Response common.Response
	reader = strings.NewReader(loginJSON) //Convert string to reader
	requestLogin, err := http.NewRequest("POST", Endpoint+"/users/login", reader)
	res, err := http.DefaultClient.Do(requestLogin)
	if err != nil {
		t.Error(err)
	}

	// Valid login first
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&Response)
	if err != nil {
		t.Error(err)
	}

	validToken := Response.Token

	tests := []TestStructGetMessages{
		{
			"Adding message with no token",
			fmt.Sprintf(`
{
"Message":"%s",
"Tags":""
}`, "Test message with no token"),
			403,
			[]string{""},
			"Your are not allowed to add new messages.",
		},
		{
			"Adding message with empty token",
			fmt.Sprintf(`
{
"Token":"",
"Message":"%s",
"Tags":""
}`, "Test message with empty token"),
			403,
			[]string{""},
			"Your are not allowed to add new messages.",
		},
		{
			"Adding message with invalid token",
			fmt.Sprintf(`
{
"Token":"_%s_",
"Message":"%s",
"Tags":""
}`, validToken, "Test message with invalid token"),
			403,
			[]string{""},
			"Your are not allowed to add new messages.",
		},
		{
			"Adding message with no tags",
			fmt.Sprintf(`
{
"Token":"%s",
"Message":"%s",
"Tags":""
}`, validToken, "Test message 1 with no tags"),
			201,
			[]string{""},
			"You have added a new message.",
		},
		{
			"Adding message with multiple tags",
			fmt.Sprintf(`
{
"Token":"%s",
"Message":"%s",
"Tags":["tag1","tag2","tag3"]
}`, validToken, "Test message 2 with tags"),
			201,
			[]string{"tag1", "tag2", "tag3"},
			"You have added a new message.",
		},
		{
			"Adding too short message",
			fmt.Sprintf(`
{
"Token":"%s",
"Message":"1",
"Tags":""
}`, validToken),
			422,
			[]string{""},
			"Your message is too short or too long (2/180)",
		},
		{
			"Adding too long message",
			fmt.Sprintf(`
{
"Token":"%s",
"Message":"ads sad sad ada da da dad ad ad sa dsad ad ad ad sad sads adsadsad ad sad sads d adsa dsa dsa a sad
asd asd sad sa dad sad sad sa dsa dsa d few f dsfd fds fds dsf dsf ds fds fs fds fsd fds fds fs  f sfs fs fs fs
sf sdf sf ds fsf dsf dsf dsf dsf dsfdsfdsf dsfds fsd sd s  sfdsfdsfds dsdfdsfdsfdsds fdsfdsfdsfds ffsfdsfdsfdsf
sdfdsfsds fdsfdsfsdffdsds fdsfdsfdsfdsf dsfdsdsfdssewrdsdffdf dsfdsfdsdsfdsfs fdsfdsfdsfdsfdsfd sfdsdsfsdsfsds",
"Tags":""
}`, validToken),
			422,
			[]string{""},
			"Your message is too short or too long (2/180)",
		},
		{
			"Adding message with invalid JSON",
			fmt.Sprintf(`

Token":'",
"Message","%s",}:{}{
"Tags:"
}`, "Completely Invalid JSON"),
			422,
			[]string{""},
			"Your message is too short or too long (2/180)",
		},
	}

	for _, testCase := range tests {
		t.Log(testCase.Name)
		var ResponseTestCase common.Response
		var reader io.Reader
		reader = strings.NewReader(testCase.Message) //Convert string to reader
		requestLogin, err := http.NewRequest("POST", Endpoint+"/messages", reader)
		res, err := http.DefaultClient.Do(requestLogin)
		if err != nil {
			t.Error(err)
		}

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&ResponseTestCase)
		if err != nil {
			t.Error(err)
		}

		tagString := &strings.Builder{}
		for i, t := range testCase.expectedTags {
			if i == 0 {
				tagString.WriteString(t)
			} else {
				tagString.WriteString(fmt.Sprintf(", %s", t))
			}
		}

		assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		assert.Equal(t, testCase.expectedResponseBody, ResponseTestCase.Message)
	}
}

func TestAddingMessagesNoPrivilegedUser(t *testing.T) {

	loginJSON := `{"Username":"admin2","Password":"zecret222"}`
	var reader io.Reader
	var Response common.Response
	reader = strings.NewReader(loginJSON) //Convert string to reader
	requestLogin, err := http.NewRequest("POST", Endpoint+"/users/login", reader)
	res, err := http.DefaultClient.Do(requestLogin)
	if err != nil {
		t.Error(err)
	}

	// Valid login first
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&Response)
	if err != nil {
		t.Error(err)
	}

	validToken := Response.Token

	tests := []TestStructGetMessages{
		{
			"No proviledged user - trying to add with no token",
			fmt.Sprintf(`
{
"Token":"%s",
"Message":"%s",
"Tags":""
}`, validToken, "Test message with no token"),
			403,
			[]string{""},
			"Your are not allowed to add new messages.",
		},
	}

	for _, testCase := range tests {
		t.Log(testCase.Name)
		var ResponseTestCase common.Response
		var reader io.Reader
		reader = strings.NewReader(testCase.Message) //Convert string to reader
		requestLogin, err := http.NewRequest("POST", Endpoint+"/messages", reader)
		res, err := http.DefaultClient.Do(requestLogin)
		if err != nil {
			t.Error(err)
		}

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&ResponseTestCase)
		if err != nil {
			t.Error(err)
		}

		tagString := &strings.Builder{}
		for i, t := range testCase.expectedTags {
			if i == 0 {
				tagString.WriteString(t)
			} else {
				tagString.WriteString(fmt.Sprintf(", %s", t))
			}
		}

		assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		assert.Equal(t, testCase.expectedResponseBody, ResponseTestCase.Message)
	}
}

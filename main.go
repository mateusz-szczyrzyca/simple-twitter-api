package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"twitter/common"
	"twitter/database"
)

func main() {
	// Default level for this example is info, unless debug flag is present
	endpoint := flag.String("endpoint", "localhost:58123", "set endpoint on this host:port")
	debug := flag.Bool("debug", false, "sets log level to debug")
	dsn := flag.String("dsn", "", "sets dns for database")

	flag.Parse()
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	database := database.NewDB(logger, *dsn)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger.Info().Msg("server successfully started")

	router := mux.NewRouter()
	var responseCode int
	var responseMessage string
	var tokenMessage string

	//
	// POST: /users/login
	//
	// Expected headers:
	// - None
	//
	// Expected Body JSON:
	// - {"Username":"username", "Password":"password"}
	//
	// Return codes:
	// - 200, login successfull
	// - 401, invalid username/password
	//
	// Response:
	// - JSON with authorization token
	//
	router.HandleFunc("/users/login", func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var User common.UserLogin
		err := decoder.Decode(&User)
		if err != nil {
			logger.Err(err).Msg("error with decoding json request.")
		}

		token := common.GenerateToken()
		hash := common.HashPassword(User.Password)
		logger.Info().
			Str("METHOD", "POST /login").
			Str("Username", User.Username).
			Str("Password", User.Password).
			Str("hashPass", hash).
			Msg("Login action executed")

		responseMessage = fmt.Sprintf("Congratulations, you've provided correct credentials!")
		responseCode = 200
		tokenMessage = token

		queryLogin := "UPDATE users SET token = $1 WHERE username = $2 AND password = $3"
		if database.RowsAffected(queryLogin, token, User.Username, hash) != 1 {
			responseMessage = fmt.Sprintf("Invalid login or password.")
			logger.Warn().Msg(responseMessage)
			responseCode = 401
			tokenMessage = ""
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(responseCode)
		err = json.NewEncoder(writer).Encode(&common.Response{
			Message: responseMessage,
			Token:   tokenMessage,
		})
		if err != nil {
			logger.Err(err)
		}
	}).Methods("POST")

	//
	// POST: /users/logout
	//
	// Expected Body JSON:
	// - {"Token":"TOKENCODE"}
	//
	// Return codes:
	// - 200, logout successful
	// - 403, invalid token
	//
	// Response:
	// - JSON with message
	//
	router.HandleFunc("/users/logout", func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var User common.UserLogin
		err := decoder.Decode(&User)
		if err != nil {
			logger.Err(err).Msg("error with decoding json request.")
		}

		logger.Info().
			Str("METHOD", "POST /logout").
			Str("Token", User.Token).
			Msg("Logout action performed")

		responseMessage = fmt.Sprintf("You have successfully logged out.")
		responseCode = 200
		tokenMessage = ""

		queryLogin := "UPDATE users SET token = NULL WHERE token = $1"
		if database.RowsAffected(queryLogin, User.Token) != 1 {
			responseMessage = fmt.Sprintf("Invalid token.")
			responseCode = 403
			tokenMessage = ""
		}

		logger.Info().Msg(responseMessage)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(responseCode)
		err = json.NewEncoder(writer).Encode(&common.Response{
			Message: responseMessage,
			Token:   tokenMessage,
		})
		if err != nil {
			logger.Err(err).Msg("problem with JSON encode")
		}
	}).Methods("POST")

	//
	// POST: /messages
	//
	// Expected Body JSON:
	// - {"Token":"TOKENCODE", "message":"Message content"}
	//
	// Return codes:
	// - 201, message added
	// - 401, not logged/wrong token
	// - 403, no permissions
	// - 422, too short/long message
	//
	// Response:
	// - JSON with response
	//
	router.HandleFunc("/messages", func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var Message common.Message
		err := decoder.Decode(&Message)
		if err != nil {
			logger.Err(err).Msg("error with decoding json request.")
		}

		// Default responses
		responseMessage = fmt.Sprintf("You have added a new message.")
		responseCode = 201
		tokenMessage = ""

		// TODO: fix this someday
		queryUser := "SELECT uuid,status FROM users WHERE token = $1"
		u := database.FetchUser(queryUser, Message.Token)

		logger.Info().
			Str("METHOD", "POST /messages").
			Str("Token", Message.Token).
			Str("Message", Message.Message).
			Strs("Tags", Message.Tags).
			Str("DB UUID", u.Uuid).
			Str("DB Status", u.Status).
			Msg("Posting new message")

		// Not logged?
		if u.Uuid == "" || u.Status == "" {
			responseCode = 401
			responseMessage = "Your are not logged. Please fetch your auth code first."
		}

		// Wrong permissions
		// Only 'u' can add new messages
		if u.Status != "u" {
			responseCode = 403
			responseMessage = "Your are not allowed to add new messages."
		}

		// Too long or too short message
		if utf8.RuneCountInString(Message.Message) < 2 || utf8.RuneCountInString(Message.Message) > 180 {
			responseCode = 422
			responseMessage = "Your message is too short or too long (2/180)"
		}

		// No errors reported - process messages
		if responseCode == 201 {
			// Assembling tags
			tagString := &strings.Builder{}
			for i, t := range Message.Tags {
				if i == 0 {
					tagString.WriteString(t)
				} else {
					tagString.WriteString(fmt.Sprintf(", %s", t))
				}
			}

			// Adding message to database
			addMessageQuery := "INSERT INTO messages(userid, datetime, tags, text) VALUES($1, now(), $2, $3)"
			if database.RowsAffected(addMessageQuery, u.Uuid, tagString.String(), Message.Message) == 0 {
				logger.Warn().
					Str("METHOD", "POST /messages").
					Msg("there was a problem with adding new message")
			}

		}
		logger.Info().Msg(responseMessage)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(responseCode)
		err = json.NewEncoder(writer).Encode(&common.Response{
			Message: responseMessage,
		})
		if err != nil {
			logger.Err(err).Msg("problem with JSON encode")
		}
	}).Methods("POST")

	//
	// GET: /messages
	//
	// Expected Body JSON:
	// - {"Token":"TOKENCODE", "tags":"", "time_from":"", "time_to":""}
	//
	// Return codes:
	// - 200, messages list
	// - 403, no permissions for time filtering
	//
	// Response:
	// - JSON with response
	//
	router.HandleFunc("/messages", func(writer http.ResponseWriter, request *http.Request) {
		// Default responses
		var timeFrameRequest bool
		var tagsRequest bool
		var limitRequest bool
		responseMessage = fmt.Sprintf("Here's the list of messages.")
		responseCode = 200

		decoder := json.NewDecoder(request.Body)
		var Message common.Message
		err := decoder.Decode(&Message)
		if err != nil {
			logger.Err(err).Msg("error with decoding json request.")
			responseCode = 422
			responseMessage = "Invalid JSON request."
		}

		queryUser := "SELECT uuid,status FROM users WHERE token = $1"
		u := database.FetchUser(queryUser, Message.Token)

		// Not authorized?
		// Limit for 100 results
		if u.Uuid == "" || u.Status == "" {
			limitRequest = true
		}

		// Wrong permissions for timeline
		if Message.TimeFrom != "" || Message.TimeTo != "" {
			if u.Status != "a" {
				responseCode = 403
				responseMessage = "Your are not allowed to filter by timeline."
			}

			// Checking date format
			if u.Status == "a" {
				if len(Message.TimeFrom) != 10 || len(Message.TimeTo) != 10 {
					re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
					if !re.MatchString(Message.TimeFrom) || !re.MatchString(Message.TimeTo) {
						responseCode = 422
						responseMessage = "Wrong date format. It should be (YYYY-DD-MM)"
					}
				} else {
					timeFrameRequest = true
				}
			}
		}
		// No errors reported - process request
		if responseCode == 200 {

			// Assembling tags for query
			tagString := &strings.Builder{}
			for i, t := range Message.Tags {
				if i == 0 {
					tagString.WriteString(fmt.Sprintf("'%%%s%%'", t))
				} else {
					tagString.WriteString(fmt.Sprintf(" || '%%%s%%'", t))
				}
			}

			if len(tagString.String()) > 1 {
				tagsRequest = true
			}

			// Composite query
			queryMessage := &strings.Builder{}
			queryMessage.WriteString("SELECT datetime,tags,text FROM messages")

			if timeFrameRequest {
				queryMessage.WriteString(fmt.Sprintf(" WHERE datetime BETWEEN '%s' AND '%s'", Message.TimeFrom,
					Message.TimeTo))
			}

			if tagsRequest {
				if !timeFrameRequest {
					queryMessage.WriteString(fmt.Sprintf(" WHERE tags LIKE (%s)", tagString))
				} else {
					queryMessage.WriteString(fmt.Sprintf(" AND tags LIKE (%s)", tagString))
				}
			}

			queryMessage.WriteString(" ORDER BY datetime DESC")
			if limitRequest {
				queryMessage.WriteString(" LIMIT 10")
			}

			logger.Debug().
				Str("METHOD", "GET /messages").
				Str("SQL Query", queryMessage.String()).Msg("parameters processed and query generated.")

			messages := database.FetchMessages(queryMessage.String())
			err = json.NewEncoder(writer).Encode(&messages)
			if err != nil {
				logger.Err(err).Msg("problem with JSON encode")
			}
		}

		logger.Info().
			Str("METHOD", "GET /messages").
			Str("Token", Message.Token).
			Str("Time FROM", Message.TimeFrom).
			Str("Time TO", Message.TimeTo).
			Strs("Tags", Message.Tags).
			Int("Response CODE", responseCode).
			Msg("Getting messages")
		logger.Info().Msg(responseMessage)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(responseCode)
		if responseCode != 200 {
			err = json.NewEncoder(writer).Encode(&common.Response{
				Message: responseMessage,
			})
			if err != nil {
				logger.Err(err).Msg("problem with JSON encode")
			}
		}
	}).Methods("GET")

	// Starting server
	err := http.ListenAndServe(*endpoint, router)
	if err != nil {
		panic(err)
	}
}

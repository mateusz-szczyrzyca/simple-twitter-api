## POST /users/login

It's used for obtain auth token if provided user and pass is correct.

Request format: 

```json
{
    "Username":"user",
    "Password":"password"
}
```

Response format:

```json
{
    "message": "Message if login is correct",
    "token": "token_if_available"
}
```

Response codes:

```
200, login successfull
401, invalid username/password
```

No custom headers are required.

#
## POST /users/logout

It's used to reset active token if provided token will be found in the database

Request format:

```json
{
    "TOKEN":"TOKENID"
}
```

Response format:

```json
{
    "message": "Message if login is correct",
    "token": "empty_if_ok"
}
```

Response codes:

```
200, logout  successfull
403, invalid token
```

No custom headers are required.

#
## GET /messages

It fetches messages from database. If user is not authorized (`Token` empty or invalid), then 
system returns only 10 messages. 

`Tags` can be empty or ommited. If more tags is specified then system will show such messages 
that contain ALL of them.

`TimeFrom` and `TimeTo` options only support dates in format `YYYY-MM-DD` and are available for 
users with status = `a`, hence normal authorized users (`u`) cannot use this option. 

Request format:

```json
{
    "Token":"TOKENID",
    "TimeFrom":"2016-01-01",
    "TimeTo":"2016-12-01",
    "Tags":["tag1","tag2"]
}
```

Response format:

```json
[
    {"Datetime":"2016-06-06T10:10:10.555555Z","Tags":["tag1,tag2"],"Message":"Content of message 8"},
    {"Datetime":"2016-03-26T10:10:10.555555Z","Tags":["tag1,tag2"],"Message":"Content of message 19"}
]
```

Response codes:

```
200, messages list
403, no sufficient permission for this action
422, invalid json
```

## POST /messages

It's used to add a new message by authorized users (required status = `u`)

Request format:

```json
{
	"Token":"TOKENID",
	"Message":"Message text",
	"Tags":["tag1","tag2","tag3"]
}
```

Response format:

```json
{
    "message": "Message if post was added",
    "token": ""
}
```

Response codes:

```
201, message added
401, not logged or wrong token
403, no sufficient permission for this action
422, too short/long message
```
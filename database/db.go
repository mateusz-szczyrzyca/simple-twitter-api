package database

import (
	"database/sql"

	"github.com/rs/zerolog"
	"twitter/common"
	"twitter/user"
)

type Database interface {
	RowsAffected(string, ...interface{}) int64
	FetchMessages(string, ...interface{}) []common.Message
	FetchUser(string, ...interface{}) user.SessionUser
}

type DB struct {
	DB     *sql.DB
	Logger zerolog.Logger
}

func NewDB(logger zerolog.Logger, dsn string) Database {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Err(err).Msg("cannot connect to database")
	}

	err = db.Ping()
	if err != nil {
		logger.Err(err).Msg("cannot ping database")
	}
	return &DB{
		DB:     db,
		Logger: logger,
	}
}

func (d *DB) DBHandler() *sql.DB {
	return d.DB
}

func (d *DB) RowsAffected(query string, args ...interface{}) int64 {
	var num int64
	db := d.DB
	logger := d.Logger

	row, err := db.Exec(query, args...)
	if err != nil {
		logger.Err(err).Str("funcName", "RowsAffected()").Msg("problem with executing SQL query")
	}

	num, err = row.RowsAffected()
	if err != nil {
		logger.Err(err).Str("funcName", "RowsAffected()").Msg("problem with executing SQL query")
	}

	return num
}

func (d *DB) FetchMessages(query string, args ...interface{}) []common.Message {
	Messages := make([]common.Message, 0)

	db := d.DB
	logger := d.Logger
	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Err(err).
			Str("func", "FetchMessages()").
			Msg("problem with SQL query")
		return Messages
	}

	for rows.Next() {
		var datetime string
		var tags string
		var text string
		if err := rows.Scan(&datetime, &tags, &text); err != nil {
			logger.Err(err).
				Str("func", "FetchMessages()").
				Msg("something wrong with taking messages from database.")
			continue
		}
		//datetime, tags, text
		Messages = append(Messages, common.Message{
			Datetime: datetime,
			Tags:     []string{tags},
			Message:  text,
		})
	}
	return Messages
}

func (d *DB) FetchUser(query string, args ...interface{}) user.SessionUser {
	var User user.SessionUser

	db := d.DB
	logger := d.Logger

	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Err(err).
			Str("func", "FetchUser()").
			Msg("problem with SQL query")
		return User
	}

	for rows.Next() {
		if err := rows.Scan(&User.Uuid, &User.Status); err != nil {
			logger.Err(err).
				Str("func", "FetchUser()").
				Msg("something wrong with taking user uuid or status:w" +
					" from database")
		}
		break
	}

	return User
}

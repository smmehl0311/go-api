package dao

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/smmehl0311/go-api/db/queries"
)

type User struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func InsertUser(db *sql.DB, username string, password string) error {
	_, err := db.Exec(queries.InsertUserQuery, username, password)
	return err
}

func InsertToken(db *sql.DB, username string, token []byte) error {
	_, err := db.Exec(queries.InsertTokenQuery, username, token)
	return err
}

func GetTokensForUser(db *sql.DB, username string) (res *sql.Rows, err error) {
	res, err = db.Query(queries.GetTokensQuery, username)
	return
}

func DeleteToken(db *sql.DB, token []byte) (err error) {
	_, err = db.Exec(queries.DeleteTokenQuery, token)
	return
}

func DeleteTokensByUser(db *sql.DB, user string) (err error) {
	_, err = db.Exec(queries.DeleteTokensByUser, user)
	return
}

func AuthenticateUser(db *sql.DB, username string, password string) (res *sql.Rows, err error) {
	res, err = db.Query(queries.AuthenticateUserQuery, username, password)
	return
}

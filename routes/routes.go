package routes

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/smmehl0311/go-api/db/dao"
)

type PostResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func InsertUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user User
		json.NewDecoder(r.Body).Decode(&user)
		log.Print(user)

		err := dao.InsertUser(db, user.Username, user.Password)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(PostResponse{true, "Created User: " + user.Username})
		} else {
			pqErr := err.(*pq.Error)
			if pqErr.Constraint == "user_pkey" {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "That username already exists"})
			} else {
				log.Print(pqErr)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Error inserting into database"})
			}
		}
	}
}

func CheckToken(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cookie, err := r.Cookie("auth-token")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Invalid Cookie"})
		} else {
			stringSplit := strings.Split(cookie.Value, ":")
			username := stringSplit[0]
			authToken := stringSplit[1]
			res, err := dao.GetTokensForUser(db, username)
			defer res.Close()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Invalid Cookie"})
			} else {
				var token []byte
				var insertedDate time.Time
				foundMatch := false
				for res.Next() {
					err := res.Scan(&token, &insertedDate)
					if err == nil {
						log.Println(insertedDate)
						log.Println(insertedDate.AddDate(0, 0, 7))
						if insertedDate.AddDate(0, 0, 7).Before(time.Now()) {
							log.Println("deleting token")
							err := dao.DeleteToken(db, token)
							if err != nil {
								log.Println(err)
							}
						} else {
							decodedRequestToken, err := base64.URLEncoding.DecodeString(authToken)
							if err != nil {
								log.Println(err)
								break
							}
							if bytes.Equal(token, decodedRequestToken) {
								log.Println("found token match")
								foundMatch = true
								w.WriteHeader(http.StatusNoContent)
								break
							}
						}
					} else {
						log.Println(err)
					}
				}
				if foundMatch == false {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Invalid Cookie"})
				}
			}
		}
	}
}

func AuthenticateUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user User
		json.NewDecoder(r.Body).Decode(&user)

		res, err := dao.AuthenticateUser(db, user.Username, user.Password)
		defer res.Close()
		if err == nil {
			if res.Next() {
				token, err := GetToken()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Internal Server Error"})
				} else {
					err := dao.InsertToken(db, user.Username, token)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Error reaching the database"})
						log.Println(err)
					} else {
						tokenString := base64.URLEncoding.EncodeToString(token)
						cookie := http.Cookie{
							Name:   "auth-token",
							Value:  user.Username + ":" + tokenString,
							MaxAge: 60 * 15,
							Path:   "/"}
						http.SetCookie(w, &cookie)
						w.WriteHeader(http.StatusNoContent)
					}
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Invalid Credentials"})
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Error reaching the database"})
		}
	}
}

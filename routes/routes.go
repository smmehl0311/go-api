package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

func AuthenticateUser(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user User
		json.NewDecoder(r.Body).Decode(&user)
		log.Print(user)

		res, err := dao.AuthenticateUser(db, user.Username, user.Password)
		if err == nil {
			if res.Next() {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(PostResponse{Success: false, Message: "Error reaching the database"})
		}
	}
}

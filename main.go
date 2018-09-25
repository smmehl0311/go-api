package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/smmehl0311/go-api/routes"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/users", routes.InsertUser(db)).Methods("POST")
	router.HandleFunc("/users/authenticate", routes.AuthenticateUser(db)).Methods("POST")
	router.HandleFunc("/token", routes.CheckToken(db)).Methods("POST")
	fmt.Printf("Listening at port 8000...\n")
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	corsHandler := handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins, handlers.AllowCredentials())
	log.Fatal(http.ListenAndServe(":8000", corsHandler(router)))
}

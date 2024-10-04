package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/perfect-less/go-session-auth/handler"
)

const HOST = ""
const PORT = "8082"
var ALLOWED_ORIGINS = []string{"http://localhost:3000"}

func main() {
	fmt.Println("Starting back-end server.")

    router := mux.NewRouter()

	router.HandleFunc("/", handler.Welcome_handler)
	router.HandleFunc("/login", handler.Login_handler)
	router.HandleFunc("/session_check", handler.CheckSession_handler)
	router.HandleFunc("/refresh_session", handler.RefreshSession_handler)
    router.HandleFunc("/logout", handler.Logout_handler)

    // Cors Options
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Language", "Origin"})
    originsOk := handlers.AllowedOrigins(ALLOWED_ORIGINS)
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
    credentialsOK := handlers.AllowCredentials()

	hostname := "localhost"
	if len(HOST) > 0 {
		hostname = HOST
	}
	log.Printf("listening to http://%s:%s", hostname, PORT)
	log.Fatal(
        http.ListenAndServe(
            fmt.Sprintf("%s:%s", HOST, PORT),
            handlers.CORS(headersOk, originsOk, methodsOk, credentialsOK)(router),
        ),
    )
}

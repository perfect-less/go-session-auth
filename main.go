package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/perfect-less/go-session-auth/handler"
)

const HOST = ""
const PORT = "8082"

func main() {
	fmt.Println("Starting back-end server.")

	http.HandleFunc("/", handler.Welcome_handler)
	http.HandleFunc("/login", handler.Login_handler)
	http.HandleFunc("/session_check", handler.CheckSession_handler)
	http.HandleFunc("/refresh_session", handler.RefreshSession_handler)
	http.HandleFunc("/logout", handler.Logout_handler)

	hostname := "localhost"
	if len(HOST) > 0 {
		hostname = HOST
	}
	log.Printf("listening to http://%s:%s", hostname, PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", HOST, PORT), nil))
}

package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
)

var users []string = []string{"usersatu", "userdua"}
var sessions map[string]string = map[string]string{}


func createNewSession(username string) string {
    var generated bool = false
    var newToken string
    for !generated {
        newToken = uuid.NewString()
        _, exists := sessions[newToken]
        if !exists {
            sessions[newToken] = username
            generated = true
        }
    }
    
    return newToken
}


func Welcome_handler(w http.ResponseWriter, r *http.Request) {
    var sleep_time time.Duration = 5 // seconds
    log.Println("getting request on welcome_handler")

    time.Sleep(sleep_time * time.Second)

    log.Println("responding")
    io.WriteString(w, "Welcome")
}


func Login_handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("getting %s request on Login_handler\n", r.Method)

    if r.Method != "POST" {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Only POST request allowed.",
                    http.StatusBadRequest,
                    ),
                )
        w.WriteHeader(http.StatusBadRequest)
        w.Write(error_message)
        return
    }

    err := r.ParseForm()
    if err != nil {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Trouble when Parsing The Form.",
                    http.StatusInternalServerError,
                    ),
                )
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(error_message)
        return
    }

    username := r.FormValue("username")
    log.Printf("username: %s\n", username)
    if username == "" {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Insufficient Params.",
                    http.StatusBadRequest,
                    ),
                )
        w.WriteHeader(http.StatusBadRequest)
        w.Write(error_message)
        return
    }
    
    var userexists bool = slices.Contains(users, username)
    if !userexists {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - User Doesn't Exist.",
                    http.StatusUnauthorized,
                    ),
                )
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(error_message)
        return
    }

    session_token := createNewSession(username)
    session_cookie := http.Cookie{
        Name: "session-cookie",
        Value: session_token,
        Expires: time.Now().Add(time.Minute*2),
    }
    http.SetCookie(w, &session_cookie)
}


func CheckSession_handler(w http.ResponseWriter, r *http.Request) {
    log.Println("getting request on CheckSession_handler")

    session_cookie, err := r.Cookie("session-cookie")
    if err != nil {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Missing Session Cookie.",
                    http.StatusUnauthorized,
                    ),
                )
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(error_message)
        return
    }

    username, exists := sessions[session_cookie.Value]
    if !exists {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Invalid Session.",
                    http.StatusUnauthorized,
                    ),
                )
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(error_message)
        return
    }
    
    io.WriteString(w, fmt.Sprintf("Valid session for username: %s", username))
}


func Logout_handler(w http.ResponseWriter, r *http.Request) {
    log.Println("getting request on Logout_handler")
    session_cookie, err := r.Cookie("session-cookie")
    if err != nil {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Missing Session Cookie.",
                    http.StatusUnauthorized,
                    ),
                )
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(error_message)
        return
    }
    
    _, exists := sessions[session_cookie.Value]
    if !exists {
        error_message := []byte(
                fmt.Sprintf(
                    "%d - Invalid Session.",
                    http.StatusUnauthorized,
                    ),
                )
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(error_message)
        return
    }

    delete(sessions, session_cookie.Value)
    io.WriteString(
        w,
        fmt.Sprintf(
            "Successfully invalidated session %s",
            session_cookie.Value,
        ),
    )
}

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

const TOKEN_EXPIRE_TIME = 20 * 60 // Seconds

type User struct {
    userid int
    username string
    password string
    Name string
}


type Session struct {
    sessionid int
    session_token string
    userid int
    expire time.Time
}

func (s Session) isExpired () bool {
    return time.Now().After(s.expire)
}


var users []User = []User{
    User{
        userid: 1,
        username: "usersatu",
        password: "password",
        Name: "User Satu",
    },
    User{
        userid: 2,
        username: "userdua",
        password: "pw",
        Name: "User Dua",
    },
}
var sessions []Session = []Session {}


func getSessionByToken(session_token string) (*Session, bool) {
    for i := range sessions {
        session := &sessions[i]
        if session.session_token == session_token {
            return session, true
        }
    }

    return &Session{}, false
}


func getUserByUserID(userid int) (User, bool) {
    for _, user := range users {
        if user.userid == userid {
            return user, true
        }
    }
    
    return User{}, false
}


func getUserByUsername(username string) (User, bool) {
    for _, user := range users {
        if user.username == username {
            return user, true
        }
    }
    
    return User{}, false
}


func validateUser(username string, password string) bool {
    user, exists := getUserByUsername(username)
    if !exists {
        return false
    }

    if user.password != password {
        return false
    }

    return true
}


func createNewSession(username string) Session {
    var generated bool = false
    var newToken string

    var newSessionID int = 0
    for _, session := range sessions {
        if session.sessionid >= newSessionID {
            newSessionID = session.sessionid + 1
        }
    }

    var newSession Session
    for !generated {
        newToken = uuid.NewString()
        _, exists := getSessionByToken(newToken)
        if !exists {
            user, _ := getUserByUsername(username)
            userid := user.userid
            newSession = Session{
                sessionid: newSessionID,
                session_token: newToken,
                userid: userid,
                expire: time.Now().Add(TOKEN_EXPIRE_TIME * time.Second),
            }
            sessions = append(sessions, newSession)
            generated = true
        }
    }
    
    return newSession
}


func removeSessionByToken(session_token string) error {
    sessions = slices.DeleteFunc(
        sessions,
        func(s Session) bool {
            return s.session_token == session_token
        },
    )
    return nil
}


func refreshSessionByToken(session_token string) error {
    session, exists := getSessionByToken(session_token)
    if !exists {
        return fmt.Errorf("Session invalid or missing.")
    }

    session.expire = time.Now().Add(TOKEN_EXPIRE_TIME * time.Second)
    return nil
}


func createSessionCookieFromSession(session *Session) (http.Cookie, error) {
    session_cookie := http.Cookie{
        Name: "session-cookie",
        Value: session.session_token,
        Expires: session.expire,
    }
    return session_cookie, nil
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
        http.Error(w, "Only POST request allowed.", http.StatusBadRequest)
        return
    }

    err := r.ParseMultipartForm(2 << 20)
    if err != nil {
        log.Printf(err.Error())
        http.Error(w, "Trouble when Parsing The Form.", http.StatusInternalServerError)
        return
    }

    username := r.PostFormValue("username")
    password := r.PostFormValue("password")
    log.Printf("username: %s\n", username)
    log.Printf("password: %s\n", password)
    if username == "" || password == "" {
        http.Error(w, "Insufficient Params.", http.StatusBadRequest)
        return
    }
    
    var credential_valid bool = validateUser(username, password)
    if !credential_valid {
        http.Error(w, "Invalid Credentials.", http.StatusUnauthorized)
        return
    }

    session := createNewSession(username)
    session_cookie, err := createSessionCookieFromSession(&session)
    if err != nil {
        http.Error(w, "Trouble when Creating Cookie.", http.StatusInternalServerError)
        return
    }
    http.SetCookie(w, &session_cookie)
}


func CheckSession_handler(w http.ResponseWriter, r *http.Request) {
    log.Println("getting request on CheckSession_handler")

    session_cookie, err := r.Cookie("session-cookie")
    if err != nil {
        http.Error(w, "Missing Session Cookie.", http.StatusUnauthorized)
        return
    }

    session, exists := getSessionByToken(session_cookie.Value)
    userid := session.userid
    if !exists {
        http.Error(w, "Invalid Session.", http.StatusUnauthorized)
        return
    }
    
    user, exists := getUserByUserID(userid)
    if !exists {
        http.Error(w, "invalid Credentials.", http.StatusUnauthorized)
        return
    }

    if session.isExpired() {
        removeSessionByToken(session.session_token)

        http.Error(w, "invalid Credentials.", http.StatusUnauthorized)
        return
    }

    io.WriteString(w, fmt.Sprintf("Valid session for user: %s", user.Name))
}


func RefreshSession_handler(w http.ResponseWriter, r *http.Request) {
    log.Println("getting request on RefreshSession_handler")

    session_cookie, err := r.Cookie("session-cookie")
    if err != nil {
        http.Error(w, "Missing Session Cookie.", http.StatusUnauthorized)
        return
    }

    session, exists := getSessionByToken(session_cookie.Value)
    if !exists {
        http.Error(w, "Invalid Session.", http.StatusUnauthorized)
        return
    }
    
    if session.isExpired() {
        removeSessionByToken(session.session_token)

        http.Error(w, "Session Already Expired.", http.StatusUnauthorized)
        return
    }

    err = refreshSessionByToken(session.session_token)
    if err != nil {
        http.Error(w, "Error when Refreshing Token.", http.StatusInternalServerError)
        return
    }
    
    session, _ = getSessionByToken(session.session_token)
    new_session_cookie, err := createSessionCookieFromSession(session)
    if err != nil {
        http.Error(w, "Trouble when Creating Cookie.", http.StatusInternalServerError)
        return
    }
    http.SetCookie(w, &new_session_cookie)
    io.WriteString(w, "Successfuly refreshing token")
}



func Logout_handler(w http.ResponseWriter, r *http.Request) {
    log.Println("getting request on Logout_handler")
    session_cookie, err := r.Cookie("session-cookie")
    if err != nil {
        http.Error(w, "Missing Session Cookie.", http.StatusUnauthorized)
        return
    }
    
    _, exists := getSessionByToken(session_cookie.Value)
    if !exists {
        http.Error(w, "Invalid Session.", http.StatusUnauthorized)
        return
    }

    removeSessionByToken(session_cookie.Value)
    io.WriteString(
        w,
        fmt.Sprintf(
            "Successfully invalidated session %s",
            session_cookie.Value,
        ),
    )
}

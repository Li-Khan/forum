package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type session struct {
	mx sync.Mutex
	// key - uuid	| value - login
	mapCookie map[string]string
}

const cookieName string = "Forum"

var cookie = session{
	mapCookie: make(map[string]string),
}

func addCookie(w http.ResponseWriter, r *http.Request, login string) error {
	u := uuid.NewV4()

	cookie.mx.Lock()
	cookie.mapCookie[u.String()] = login
	cookie.mx.Unlock()

	expire := time.Now().AddDate(0, 0, 1)
	c := &http.Cookie{
		Name:     cookieName,
		Value:    u.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  expire,
	}
	http.SetCookie(w, c)

	fmt.Println(cookie.mapCookie)
	return nil
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	var ok bool
	if err == nil {
		_, ok = cookie.mapCookie[c.Value]
	}
	return ok
}

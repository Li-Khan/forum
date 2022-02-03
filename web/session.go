package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type session struct {
	mx        sync.Mutex
	mapCookie map[string]string
}

const cookieName string = "Forum"

var cookie = session{
	mapCookie: make(map[string]string),
}

func addCookie(w http.ResponseWriter, r *http.Request, login string) error {
	c, err := r.Cookie(cookieName)
	if err != nil {
		u := uuid.NewV4()

		cookie.mx.Lock()
		cookie.mapCookie[u.String()] = login
		cookie.mx.Unlock()

		expire := time.Now().AddDate(0, 0, 1)
		c = &http.Cookie{
			Name:     cookieName,
			Value:    u.String(),
			Path:     "/",
			HttpOnly: true,
			Expires:  expire,
		}
		http.SetCookie(w, c)
	}
	fmt.Println(cookie.mapCookie)
	return nil
}

func deleteCookie(w http.ResponseWriter, r *http.Request) {
	cookie.mx.Lock()
	defer cookie.mx.Unlock()
	// r.Cookie - не вернет ошибку потому что перед тем как вызвать deleteCookie
	// вызывается isSession
	c, _ := r.Cookie(cookieName)

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	delete(cookie.mapCookie, c.Value)
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	var ok bool
	if err == nil {
		_, ok = cookie.mapCookie[c.Value]
	}
	return ok
}

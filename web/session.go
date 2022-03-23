package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

const cookieName string = "Forum"

var cookie sync.Map

type expires struct {
	exp map[interface{}]time.Time
	mu  sync.Mutex
}

var expireSession = expires{exp: make(map[interface{}]time.Time)}

func addCookie(w http.ResponseWriter, r *http.Request, login string) {
	expireSession.mu.Lock()
	defer expireSession.mu.Unlock()

	u := uuid.NewV4()
	oneUser(login, u.String())

	cookie.Store(u.String(), login)
	expire := time.Now().AddDate(0, 0, 1)

	expireSession.exp[u.String()] = expire

	c := &http.Cookie{
		Name:     cookieName,
		Value:    u.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  expire,
	}
	http.SetCookie(w, c)
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	var ok bool
	if err == nil {
		_, ok = cookie.Load(c.Value)
	}
	return ok
}

func oneUser(login, uuid string) {
	cookie.Range(func(key, value interface{}) bool {
		if login == fmt.Sprint(value) {
			cookie.Delete(key)
		}
		return true
	})
}

func getLogin(r *http.Request) string {
	c, _ := r.Cookie(cookieName)
	value, _ := cookie.Load(c.Value)
	login := fmt.Sprint(value)
	return login
}

// SessionGC ...
func SessionGC() {
	for {
		cookie.Range(func(key, value interface{}) bool {
			expireSession.mu.Lock()
			if time.Now().Unix() > expireSession.exp[key].Unix() {
				cookie.Delete(key)
			}
			expireSession.mu.Unlock()
			return true
		})
		time.Sleep(time.Second * 1)
	}
}

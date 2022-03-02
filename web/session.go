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

var expireSession expires

func addCookie(w http.ResponseWriter, r *http.Request, login string) {
	expireSession.mu.Lock()
	defer expireSession.mu.Unlock()

	expireSession.exp = make(map[interface{}]time.Time)

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

// SessionGC ...
func SessionGC() {
	for {
		cookie.Range(func(key, value interface{}) bool {
			expireSession.mu.Lock()
			if expireSession.exp[key].Unix() < time.Now().Unix() {
				cookie.Delete(key)
			}
			expireSession.mu.Unlock()
			return true
		})
		time.Sleep(time.Second * 10)
	}
}

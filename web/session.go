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

func addCookie(w http.ResponseWriter, r *http.Request, login string) {
	u := uuid.NewV4()

	oneUser(login, u.String())

	cookie.Store(u.String(), login)

	expire := time.Now().AddDate(0, 0, 1)
	c := &http.Cookie{
		Name:     cookieName,
		Value:    u.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  expire,
	}
	http.SetCookie(w, c)

	// не забыть удалить все что ниже
	cookie.Range(func(key, value interface{}) bool {
		fmt.Printf("key = %v\tvalue - %v", key, value)
		return true
	})
	fmt.Println()
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

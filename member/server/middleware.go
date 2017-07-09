package server

import (
	"context"
	"net/http"
	"strings"
)

func getSessionFromRequest(r *http.Request) *Session {
	val := r.Header.Get("Authorization")
	if val == "" {
		return nil
	}
	bits := strings.Split(val, " ")
	if len(bits) != 2 {
		return nil
	}
	if strings.ToLower(bits[0]) != "basic" {
		return nil
	}
	session := SessionByKey(bits[1])
	if session != nil {
		go session.Touch()
	}
	return session
}

func Middleware(next http.Handler) http.Handler {
	// Set Session in the http.Request context by checking the
	// Authorization HTTP header.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := getSessionFromRequest(r)
		if session != nil {
			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

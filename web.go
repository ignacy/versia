package main

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

var (
	pageName = whmux.NewStringArg()
)

type ErrHandler interface {
	HandleError(w http.ResponseWriter, req *http.Request, err error)
}

type userKey int

func RequireUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, err := GetUser(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, userKey(0), user)
		h.ServeHTTP(w, req.WithContext(ctx))
	})
}

func main() {
	mux := http.NewServerMux()
	mux.Handle("/hello/", http.HandlerFunc(hello))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static-assets"))))
	http.ListenAndServe(":8080", mux)
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello, world!"))
}

type Context interface {
	Done() <-chan struct{}
	Err() error
	Deadline() (deadline time.Time, ok bool)

	Value(key interface{}) interface{}
}

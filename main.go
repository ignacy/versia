package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ignacy/versia/storage"
	_ "github.com/lib/pq"
)

func invoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/invoice/"):]
	i, _ := strconv.Atoi(id)
	versions := storage.FindVersions(i)

	fmt.Fprintf(w, "<h1>%s</h1>", id)
	for _, v := range versions {
		fmt.Fprintf(w, "<h4>%s %s %s</h4>", v.Event, v.Id, v.Whodunnit)
		fmt.Fprintf(w, "<div>%s</div>", v.Object)
		fmt.Fprintf(w, "<div>%s</div>", v.Object_changes)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello there!</h1>")
}

func main() {
	http.HandleFunc("/", BasicAuth(handleIndex, "admin", "123456", "Auth"))
	http.HandleFunc("/invoice/", BasicAuth(invoiceHandler, "admin", "123456", "Auth"))
	http.ListenAndServe(":8080", nil)
}

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

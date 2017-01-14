package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/webhelp.v1/whlog"

	"github.com/alecthomas/template"
	"github.com/ignacy/versia/storage"
	_ "github.com/lib/pq"
)

var (
	basicAuth = os.Getenv("VERSIA_ADMIN_PASSWORD")
)

func invoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/invoice/"):]
	i, _ := strconv.Atoi(id)
	versions := storage.FindVersions(i)

	t, err := template.ParseFiles("versions.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, versions)
	if err != nil {
		panic(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	invoices := storage.ListInvoices()
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, invoices)
	if err != nil {
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(BasicAuth(handleIndex, "admin", basicAuth, "Auth")))
	mux.Handle("/invoice/", http.HandlerFunc(BasicAuth(invoiceHandler, "admin", basicAuth, "Auth")))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static-assets"))))
	whlog.ListenAndServe(":"+port, whlog.LogResponses(whlog.Default, mux))
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

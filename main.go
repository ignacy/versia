package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/webhelp.v1/whlog"

	"github.com/alecthomas/template"
	"github.com/ignacy/versia/storage"
	_ "github.com/lib/pq"
)

var (
	basicAuth = os.Getenv("VERSIA_ADMIN_PASSWORD")
	modelName = os.Getenv("VERSIA_MODEL_NAME")
)

func invoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"+modelName+"/"):]
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
	models := storage.ListModels()
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, models)
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
	mux.Handle("/", http.HandlerFunc(Etags(BasicAuth(handleIndex, "admin", basicAuth, "Auth"), "list")))
	mux.Handle("/"+modelName+"/", http.HandlerFunc(BasicAuth(invoiceHandler, "admin", basicAuth, "Auth")))
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

func Etags(handler http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := `"` + key + `"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age=7200") // 2 hours

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		handler(w, r)
	}
}

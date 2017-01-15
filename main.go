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
	basicAuth       = os.Getenv("VERSIA_ADMIN_PASSWORD")
	modelName       = os.Getenv("VERSIA_MODEL_NAME")
	pgString        = os.Getenv("VERSIA_PG_STRING")
	username        = "admin"
	basicAuthPrompt = "Authorization:"
)

func modelHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"+modelName+"/"):]
	i, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	log.Printf("Showing version details for id = %d", i)

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

func BasicAuth(handler http.HandlerFunc, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+basicAuthPrompt+`"`)
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

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	storage.InitDB(pgString)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(Etags(BasicAuth(handleIndex, basicAuth), "list")))
	mux.Handle("/"+modelName+"/", http.HandlerFunc(BasicAuth(modelHandler, basicAuth)))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static-assets"))))
	whlog.ListenAndServe(":"+port, whlog.LogResponses(whlog.Default, mux))
}

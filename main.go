package main

import (
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

func main() {
	http.HandleFunc("/invoice/", invoiceHandler)
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ignacy/versia/storage"
	_ "github.com/lib/pq"
)

func main() {
	versions := storage.FindVersions(29)
	for _, v := range versions {
		spew.Dump(v)
	}
}

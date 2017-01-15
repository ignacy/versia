package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ignacy/versia/storage"
)

type mockDB struct{}

func (mdb *mockDB) ListModels() ([]storage.Model, error) {
	models := make([]storage.Model, 0)
	models = append(models, storage.Model{10})
	models = append(models, storage.Model{20})
	return models, nil
}

func (mdb *mockDB) FindVersions(id int) ([]storage.Version, error) {
	return make([]storage.Version, 0), nil
}

func TestListModels(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	env := Env{db: &mockDB{}}
	http.HandlerFunc(env.handleIndex).ServeHTTP(rec, req)

	expectedFirst := `10"><small>10 </small></a>`
	expectedSecond := `20"><small>20 </small></a>`

	if !strings.Contains(rec.Body.String(), expectedFirst) {
		t.Errorf("\n...expected to see %v\n...in = %v", expectedFirst, rec.Body.String())
	}

	if !strings.Contains(rec.Body.String(), expectedSecond) {
		t.Errorf("\n...expected to see %v\n...in = %v", expectedSecond, rec.Body.String())
	}
}

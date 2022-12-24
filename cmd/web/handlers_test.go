package main

import (
	"github.com/manny-e1/snippetbox/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	//rr := httptest.NewRecorder()
	//r, err := http.NewRequest(http.MethodGet, "/", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//Ping(rr, r)
	//rs := rr.Result()
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

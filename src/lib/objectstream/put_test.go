package objectstream

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func putHandler(w http.ResponseWriter, r *http.Request) {
	// b,_ := ioutil.ReadAll(r.Body)
	b, _ := io.ReadAll(r.Body)
	if string(b) == "test" {
		return
	}
	w.WriteHeader(http.StatusForbidden)
}

func TestPut(t *testing.T) {
	s:= httptest.NewServer(http.HandlerFunc(putHandler))
	defer s.Close()

	ps:=NewPutStream(s.URL[7:],"any")
	io.WriteString(ps, "test")
	err := ps.Close()
	if err != nil {
		t.Error(err)
	}
}
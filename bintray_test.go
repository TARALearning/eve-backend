package eve

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_BinTrayDeleteFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		if (r.Method != http.MethodDelete) || (r.URL.Path != "/content/subject/repo//test/filepath.txt") {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server ERROR"))
		}
	}))
	defer ts.Close()
	binTrayURL = ts.URL
	binTrayUser = "testuser"
	binTrayPassword = "secret"
	err := BinTrayDeleteFile("subject", "repo", "/test/filepath.txt")
	if err != nil {
		t.Error(err)
	}
}

func Test_BinTrayPublishFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		if (r.Method != http.MethodPut) || (r.URL.Path != "/file_metadata/subject/repo//test/filepath.txt") {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server ERROR"))
		}
	}))
	defer ts.Close()
	binTrayURL = ts.URL
	binTrayUser = "testuser"
	binTrayPassword = "secret"
	err := BinTrayPublishFile("subject", "repo", "/test/filepath.txt")
	if err != nil {
		t.Error(err)
	}
}

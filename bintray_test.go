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

func Test_BinTrayDownloadFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		if (r.Method != http.MethodGet) || (r.URL.Path != "/subject/repo//test/filepath.txt") {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server ERROR"))
		}
		w.WriteHeader(200)
		w.Write([]byte("content"))
	}))
	defer ts.Close()
	binTrayDownloadURL = ts.URL
	binTrayUser = "testuser"
	binTrayPassword = "secret"
	content, err := BinTrayDownloadFile("subject", "repo", "/test/filepath.txt")
	if err != nil {
		t.Error(err)
	}
	if string(content) != "content" {
		t.Error("BinTrayDownloadFile does not work as expected")
	}
}

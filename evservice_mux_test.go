package eve

import (
	"net/http"
	"testing"
)

func Test_MuxValueNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1/users/test", nil)
	if err != nil {
		t.Error(err)
	}
	testuser := MuxValue(req, "user")
	if testuser != "" {
		t.Error("MuxValue does not work as expected")
	}
}

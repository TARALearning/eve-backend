package eve

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MuxValue returns the given gorilla mux uri value
func MuxValue(r *http.Request, key string) string {
	vars := mux.Vars(r)
	if val, ok := vars[key]; ok {
		return val
	}
	return ""
}

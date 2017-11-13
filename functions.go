package eve

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MuxValue(r *http.Request, key string) string {
	vars := mux.Vars(r)
	if val, ok := vars[key]; ok {
		return val
	}
	return ""
}

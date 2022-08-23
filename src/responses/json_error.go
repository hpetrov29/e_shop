package responses

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	response := make(map[string]string)
	response["message"] = error
	jsonResp, _ := json.Marshal(response)
	w.WriteHeader(code)
	w.Write(jsonResp)
}

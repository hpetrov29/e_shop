package responses

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, message string, content interface{}, code int) {
	payload := map[string]interface{}{"success": 1, "message": message, "paylaod": content}
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

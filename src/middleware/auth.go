package middleware

import (
	"net/http"

	"github.com/fnmzgdt/e_shop/src/responses"
)

func Authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwt, err := r.Cookie("Auth-token")
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		result, err := Validate(jwt.Value)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		res2 := result.(map[string]interface{})
		userId := res2["userId"].(string)
		r.Header.Add("userId", userId)
		handler.ServeHTTP(w, r)
	}
}

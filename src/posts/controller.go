package posts

import (
	"encoding/json"
	"net/http"

	"github.com/fnmzgdt/e_shop/src/responses"
)

func writePost(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//get userId from token in middleware
		post := NewPost(999) //argument = id of writer
		_ = json.NewDecoder(r.Body).Decode(&post)

		err := post.checkFields()
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = s.WritePost(&post)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responses.JSONResponse(w, "result", "Successful entry.", 200)
	}
}

package posts

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writePost(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var post Post
		_ = json.NewDecoder(r.Body).Decode(&post)
		err := s.WritePost(&post)
		if err == nil {
			response, _ := json.Marshal("success")
			w.WriteHeader(200)
			w.Write(response)
		}
		fmt.Println(err)
		return
	}
}

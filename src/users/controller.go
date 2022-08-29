package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/fnmzgdt/e_shop/src/responses"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		user := NewUser()
		_ = json.NewDecoder(r.Body).Decode(&user)

		if err := user.checkFields(); err != nil {
			responses.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = string(password[:])
		userId, err := s.InsertUser(&user)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = ""
		claims := user.createClaims(userId)
		jwt, err := middleware.NewJWT(time.Minute*30, claims)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionId := claims.SessionUUID
		claims.SessionUUID = ""
		claimsJson, err := json.Marshal(claims)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.CreateSession(userId, sessionId, string(claimsJson))
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{Name: "Auth-token", Value: jwt, Path: "/", Expires: time.Now().Add(time.Minute * 60), Secure: true, HttpOnly: true}
		http.SetCookie(w, &cookie)

		responses.JSONResponse(w, "Successful registration.", []User{user}, 200)
		return
	}
}

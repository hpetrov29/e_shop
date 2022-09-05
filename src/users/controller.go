package users

import (
	"encoding/json"
	"net/http"
	"strings"
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
			if strings.Split(err.Error(), ":")[0] == "Error 1062" {
				responses.JSONError(w, "An account with this email already exists.", http.StatusInternalServerError)
				return
			}
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = ""
		claims := user.createClaims(userId)
		refreshToken, err := middleware.NewJWT(time.Hour*24*356, map[string]interface{}{"sessionId": claims.SessionUUID, "userId": claims.UserId})
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionId := claims.deleteSessionId()
		accessToken, err := middleware.NewJWT(time.Minute*5, claims)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claimsJson, err := json.Marshal(claims)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = s.CreateSession(userId, sessionId, string(claimsJson)); err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		refreshCookie := http.Cookie{Name: "refreshToken", Value: refreshToken, Path: "/", Expires: time.Now().Add(time.Hour * 24 * 356), Secure: true, HttpOnly: true}
		http.SetCookie(w, &refreshCookie)

		accessCookie := http.Cookie{Name: "accessToken", Value: accessToken, Path: "/", Expires: time.Now().Add(time.Minute * 5), Secure: true, HttpOnly: true}
		http.SetCookie(w, &accessCookie)

		responses.JSONResponse(w, "Successful registration.", []User{user}, 200)
		return
	}
}

func login(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responses.JSONError(w, "Method type not allowed.", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		userLogin := UserLogin{}
		_ = json.NewDecoder(r.Body).Decode(&userLogin)
		password, err := s.GetPasswordFromEmail(&userLogin)
		if err != nil {
			responses.JSONError(w, "Wrong email or password.", http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userLogin.Password)); err != nil {
			responses.JSONError(w, "Wrong email or password.", http.StatusBadRequest)
			return
		}
		claims, err := s.GetClaimsFromEmail(&userLogin)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claims.addSessionId()
		refreshToken, err := middleware.NewJWT(time.Hour*24*356, map[string]interface{}{"sessionId": claims.SessionUUID, "userId": claims.UserId})
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionUUID := claims.deleteSessionId()
		accessToken, err := middleware.NewJWT(time.Minute*5, *claims)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claimsJson, err := json.Marshal(claims)
		if err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = s.CreateSession(claims.UserId, sessionUUID, string(claimsJson)); err != nil {
			responses.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refreshCookie := http.Cookie{Name: "refreshToken", Value: refreshToken, Path: "/", Expires: time.Now().Add(time.Hour * 24 * 356), Secure: true, HttpOnly: true}
		http.SetCookie(w, &refreshCookie)

		accessCookie := http.Cookie{Name: "accessToken", Value: accessToken, Path: "/", Expires: time.Now().Add(time.Minute * 5), Secure: true, HttpOnly: true}
		http.SetCookie(w, &accessCookie)

		responses.JSONResponse(w, "Successful Login.", []UserClaims{*claims}, 200)
		return
	}
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"to-do/internal/http/cookie"
	"to-do/pkg/response"

	"github.com/go-playground/validator"
)

func (h *Handlers) SighUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user cookie.User
	json.NewDecoder(r.Body).Decode(&user)
	err := validator.New().Struct(user)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response.NewError("Validate error"))
		return
	}

	id, err := h.db.Register(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response.NewError("User is have on server"))
		return
	}

	cookie, err := cookie.GetCookieUser(h.cfg.Cookie.Key, cookie.User{
		UserID:   id,
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(response.NewError("Cookie error"))
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response.NewOk())
}

func (h *Handlers) SighIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user cookie.User
	json.NewDecoder(r.Body).Decode(&user)
	err := validator.New().Struct(user)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response.NewError("Validate error"))
		return
	}

	id, err := h.db.Login(user.Username, user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(response.NewError("User not found"))
			return
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response.NewError("Database error"))
		return
	}

	cookieData, err := cookie.GetCookieUser(h.cfg.Cookie.Key, cookie.User{
		UserID:   int64(id),
		Username: user.Username,
		Password: user.Password,
	})

	http.SetCookie(w, cookieData)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response.NewOk())
}

func (h *Handlers) WhoI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := getCookie(r, h.cfg)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(response.NewError("Please, autarithazion!"))
			return
		} else {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(response.NewError("Database error"))
			return
		}
	}

	json.NewEncoder(w).Encode(cookie.User{UserID: user.UserID, Username: user.Username})
}

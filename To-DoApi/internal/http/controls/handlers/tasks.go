package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"to-do/internal/sqlite"
	"to-do/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
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

	tasks, err := h.db.GetTasks(int(user.UserID))
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response.NewError("Database error"))
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idData := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idData)
	if err != nil {
		w.WriteHeader(204)
		json.NewEncoder(w).Encode(response.NewError("Invalid ID"))
		return
	}

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

	task, err := h.db.GetTask(int(user.UserID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(response.NewError("Not found"))
			return
		}
	}

	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) SetTask(w http.ResponseWriter, r *http.Request) {
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

	var task sqlite.Task
	json.NewDecoder(r.Body).Decode(&task)
	err = validator.New().Struct(task)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewError("Validate error"))
		return
	}

	if task.Priorety == "" {
		task.Priorety = "normal"
	}

	task.Date = strings.Split(time.Now().String(), " ")[0]

	id, err := h.db.AddTask(int(user.UserID), task)
	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(response.NewError("User not found"))
		}
		json.NewEncoder(w).Encode(response.NewError("Database error"))
		return
	}

	json.NewEncoder(w).Encode(response.NewMessage(fmt.Sprint("ID: ", id)))
}

func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idData := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idData)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewMessage("error ID"))
		return
	}

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

	var task sqlite.Task
	json.NewDecoder(r.Body).Decode(&task)

	if task.Title != "" {
		err := h.db.UpdateTitle(id, int(user.UserID), task.Title)
		if err != nil {
			log.Print(err)
			json.NewEncoder(w).Encode(response.NewError("Update title error"))
			return
		}
	}
	if task.Description != "" {
		err := h.db.UpdateDesc(id, int(user.UserID), task.Description)
		if err != nil {
			json.NewEncoder(w).Encode(response.NewError("Update description error"))
			return
		}
	}
	if task.Priorety != "" {
		err := h.db.UpdatePriorety(id, int(user.UserID), task.Priorety)
		if err != nil {
			json.NewEncoder(w).Encode(response.NewError("Update priorety error"))
			return
		}
	}

	json.NewEncoder(w).Encode(response.NewOk())
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idData := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idData)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewMessage("Error ID"))
		return
	}

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

	err = h.db.DeleteTask(id, int(user.UserID))
	if err != nil {
		json.NewEncoder(w).Encode(response.NewError("Error delete"))
		return
	}
}

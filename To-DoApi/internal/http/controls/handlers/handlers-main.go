package handlers

import (
	"encoding/json"
	"net/http"
	"to-do/internal/config"
	"to-do/internal/http/controls/interfaces"
	"to-do/internal/http/cookie"
	"to-do/internal/logger"
	"to-do/internal/sqlite"
	mwlogger "to-do/pkg/mw_logger"
	"to-do/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	db     *sqlite.Database
	cfg    *config.Config
	logger *logger.Logger
}

func NewHandler(db *sqlite.Database, cfg *config.Config, logger *logger.Logger) interfaces.Handler {
	return &Handlers{db: db, cfg: cfg, logger: logger}
}

func (h *Handlers) Register(r *chi.Mux) {
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(mwlogger.New(h.logger.MW, h.cfg))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.NewOk())
	})

	r.Get("/tasks", h.GetTasks)
	r.Get("/tasks/{id}", h.GetTask)
	r.Post("/tasks", h.SetTask)
	r.Patch("/tasks/{id}", h.PatchTask)
	r.Delete("/tasks/{id}", h.DeleteTask)

	r.Get("/auth/me", h.WhoI)
	r.Post("/auth/sigh-in", h.SighIn)
	r.Post("/auth/sigh-up", h.SighUp)
}

func getCookie(r *http.Request, cfg *config.Config) (*cookie.User, error) {
	cookieData, err := r.Cookie("session_cookie")
	if err != nil {
		return nil, err
	}

	dataUser, err := cookie.ValidateCookieUser(cfg.Cookie.Key, cookieData.Value)
	if err != nil {
		return nil, err
	}

	var user cookie.User
	json.Unmarshal([]byte(dataUser), &user)
	return &user, nil
}

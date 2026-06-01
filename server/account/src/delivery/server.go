package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/BalabanovA898/project/account/src/domain"
	"github.com/BalabanovA898/project/account/src/usecase"
)

type Server struct {
	authUseCase *usecase.AuthUseCase
	userUseCase *usecase.UserUseCase
	mux         *http.ServeMux
}

func NewServer(authUseCase *usecase.AuthUseCase, userUseCase *usecase.UserUseCase) *Server {
	s := &Server{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		mux:         http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/api/v1/auth/register", s.handleRegister)
	s.mux.HandleFunc("/api/v1/auth/login", s.handleLogin)
	s.mux.HandleFunc("/api/v1/auth/refresh", s.handleRefresh)
	s.mux.HandleFunc("/api/v1/users/", s.handleUserByID)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": message})
}

func writeUseCaseError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrConflict:
		writeError(w, http.StatusConflict, err.Error())
	case domain.ErrInvalidCredential, domain.ErrUnauthorized:
		writeError(w, http.StatusUnauthorized, err.Error())
	case domain.ErrNotFound:
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

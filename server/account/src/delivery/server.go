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
	teamUseCase *usecase.TeamUseCase
	jwtSecret   string
	mux         *http.ServeMux
}

func NewServer(authUseCase *usecase.AuthUseCase, userUseCase *usecase.UserUseCase, teamUseCase *usecase.TeamUseCase, jwtSecret string) *Server {
	s := &Server{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		teamUseCase: teamUseCase,
		jwtSecret:   jwtSecret,
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

	// Protected routes
	s.mux.HandleFunc("/api/v1/users/", s.AuthMiddleware(s.handleUserByID))
	s.mux.HandleFunc("/api/v1/teams", s.AuthMiddleware(s.handleCreateTeam))
	s.mux.HandleFunc("/api/v1/teams/", s.handleTeamsRouter)
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
	// Проверяем ValidationError
	if ve, ok := err.(*domain.ValidationError); ok {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":  "validation failed",
			"field":  ve.Field,
			"detail": ve.Message,
		})
		return
	}

	// Проверяем другие ошибки
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

func (s *Server) handleTeamsRouter(w http.ResponseWriter, r *http.Request) {
	// Маршрутизируем /api/v1/teams/{id}/* пути
	path := r.URL.Path

	// /api/v1/teams/{id}
	if !contains(path, "/join") && !contains(path, "/leave") {
		s.AuthMiddleware(s.handleGetTeam)(w, r)
		return
	}

	// /api/v1/teams/{id}/join
	if contains(path, "/join") {
		s.AuthMiddleware(s.handleJoinTeam)(w, r)
		return
	}

	// /api/v1/teams/{id}/leave
	if contains(path, "/leave") {
		s.AuthMiddleware(s.handleLeaveTeam)(w, r)
		return
	}

	writeError(w, http.StatusNotFound, "not found")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

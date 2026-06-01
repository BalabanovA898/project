package delivery

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) handleCreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	team, err := s.teamUseCase.CreateTeam(r.Context(), req.Name, userID)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, team)
}

func (s *Server) handleGetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	teamID := strings.TrimPrefix(r.URL.Path, "/api/v1/teams/")
	teamID = strings.Split(teamID, "/")[0]

	team, err := s.teamUseCase.GetTeam(r.Context(), teamID)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, team)
}

func (s *Server) handleJoinTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Извлекаем teamID из пути /api/v1/teams/{id}/join
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/teams/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}
	teamID := parts[0]

	if err := s.teamUseCase.JoinTeam(r.Context(), teamID, userID); err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "successfully joined team"})
}

func (s *Server) handleLeaveTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Извлекаем teamID из пути /api/v1/teams/{id}/leave
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/teams/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}
	teamID := parts[0]

	if err := s.teamUseCase.LeaveTeam(r.Context(), teamID, userID); err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "successfully left team"})
}

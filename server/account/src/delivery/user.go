package delivery

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BalabanovA898/project/account/src/domain"
	"github.com/BalabanovA898/project/account/src/usecase"
)

func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	if id == "" {
		writeError(w, http.StatusNotFound, "missing user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := s.userUseCase.GetUserByID(r.Context(), id)
		if err != nil {
			writeUseCaseError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, toUserPreviewDTO(user))
	case http.MethodPatch:
		var req struct {
			Username    *string              `json:"username,omitempty"`
			Preferences *domain.Preferences `json:"preferences,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		user, err := s.userUseCase.UpdateUser(r.Context(), id, usecase.UserPatch{Username: req.Username, Preferences: req.Preferences})
		if err != nil {
			writeUseCaseError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, toUserPreviewDTO(user))
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

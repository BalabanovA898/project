package delivery

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, tokens, err := s.authUseCase.Register(r.Context(), req.Email, req.Password, req.Username)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, tokens)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, tokens, err := s.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := s.authUseCase.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

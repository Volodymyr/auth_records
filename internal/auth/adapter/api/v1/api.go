package api

import (
	"auth_records/internal/auth/adapter/api/v1/dto"
	"auth_records/pkg/utils"
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (string, error)
}

type api struct {
	log         *zap.Logger
	authUseCase AuthUseCase
}

func New(log *zap.Logger, authUseCase AuthUseCase) *api {
	return &api{
		log,
		authUseCase,
	}
}

func (a *api) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	loginRequest := &dto.LoginRequest{}

	err := utils.DecodeJSON(a.log, w, r, loginRequest)
	if err != nil {
		return
	}

	ctx := r.Context()

	token, err := a.authUseCase.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

package api

import (
	"auth_records/internal/auth/adapter/api/v1/dto"
	"auth_records/pkg/utils"
	"context"
	"encoding/json"
	"net/http"

	pb "auth_records/pkg/records_grpc/v1"

	"go.uber.org/zap"
)

type RecordsClient interface {
	GetRandomRecords(token string) ([]*pb.Record, error)
}

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (string, error)
}

type api struct {
	log           *zap.Logger
	authUseCase   AuthUseCase
	recordsСlient RecordsClient
}

func New(log *zap.Logger, authUseCase AuthUseCase, recordsСlient RecordsClient) *api {
	return &api{
		log,
		authUseCase,
		recordsСlient,
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

	records, err := a.recordsСlient.GetRandomRecords(token)
	if err != nil {
		a.log.Error("Failed to get records", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to fetch records"))
		return
	}

	recordsResponse := dto.RecordsResponse{
		Records: records,
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(recordsResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

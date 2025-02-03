package dto

import pb "auth_records/pkg/records_grpc/v1"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RecordsResponse struct {
	Records []*pb.Record `json:"records"`
}

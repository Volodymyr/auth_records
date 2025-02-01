package utils

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func DecodeJSON(log *zap.Logger, w http.ResponseWriter, r *http.Request, requestMapper any) error {
	err := json.NewDecoder(r.Body).Decode(requestMapper)
	if err != nil {
		log.Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return err
	}
	r.Body.Close()

	return nil
}

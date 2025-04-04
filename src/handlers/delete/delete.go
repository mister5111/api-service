package delete

import (
	cv "api-service/src/lib/custom_validator"
	"api-service/src/storage"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type delRequest struct {
	Alias string `json:"alias" validate:"required,alias"`
}

type deleted struct {
	Message string `json:"Message"`
	Alias   string `json:"Alias"`
	Status  string `json:"Status"`
}

type deleteRows interface {
	Delete(alias string) error
}

func Del(log *slog.Logger, rows deleteRows) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		bodyText, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("BodyText", slog.String("err", err.Error()))
		}
		defer r.Body.Close()

		var delReq delRequest
		err = json.Unmarshal(bodyText, &delReq)
		if err != nil {
			log.Error("Unmarshal", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Request body empty or invalid")
			return
		}

		log.Info("Request body decoded", "Request", delReq)

		if !cv.AliasValidator(delReq.Alias) {
			log.Error("Validation error", slog.String("err", "Wrong characters in Alias"))

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Wrong characters in Alias")
			return
		}

		err = rows.Delete(delReq.Alias)
		if errors.Is(err, storage.ErrALIASNotFound) {
			log.Error("Delete Alias", slog.String("Alias", delReq.Alias), slog.String("err", "not found"))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Alias not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deleted{
			Message: "Alias deleted",
			Alias:   delReq.Alias,
			Status:  "Successful",
		})
	}

}

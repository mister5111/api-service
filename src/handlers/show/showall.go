package show

import (
	cv "api-service/src/lib/custom_validator"
	"api-service/src/storage"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func ShowAll(log *slog.Logger, show getRowsAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ShowAll"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		defer r.Body.Close()

		bodyText := json.NewDecoder(r.Body)
		bodyText.DisallowUnknownFields()

		var showReq showFromAlias
		err := bodyText.Decode(&showReq)
		if err != nil {
			log.Error("Decode", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Request body empty or invalid")
			return
		}

		if !cv.AliasValidator(showReq.Alias) {
			log.Error("Validation error", slog.String("err", "Wrong characters in Alias"))

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Wrong characters in Alias")
			return
		}

		results, err := show.ShowAll()
		if errors.Is(err, storage.ErrALIASNotFound) {
			log.Error("Get All Rows", slog.String("err", storage.ErrALIASNotFound.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Get All Rows Error")
			return
		}

		log.Info("Get All Rows", "Rows", results)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	}
}

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

func Show(log *slog.Logger, show getRows) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Show"

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

		results, err := show.ShowAlias(showReq.Alias)
		if errors.Is(err, storage.ErrALIASNotFound) {
			log.Error("Show Alias", slog.String("Alias", showReq.Alias), slog.String("err", storage.ErrALIASNotFound.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Alias not found")
			return
		}

		log.Info("Get Alias", "Rows", results)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	}
}

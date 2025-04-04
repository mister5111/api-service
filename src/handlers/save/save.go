package save

import (
	resp "api-service/src/lib/response"
	"api-service/src/storage"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type saved struct {
	Message string `json:"Message"`
	Alias   string `json:"Alias"`
	Url     string `json:"Url"`
	IDRows  int64  `json:"ID"`
	Status  string `json:"Status"`
}

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		bodyText, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("BodyText", slog.String("err", err.Error()))
		}
		defer r.Body.Close()

		var req Request
		err = json.Unmarshal(bodyText, &req)
		if err != nil {
			log.Error("Unmarshal", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("request body empty or invalid")
			return
		}

		log.Info("Request body decoded", "Request", req)

		if err := validator.New().Struct(req); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			log.Error("Validation error", "Error", resp.ValidationError(validationErrs))

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp.ValidationError(validationErrs))
			return
		}

		id, err := urlSaver.SaveURL(req.URL, req.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL), slog.String("alias", req.Alias))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("url already exists")
			return
		}
		if err != nil {
			log.Error("Error saving url", slog.String("err", err.Error()))

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(fmt.Sprintf("Error saving Alias: %s must be unique", req.Alias))
			return
		}

		log.Info("Saved", slog.Int64("id", id), slog.String("alias", req.Alias))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(saved{
			Message: "Saved",
			Alias:   req.Alias,
			Url:     req.URL,
			IDRows:  id,
			Status:  "OK",
		})
	}
}

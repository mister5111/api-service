package example

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type req struct {
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

type alias struct {
	Alias string `json:"alias"`
}

func Example(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Example"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		exampleResponse := req{
			Alias: "shevchenko",
			URL:   "https://shevchenko.cc",
		}

		exampleDel := alias{
			Alias: "shevchenko",
		}

		exampleAll := alias{
			Alias: "all",
		}

		exampleOneRows := alias{
			Alias: "shevchenko",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Example save POST request: http://localhost:8080/save")
		json.NewEncoder(w).Encode(exampleResponse)

		json.NewEncoder(w).Encode("Example delete POST request: http://localhost:8080/del")
		json.NewEncoder(w).Encode(exampleDel)

		json.NewEncoder(w).Encode("Example show all rows GET request: http://localhost:8080/all")
		json.NewEncoder(w).Encode(exampleAll)

		json.NewEncoder(w).Encode("Example show one rows GET request: http://localhost:8080/all")
		json.NewEncoder(w).Encode(exampleOneRows)

		log.Info("Request Example", slog.String("status", "OK"))
	}
}

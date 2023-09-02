package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"
)

type ChannType int8

const (
	BI_CHANN  ChannType = 0
	TRI_CHANN ChannType = 1
)

func GetQueryHandler(c proxy.QueryClient, ct ChannType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}
		req := &proxy.QueryRequest{}
		err = json.Unmarshal(body, req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
			return
		}
		var resp *proxy.QueryResult
		switch ct {
		case BI_CHANN:
			resp, err = c.QueryBiChann(r.Context(), req)
		case TRI_CHANN:
			resp, err = c.QueryTriChann(r.Context(), req)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "while query: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

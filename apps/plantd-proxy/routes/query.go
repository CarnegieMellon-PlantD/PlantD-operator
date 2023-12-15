package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/proxy"
)

// queryRawHandler returns an HTTP handler function that receives requests for raw data and responds data of
// corresponding type.
func queryRawHandler(qa *proxy.QueryAgent, st proxy.SourceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}

		var resp any
		switch st {
		case proxy.Redis:
			req := &proxy.RedisRequest{}
			err = json.Unmarshal(body, req)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
				return
			}
			resp, err = qa.RedisGet(r.Context(), req)
		}
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while querying: " + err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// queryHandler returns an HTTP handler function that receives requests for bi-channel or tri-channel data and responds
// data of corresponding type.
func queryHandler(qa *proxy.QueryAgent, ct proxy.ChanType, st proxy.SourceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while reading request body: " + err.Error()})
			return
		}

		var resp any
		switch st {
		case proxy.Prometheus:
			req := &proxy.PromRequest{}
			err = json.Unmarshal(body, req)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
				return
			}
			switch ct {
			case proxy.BiChan:
				resp, err = qa.PromQuery(r.Context(), req)
			case proxy.TriChan:
				resp, err = qa.PromQueryRange(r.Context(), req)
			}

		case proxy.RedisTimeSeries:
			req := &proxy.RedisTSRequest{}
			err = json.Unmarshal(body, req)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while unmarshalling request body: " + err.Error()})
				return
			}
			switch ct {
			case proxy.BiChan:
				resp, err = qa.RedisTSMultiGet(r.Context(), req)
			case proxy.TriChan:
				resp, err = qa.RedisTSMultiRange(r.Context(), req)
			}
		}
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proxy.ErrorResponse{Message: "while querying: " + err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

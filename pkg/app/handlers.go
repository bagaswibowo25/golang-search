package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/julienschmidt/httprouter"
)

func IngestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req types.IngestRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, logEntry := range req.Logs {
		location, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(location)
		logEntry.Timestamp = now.Format(time.RFC3339) + "Z"
		ingestionChannel <- logEntry
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}

func SearchHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	hits := searchIndexedLogs(query)

	json.NewEncoder(w).Encode(types.SearchResponse{Hits: hits})
}

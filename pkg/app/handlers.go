package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/blevesearch/bleve"
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

func IngestDocHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req types.IngestRequest
	var docEntry types.DocIngest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	indexPath := ps.ByName("idx")

	for _, docEntry.Logs = range req.Logs {
		location, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(location)
		docEntry.Logs.Timestamp = now.Format(time.RFC3339) + "Z"
		docEntry.Idx = indexPath

		ingestDocChannel <- docEntry
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

func SearchDocHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	idxPath := ps.ByName("idx")
	hits := searchDoc(idxPath, query)

	json.NewEncoder(w).Encode(types.SearchResponse{Hits: hits})
}

func CreateIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var idx types.IndexTemplate
	var err error

	idx.Index, err = bleve.New(ps.ByName("idx"), bleve.NewIndexMapping())
	if err != nil {
		http.Error(w, "Unable to create new index", http.StatusBadRequest)
	}
	idx.Index.Close()

	var responseSuccess = make(map[string]string)

	responseSuccess["status"] = "created"
	responseSuccess["indexName"] = idx.Index.Name()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(responseSuccess)

}

package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bagaswibowo25/golang-search/pkg/pubsub"
	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/julienschmidt/httprouter"
)

type JetStream pubsub.JetStream

func (ids *IndexesMetadata) CreateIndexesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	indexesName := ps.ByName("ids")
	if indexesName == "" {
		log.Fatal("")
		http.Error(w, "Invalid URL param!", http.StatusBadRequest)
	}

	ids.createNewIndexes(indexesName)
	response := map[string]string{
		"Status": "Created Successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ids *IndexesMetadata) PublishLogsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var logs types.IngestRequest

	err := json.NewDecoder(r.Body).Decode(&logs)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	logData, err := json.Marshal(&logs.Logs)
	if err != nil {
		log.Printf("Invalid JSON format for logs!")
	}
	pubsub.PublishMessage(string(logData), ids.JServer)

	response := map[string]string{
		"Status": "Success to publish logs",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ids *IndexesMetadata) SearchDocsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	hits := searchIndexedLogs(ids.Alias, query)

	if hits == nil {
		http.Error(w, "Can't find the document", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(types.SearchResponse{Hits: hits})
}

package app

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/blevesearch/bleve"
	"github.com/julienschmidt/httprouter"
)

func (ids *Indices) IngestDocHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	var req types.IngestRequest
	var idxPath bleve.Index

	if len(ids.Index) == 0 {
		log.Printf("There is no indices exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, idx := range ids.Index {
		if idx.Name == ps.ByName("idx") && idx.Active {
			idxPath = idx.Path
			log.Printf("Index is found")
			break
		}

		log.Printf("No matching index")
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, logEntry := range req.Logs {
		logEntry.Index = idxPath
		location, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(location)
		logEntry.Timestamp = now.Format(time.RFC3339) + "Z"
		ingestionChannel <- logEntry
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}

func (ids *Indices) SearchDocHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var idxPath bleve.Index

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	if len(ids.Index) == 0 {
		http.Error(w, "There is no indices exists", http.StatusNotFound)
		return
	}

	for _, idx := range ids.Index {
		if idx.Name == ps.ByName("idx") && idx.Active {
			idxPath = idx.Path
		}
	}

	hits := searchIndexedLogs(idxPath, query)

	if hits == nil {
		http.Error(w, "Can't find the document", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(types.SearchResponse{Hits: hits})
}

func (ids *Indices) CreateIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var idx types.Index
	var err error

	idx.Path, err = bleve.New(ps.ByName("idx"), bleve.NewIndexMapping())
	if err != nil {
		http.Error(w, "Unable to create new index", http.StatusBadRequest)
		return
	}
	idx.Name = ps.ByName("idx")
	idx.Active = true
	ids.Index = append(ids.Index, idx)

	w.WriteHeader(http.StatusAccepted)
	response := map[string]string{
		"Status": "Created",
	}
	json.NewEncoder(w).Encode(response)

}

func (ids *Indices) CheckIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if len(ids.Index) == 0 {
		log.Printf("There is no indices exists")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, idx := range ids.Index {
		if idx.Name == ps.ByName("idx") && idx.Active {
			response := map[string]string{
				"Status": "Index is open",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (ids *Indices) UpdateIndexStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var index types.Index
	var err error
	query := r.URL.Query().Get("index")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	idxPath := ps.ByName("idx")

	if query == "open" {
		index.Path, err = bleve.Open(idxPath)
		if err != nil {
			log.Printf("Can't open the index")
			w.WriteHeader(http.StatusBadRequest)
		}
		index.Active = true
		index.Name = idxPath
	}

	if query == "close" {
		index.Path.Close()
		index.Active = false
		index.Name = idxPath
	}

	ids.Index = append(ids.Index, index)
	w.WriteHeader(http.StatusOK)
}

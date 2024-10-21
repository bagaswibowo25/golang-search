package app

import (
	"log"
	"net/http"

	"github.com/bagaswibowo25/golang-search/pkg/types"

	"github.com/blevesearch/bleve"
	"github.com/julienschmidt/httprouter"
)

var (
	ingestionChannel chan types.LogEntry
)

type Indices struct {
	Index []types.Index
}

func OpenIndex(idxPath string) bleve.Index {
	index, err := bleve.Open(idxPath)
	if err != nil {
		log.Printf("Can't open the index")
	}
	return index
}

func StartServer(port string) error {
	var indices Indices

	workers := 5

	ingestionChannel = make(chan types.LogEntry, 100)
	for i := 0; i < workers; i++ {
		go startIngestionWorker(i)
	}

	router := httprouter.New()
	// New implementation
	router.POST("/api/v1/index/:idx", indices.CreateIndexHandler)
	router.GET("/api/v1/index/:idx", indices.CheckIndexHandler)
	router.PUT("/api/v1/index/:idx", indices.UpdateIndexStatus)
	router.POST("/api/v1/log/:idx", indices.IngestDocHandler)
	router.GET("/api/v1/log/:idx", indices.SearchDocHandler)

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}

package app

import (
	"log"
	"net/http"

	"github.com/bagaswibowo25/golang-search/pkg/types"

	"github.com/blevesearch/bleve"
	"github.com/julienschmidt/httprouter"
)

var (
	bleveIndex       bleve.Index
	ingestionChannel chan types.LogEntry
	ingestDocChannel chan types.DocIngest
)

func StartServer(port string) error {
	workers := 5

	var err error
	bleveIndex, err = bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return err
	}

	ingestionChannel = make(chan types.LogEntry, 100)
	for i := 0; i < workers; i++ {
		go startIngestionWorker(i)
	}

	ingestDocChannel = make(chan types.DocIngest)
	for i := 0; i < workers; i++ {
		go startIngestDocWorker(i)
	}

	router := httprouter.New()
	router.POST("/api/v1/ingest", IngestHandler)
	router.GET("/api/v1/search", SearchHandler)

	// New implementation
	router.POST("/api/v1/index/:idx", CreateIndexHandler)
	router.POST("/api/v1/doc/:idx", IngestDocHandler)
	router.GET("/api/v1/doc/:idx", SearchDocHandler)

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}

package app

import (
	"log"
	"net/http"

	"github.com/bagaswibowo25/golang-search/pkg/types"

	"github.com/julienschmidt/httprouter"
)

var (
	ingestionChannel chan types.LogEntry
)

type Indices struct {
	Index []types.Index
}

func StartServer(port string) error {
	var indices Indices

	var logging IndexesMetadata
	logging.startIndexes()

	workers := 5

	ingestionChannel = make(chan types.LogEntry, 100)
	for i := 0; i < workers; i++ {
		go startIngestionWorker(i)
	}

	router := httprouter.New()
	router.POST("/api/v1/index/:idx", indices.CreateIndexHandler)
	router.GET("/api/v1/index/:idx", indices.CheckIndexHandler)
	router.PUT("/api/v1/index/:idx", indices.UpdateIndexStatus)
	router.POST("/api/v1/log/:idx", indices.IngestDocHandler)
	router.GET("/api/v1/log/:idx", indices.SearchDocHandler)

	// New Implementation
	router.POST("/api/v1/indices/:ids", CreateIndexesHandler)
	router.POST("/api/v1/docs", logging.IngestDocsHandler)
	router.GET("/api/v1/docs", logging.SearchDocsHandler)

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}

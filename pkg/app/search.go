package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/bagaswibowo25/golang-search/pkg/types"

	"github.com/blevesearch/bleve"
)

func startIngestionWorker(workerID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	defer wg.Done()

	log.Printf("Worker %d started", workerID)
	for logEntry := range ingestionChannel {
		log.Printf("Worker %d processing log: %v", workerID, logEntry)

		err := bleveIndex.Index(logEntry.ID, logEntry)
		if err != nil {
			log.Printf("Worker %d encountered error indexing log: %v", workerID, err)
		}
	}
	log.Printf("Worker %d exiting", workerID)
}

func searchIndexedLogs(query string) []types.LogEntry {
	var hits []types.LogEntry

	search := bleve.NewMatchQuery(query)
	searchRequest := bleve.NewSearchRequest(search)

	searchResult, err := bleveIndex.Search(searchRequest)
	if err != nil {
		log.Printf("Error during search: %v", err)
		fmt.Println(hits)
		return hits
	}

	for _, hit := range searchResult.Hits {
		doc, err := bleveIndex.Document(hit.ID)
		if err != nil {
			log.Printf("Error retrieving document for ID %s: %v", hit.ID, err)
			continue
		}

		var logEntry types.LogEntry
		for _, field := range doc.Fields {
			switch field.Name() {
			case "id":
				logEntry.ID = string(field.Value())
			case "timestamp":
				logEntry.Timestamp = string(field.Value())
			case "message":
				logEntry.Message = string(field.Value())
			case "date":
				logEntry.Date = string(field.Value())
			}
		}

		hits = append(hits, logEntry)
	}

	return hits
}

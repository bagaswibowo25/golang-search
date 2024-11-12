package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/blevesearch/bleve"
)

type IndexMetadata struct {
	Name string `json:"name"`
	Open bool   `json:"open"`
}

type IndexesMetadata struct {
	Alias    bleve.IndexAlias
	Metadata []IndexMetadata `json:"metadata"`
}

func searchIndexedLogs(idx bleve.Index, query string) []types.LogEntry {
	var hits []types.LogEntry

	search := bleve.NewMatchQuery(query)
	searchRequest := bleve.NewSearchRequest(search)

	searchResult, err := idx.Search(searchRequest)
	if err != nil {
		log.Printf("Error during search: %v", err)
		fmt.Println(hits)
		return hits
	}

	for _, hit := range searchResult.Hits {
		doc, err := idx.Document(hit.ID)
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
			}
		}

		hits = append(hits, logEntry)
	}

	return hits
}

func (lw *loggingWorkers) ingestLogs(logsMsg []byte) {
	var logs []types.LogEntry
	err := json.Unmarshal(logsMsg, &logs)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, logEntry := range logs {
		logEntry.Index = lw.indices.Alias
		location, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(location)
		logEntry.Timestamp = now.Format(time.RFC3339) + "Z"

		lw.workerChan <- logEntry
	}
}

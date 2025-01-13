package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bagaswibowo25/golang-search/pkg/pubsub"
	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/blevesearch/bleve"
)

type IndexesMetadata struct {
	types.IndexesMetadata
	JServer *pubsub.JetStream
}

func (app *IndexesMetadata) createNewIndexes(baseIndexName string) error {
	dateSuffix := time.Now().Format("2006-01-02")
	newIndexName := fmt.Sprintf("%s-%s", baseIndexName, dateSuffix)
	metadataFile := fmt.Sprintf("metadata/%s.json", baseIndexName)

	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		index, err := bleve.New(newIndexName, bleve.NewIndexMapping())
		if err != nil {
			return fmt.Errorf("failed to create index %s: %v", newIndexName, err)
		}

		app.Alias = bleve.NewIndexAlias(index)
		app.Metadata = []types.IndexMetadata{
			{Name: newIndexName, Open: true},
		}

		err = saveAliasMetadata(metadataFile, *app)
		if err != nil {
			return fmt.Errorf("failed to save alias metadata: %v", err)
		}

		fmt.Println("New alias and single index created successfully. Metadata saved.")
		return nil
	}
	fmt.Println("Metadata file exists. Alias will be loaded.")
	return nil
}

func saveAliasMetadata(filename string, indexesMetadata IndexesMetadata) error {
	jsonData, err := json.Marshal(map[string][]types.IndexMetadata{"metadata": indexesMetadata.Metadata})
	if err != nil {
		return fmt.Errorf("failed to serialize metadata to JSON: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %v", err)
	}

	return nil
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

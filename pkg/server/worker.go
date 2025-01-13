package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bagaswibowo25/golang-search/pkg/app"
	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/blevesearch/bleve"
)

type loggingWorkers struct {
	WorkerChan chan types.LogEntry
	Indices    app.IndexesMetadata
}

func (w *loggingWorkers) startLoggingWorker(workerID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	defer wg.Done()

	log.Printf("Logging worker %d started", workerID)
	for logEntry := range w.WorkerChan {
		log.Printf("Worker %d processing log: %v", workerID, logEntry)

		err := logEntry.Index.Index(logEntry.ID, logEntry)
		if err != nil {
			log.Printf("Worker %d encountered error indexing log: %v", workerID, err)
		}
	}
	log.Printf("Logging worker %d exiting", workerID)
}

func (lw *loggingWorkers) startIndexes() error {
	var metadataFile string

	files, err := os.ReadDir("metadata")
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		metadataFile = fmt.Sprintf("metadata/%s", file.Name())
	}

	file, err := os.Open(metadataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read metadata file: %v", err)
	}
	fmt.Println(string(byteValue))

	err = json.Unmarshal(byteValue, &lw)
	if err != nil {
		log.Fatalf("failed to parse metadata file: %v", err)
	}

	alias := bleve.NewIndexAlias()

	for _, indexInfo := range lw.Indices.Metadata {
		if indexInfo.Open {
			index, err := bleve.Open(indexInfo.Name)
			if err != nil {
				log.Printf("Failed to open index %s: %v", indexInfo.Name, err)
				continue
			}
			alias.Add(index)
		} else {
			fmt.Printf("Skipping closed index: %s\n", indexInfo.Name)
		}
	}
	lw.Indices.Alias = alias

	return nil
}

func (lw *loggingWorkers) ingestLogs(logs []types.LogEntry) {
	for _, logEntry := range logs {
		logEntry.Index = lw.Indices.Alias
		location, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(location)
		logEntry.Timestamp = now.Format(time.RFC3339) + "Z"

		lw.WorkerChan <- logEntry
	}
}

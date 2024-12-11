package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/nats-io/nats.go"
)

func (app *IndexesMetadata) startIndexes() error {
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

	err = json.Unmarshal(byteValue, &app)
	if err != nil {
		log.Fatalf("failed to parse metadata file: %v", err)
	}

	alias := bleve.NewIndexAlias()

	for _, indexInfo := range app.Metadata {
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
	app.Alias = alias

	return nil
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
		app.Metadata = []IndexMetadata{
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
	jsonData, err := json.Marshal(map[string][]IndexMetadata{"metadata": indexesMetadata.Metadata})
	if err != nil {
		return fmt.Errorf("failed to serialize metadata to JSON: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %v", err)
	}

	return nil
}

func (w *loggingWorkers) loggingWorker(workerID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	defer wg.Done()

	log.Printf("Logging worker %d started", workerID)
	for logEntry := range w.workerChan {
		log.Printf("Worker %d processing log: %v", workerID, logEntry)

		err := logEntry.Index.Index(logEntry.ID, logEntry)
		if err != nil {
			log.Printf("Worker %d encountered error indexing log: %v", workerID, err)
		}
	}
	log.Printf("Logging worker %d exiting", workerID)
}

func (n *jetStream) publishMessage(message string) {
	_, err := n.js.Publish(n.c.Subject, []byte(message))
	if err != nil {
		log.Printf("Error publishing message: %v", err)
	} else {
		fmt.Printf("[%s] Published message: %s\n", n.c.ConsumerId, message)
	}
}

func (n *jetStream) subscribeMessage(cb nats.MsgHandler) {
	_, err := n.js.Subscribe(n.c.Subject, cb, nats.Durable(n.c.ConsumerId), nats.ManualAck(), nats.SkipConsumerLookup())

	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
		return
	}
}

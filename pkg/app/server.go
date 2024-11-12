package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/nats-io/nats.go"

	"github.com/julienschmidt/httprouter"
)

var (
	consumerID string
)

type jetStream struct {
	natsURL string
	subject string
	stream  string
	js      nats.JetStreamContext
	nc      *nats.Conn
}

type loggingWorkers struct {
	workerChan chan types.LogEntry
	indices    IndexesMetadata
}

func StartServer(port string) error {
	var wg sync.WaitGroup
	var indices IndexesMetadata

	err := indices.startIndexes()
	if err != nil {
		log.Printf("No existing indices found! Please create new index")
	}

	jServer := jetStream{
		natsURL: "http://localhost:4222",
		subject: "loggerSubject",
		stream:  "loggerStream",
	}

	wg.Add(1)
	jServer.StartJS(&wg)
	defer jServer.CloseJS()

	workers := 5
	lw := &loggingWorkers{
		workerChan: make(chan types.LogEntry, 100),
		indices:    indices,
	}
	for i := 0; i < workers; i++ {
		go lw.loggingWorker(i)
	}

	jServer.subscribeMessage(func(msg *nats.Msg) {
		fmt.Printf("[%s] Received message: %s\n", consumerID, string(msg.Data))
		lw.ingestLogs(msg.Data)
		msg.Ack()
	})

	router := httprouter.New()
	router.POST("/api/v1/indices/:ids", lw.indices.CreateIndexesHandler)
	router.POST("/api/v1/logs", jServer.publishLogsHandler)
	router.GET("/api/v1/logs", lw.indices.SearchDocsHandler)

	wg.Add(1)
	go StartHTTPServer(port, &wg, router)

	wg.Wait()

	return nil
}

func StartHTTPServer(port string, wg *sync.WaitGroup, router *httprouter.Router) {
	defer wg.Done()

	log.Printf("Starting HTTP server on port %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

func (jServer *jetStream) StartJS(wg *sync.WaitGroup) {
	var err error
	defer wg.Done()

	flag.StringVar(&consumerID, "consumer", os.Getenv("CONSUMER_ID"), "Unique consumer ID for the server")
	flag.Parse()
	if consumerID == "" {
		log.Fatal("Consumer ID must be provided as an argument or environment variable")
	}

	jServer.nc, err = nats.Connect(jServer.natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	jServer.js, err = jServer.nc.JetStream()
	if err != nil {
		log.Fatalf("Error enabling JetStream: %v", err)
	}

	_, err = jServer.js.AddStream(&nats.StreamConfig{
		Name:     jServer.stream,
		Subjects: []string{jServer.subject},
	})
	if err != nil {
		log.Printf("Stream may already exist: %v", err)
	}
}

func (jServer *jetStream) CloseJS() {
	jServer.nc.Drain()
}

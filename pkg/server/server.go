package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bagaswibowo25/golang-search/pkg/app"
	"github.com/bagaswibowo25/golang-search/pkg/config"
	"github.com/bagaswibowo25/golang-search/pkg/pubsub"
	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/nats-io/nats.go"

	"github.com/julienschmidt/httprouter"
)

func StartServer(httpConf *config.HttpConfig, natsConf *config.NatsConfig) error {
	var wg sync.WaitGroup
	var indices app.IndexesMetadata

	jServer := &pubsub.JetStream{
		Conf: *natsConf,
	}

	wg.Add(1)
	startJS(&wg, natsConf.ConsumerId, jServer)
	defer closeJS(jServer)
	indices.JServer = jServer

	workers := 5
	lw := &loggingWorkers{
		WorkerChan: make(chan types.LogEntry, 100),
		Indices:    indices,
	}
	for i := 0; i < workers; i++ {
		go lw.startLoggingWorker(i)
	}

	err := lw.startIndexes()
	if err != nil {
		log.Printf("No existing indices found! Please create new index")
	}

	pubsub.SubscribeMessage(func(msg *nats.Msg) {
		fmt.Printf("[%s] Received message: %s\n", natsConf.ConsumerId, string(msg.Data))
		var logs []types.LogEntry
		err := json.Unmarshal(msg.Data, &logs)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		lw.ingestLogs(logs)
		msg.Ack()
	}, jServer)

	router := httprouter.New()
	router.POST("/api/v1/indices/:ids", lw.Indices.CreateIndexesHandler)
	router.POST("/api/v1/logs", lw.Indices.PublishLogsHandler)
	router.GET("/api/v1/logs", lw.Indices.SearchDocsHandler)

	wg.Add(1)
	go StartHTTPServer(httpConf.ListenPort, &wg, router)

	wg.Wait()

	return err
}

func StartHTTPServer(port string, wg *sync.WaitGroup, router *httprouter.Router) {
	defer wg.Done()

	log.Printf("Starting HTTP server on port %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

func startJS(wg *sync.WaitGroup, consumerId string, jServer *pubsub.JetStream) {
	var err error
	defer wg.Done()

	if consumerId == "" {
		log.Fatal("Consumer ID must be provided as an argument or environment variable")
	}

	jServer.Conn, err = nats.Connect(jServer.Conf.NatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	jServer.StreamCtx, err = jServer.Conn.JetStream()
	if err != nil {
		log.Fatalf("Error enabling JetStream: %v", err)
	}

	_, err = jServer.StreamCtx.AddStream(&nats.StreamConfig{
		Name:     jServer.Conf.Stream,
		Subjects: []string{jServer.Conf.Subject},
	})
	if err != nil {
		log.Printf("Stream may already exist: %v", err)
	}
}

func closeJS(jServer *pubsub.JetStream) {
	jServer.Conn.Drain()
}

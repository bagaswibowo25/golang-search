package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bagaswibowo25/golang-search/pkg/config"
	"github.com/bagaswibowo25/golang-search/pkg/types"
	"github.com/nats-io/nats.go"

	"github.com/julienschmidt/httprouter"
)

type jetStream struct {
	streamCtx nats.JetStreamContext
	conn      *nats.Conn
	conf      config.NatsConfig
}

type loggingWorkers struct {
	workerChan chan types.LogEntry
	indices    IndexesMetadata
}

func StartServer(httpConf *config.HttpConfig, natsConf *config.NatsConfig) error {
	var wg sync.WaitGroup
	var indices IndexesMetadata

	err := indices.startIndexes()
	if err != nil {
		log.Printf("No existing indices found! Please create new index")
	}

	jServer := jetStream{
		conf: *natsConf,
	}

	wg.Add(1)
	jServer.StartJS(&wg, natsConf.ConsumerId)
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
		fmt.Printf("[%s] Received message: %s\n", natsConf.ConsumerId, string(msg.Data))
		var logs []types.LogEntry
		err := json.Unmarshal(msg.Data, &logs)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		lw.ingestLogs(logs)
		msg.Ack()
	})

	router := httprouter.New()
	router.POST("/api/v1/indices/:ids", lw.indices.CreateIndexesHandler)
	router.POST("/api/v1/logs", jServer.publishLogsHandler)
	router.GET("/api/v1/logs", lw.indices.SearchDocsHandler)

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

func (jServer *jetStream) StartJS(wg *sync.WaitGroup, consumerId string) {
	var err error
	defer wg.Done()

	if consumerId == "" {
		log.Fatal("Consumer ID must be provided as an argument or environment variable")
	}

	jServer.conn, err = nats.Connect(jServer.conf.NatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	jServer.streamCtx, err = jServer.conn.JetStream()
	if err != nil {
		log.Fatalf("Error enabling JetStream: %v", err)
	}

	_, err = jServer.streamCtx.AddStream(&nats.StreamConfig{
		Name:     jServer.conf.Stream,
		Subjects: []string{jServer.conf.Subject},
	})
	if err != nil {
		log.Printf("Stream may already exist: %v", err)
	}
}

func (jServer *jetStream) CloseJS() {
	jServer.conn.Drain()
}

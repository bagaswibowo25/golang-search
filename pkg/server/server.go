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

type ServerConfig struct {
	HttpConf *config.HttpConfig
	NatsConf *config.NatsConfig
	wg       sync.WaitGroup
	jServer  *pubsub.JetStream
}

func (srv *ServerConfig) StartServer() error {
	var wg sync.WaitGroup
	var indices app.IndexesMetadata

	srv.jServer = &pubsub.JetStream{
		Conf: *srv.NatsConf,
	}

	srv.wg.Add(1)
	srv.startJS()
	defer closeJS(srv.jServer)
	indices.JServer = srv.jServer

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

	srv.jServer.SubscribeMessage(func(msg *nats.Msg) {
		fmt.Printf("[%s] Received message: %s\n", srv.NatsConf.ConsumerId, string(msg.Data))
		var logs []types.LogEntry
		err := json.Unmarshal(msg.Data, &logs)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		lw.ingestLogs(logs)
		msg.Ack()
	})

	srv.HttpConf.Router = httprouter.New()
	srv.HttpConf.Router.POST("/api/v1/indices/:ids", lw.Indices.CreateIndexesHandler)
	srv.HttpConf.Router.POST("/api/v1/logs", lw.Indices.PublishLogsHandler)
	srv.HttpConf.Router.GET("/api/v1/logs", lw.Indices.SearchDocsHandler)

	srv.wg.Add(1)
	go srv.startHTTPServer()

	wg.Wait()

	return err
}

func (srv *ServerConfig) startHTTPServer() {
	defer srv.wg.Done()

	log.Printf("Starting HTTP server on port %s\n", srv.HttpConf.ListenPort)
	if err := http.ListenAndServe(srv.HttpConf.ListenPort, srv.HttpConf.Router); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

func (srv *ServerConfig) startJS() {
	var err error
	defer srv.wg.Done()

	if srv.NatsConf.ConsumerId == "" {
		log.Fatal("Consumer ID must be provided as an argument or environment variable")
	}

	srv.jServer.Conn, err = nats.Connect(srv.jServer.Conf.NatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	srv.jServer.StreamCtx, err = srv.jServer.Conn.JetStream()
	if err != nil {
		log.Fatalf("Error enabling JetStream: %v", err)
	}

	_, err = srv.jServer.StreamCtx.AddStream(&nats.StreamConfig{
		Name:     srv.jServer.Conf.Stream,
		Subjects: []string{srv.jServer.Conf.Subject},
	})
	if err != nil {
		log.Printf("Stream may already exist: %v", err)
	}
}

func closeJS(jServer *pubsub.JetStream) {
	jServer.Conn.Drain()
}

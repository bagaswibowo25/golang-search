package types

import "github.com/blevesearch/bleve"

type LogEntry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type IngestRequest struct {
	Logs []LogEntry `json:"logs"`
}

type DocIngest struct {
	Idx  string
	Logs LogEntry
}

type SearchResponse struct {
	Hits []LogEntry `json:"hits"`
}

type IndexTemplate struct {
	Index bleve.Index
}

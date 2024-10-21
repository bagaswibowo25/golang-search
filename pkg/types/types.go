package types

import "github.com/blevesearch/bleve"

type LogEntry struct {
	Index     bleve.Index
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type IngestRequest struct {
	Logs []LogEntry `json:"logs"`
}

type SearchResponse struct {
	Hits []LogEntry `json:"hits"`
}

type Index struct {
	Path   bleve.Index
	Active bool
	Name   string
}

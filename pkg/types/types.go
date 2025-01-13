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



type IndexMetadata struct {
	Name string `json:"name"`
	Open bool   `json:"open"`
}

type IndexesMetadata struct {
	Alias    bleve.IndexAlias
	Metadata []IndexMetadata `json:"metadata"`
}

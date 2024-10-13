package types

type LogEntry struct {
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

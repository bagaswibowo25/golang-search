package types

type LogEntry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Date      string `json:"date"`
}

type IngestRequest struct {
	Logs []LogEntry `json:"logs"`
}

type SearchResponse struct {
	Hits []LogEntry `json:"hits"`
}

#!/bin/bash

# Create Index
curl -X POST http://localhost:8080/api/v1/index/myidx

# Ingest doc to index
curl -X POST http://localhost:8080/api/v1/doc/myidx -d '{"logs":[{"id":"1","message":"My Log to myidx"}]}' -H "Content-Type: application/json"

# Query docs from index
curl "http://localhost:8080/api/v1/doc/myidx?query=My"
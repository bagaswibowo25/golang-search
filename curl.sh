#!/bin/bash

# Create Indices
curl -X POST http://localhost:8080/api/v1/indices/logging

# Ingest docs
curl -X POST http://localhost:8080/api/v1/logs -d '{"logs":[{"id":"1","message":"Hello there!"},{"id":"2","message":"Hello there2!"}]}' -H "Content-Type: application/json"

# Ingest docss
curl -X GET "http://localhost:8080/api/v1/logs?query=here"
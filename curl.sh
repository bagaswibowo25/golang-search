#!/bin/bash

# Create Index
curl -X POST http://localhost:8080/api/v1/index/myidx

# Open Index
curl -X PUT http://localhost:8080/api/v1/index/myidx?index=open

# Ingest doc to index
curl -X POST http://localhost:8080/api/v1/log/myidx -d '{"logs":[{"id":"1","message":"My Log to myidx"}]}' -H "Content-Type: application/json"

# Check Index status
curl -X GET http://localhost:8080/api/v1/index/myidx 

# Query docs from index
curl -X GET "http://localhost:8080/api/v1/log/myidx?query=My"
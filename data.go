package main

import (
	"time"
)

type Source struct {
	Message      string    `json:"message"`
	Version      string    `json:"@version"`
	Timestamp    time.Time `json:"@timestamp"`
	Host         string    `json:"host"`
	Port         uint16    `json:"port"`
	Type         string    `json:"type"`
	LSRecTime    time.Time `json:"logstash-receiver-time"`
	LSIDXTime    time.Time `json:"logstash-indexer-time"`
	LogLevel     string    `json:"loglevel"`
	LogFlow      string    `json:"logflow"`
	IDXTimestamp time.Time `json:"indexer_timestamp"`
	IDXName      string    `json:"index_name"`
}

type Hits struct {
	IDX   string   `json:"_index"`
	Type  string   `json:"_type"`
	ID    string   `json:"_id"`
	Score uint64   `json:"_score"`
	Src   []Source `json:"_source"`
}

type HitsSummary struct {
	Total    uint64 `json:"total"`
	MAXScore uint64 `json:"max_score"`
	Hits     []Hits `json:"hits"`
}



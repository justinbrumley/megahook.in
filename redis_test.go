package main

import (
	"testing"
	"time"
)

func TestRedisAddRecord(t *testing.T) {
	err := initRedis()
	if err != nil {
		t.Errorf("Error initializing redis: %v\n", err)
	}

	request := &Request{
		Body: "Test Body. Please Ignore.",
	}

	response := &Response{
		Body:       "Test Response Body.",
		StatusCode: 200,
	}

	record := &Record{
		Request:   request,
		Response:  response,
		Timestamp: time.Now().Unix(),
	}

	err = addRecord("test-redis-key", record)
	if err != nil {
		t.Errorf("Error adding record to redis: %v\n", err)
	}
}

func TestRedisGetRecords(t *testing.T) {
	err := initRedis()
	if err != nil {
		t.Errorf("Error initializing redis: %v\n", err)
	}

	records, err := getRecords("test-redis-key")
	if err != nil {
		t.Errorf("Error getting records from redis: %v\n", err)
	}

	if len(records) < 1 {
		t.Errorf("Test record not found in redis\n")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"os"
	"strconv"
	"time"
)

// Record used to store request to webhook and response from client in redis.
type Record struct {
	Request   *Request  `json:"request,omitempty"`
	Response  *Response `json:"response,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

var rClient *redis.Client

var addr string = "localhost:6379"
var password string = ""
var db int = 0

func initRedis() error {
	if os.Getenv("REDIS_HOST") != "" {
		addr = os.Getenv("REDIS_HOST")
	}

	if os.Getenv("REDIS_PASSWORD") != "" {
		password = os.Getenv("REDIS_PASSWORD")
	}

	if os.Getenv("REDIS_DB") != "" {
		var err error
		db, err = strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			db = 0
		}
	}

	rClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

// Removes old records from list in redis by key
func purgeRecords(key string) error {
	max := time.Now().AddDate(0, 0, -1).Unix() // 24hr max age
	err := rClient.ZRemRangeByScore(key, "-inf", fmt.Sprintf("%d", max)).Err()
	return err
}

// Add Record to list in redis using current timestamp as score
func addRecord(key string, record *Record) error {
	err := purgeRecords(key)
	if err != nil {
		return err
	}

	member, err := json.Marshal(record)
	if err != nil {
		return err
	}

	z := &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: member,
	}

	err = rClient.ZAdd(key, z).Err()
	if err != nil {
		return err
	}

	// Reset expire on list
	err = rClient.Expire(key, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

// Fetch slice of Records from redis by key
func getRecords(key string) ([]Record, error) {
	err := purgeRecords(key)
	if err != nil {
		return nil, err
	}

	results, err := rClient.ZRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var records []Record
	for _, item := range results {
		record := &Record{}
		err = json.Unmarshal([]byte(item), &record)
		if err != nil {
			fmt.Printf("Failed to unmarshal item: %v\n", err)
			continue
		}

		records = append(records, *record)
	}

	return records, nil
}

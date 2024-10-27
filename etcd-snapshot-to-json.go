package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	bolt "go.etcd.io/bbolt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"github.com/golang/protobuf/proto"
)

// KeyValueJSON represents a key-value pair in JSON format.
type KeyValueJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: etcd-snapshot-to-json <snapshot-file>")
		os.Exit(0)
	}

	snapshotFile := os.Args[1]

	// Open the snapshot file with bbolt
	db, err := bolt.Open(snapshotFile, 0600, nil)
	if err != nil {
		log.Fatalf("Failed to open snapshot file: %v", err)
	}
	defer db.Close()

	// Slice to store key-value pairs as JSON objects
	var kvPairs []KeyValueJSON

	// Traverse the snapshot and decode each key-value pair
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			return b.ForEach(func(k, v []byte) error {
				// Decode the value as a protobuf KeyValue message
				kv := &mvccpb.KeyValue{}
				if err := proto.Unmarshal(v, kv); err != nil {
					log.Printf("Failed to decode value for key %s: %v", k, err)
					return nil
				}

				// Add decoded key-value to the list, converting value to JSON if possible
				kvJSON := KeyValueJSON{
					Key:   string(kv.Key),
					Value: interpretValue(kv.Value),
				}
				kvPairs = append(kvPairs, kvJSON)
				return nil
			})
		})
	})
	if err != nil {
		log.Fatalf("Failed to read snapshot: %v", err)
	}

	// Marshal the entire list of key-value pairs as JSON
	jsonOutput, err := json.MarshalIndent(kvPairs, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonOutput))
}

// interpretValue returns a JSON string for the value, attempting to decode it if possible
func interpretValue(value []byte) string {
	// Check if value is valid UTF-8 and printable
	if utf8.Valid(value) && isPrintable(string(value)) {
		return string(value)
	}

	// If not printable, return the value as Base64-encoded string
	return base64.StdEncoding.EncodeToString(value)
}

// isPrintable checks if a string contains only printable characters
func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

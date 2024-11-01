package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"

	bolt "go.etcd.io/bbolt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"github.com/golang/protobuf/proto"
)

// KeyValueJSON represents a key-value pair in JSON format
type KeyValueJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	// Command-line flags
	flag.Parse()

	// Get the snapshot file path from the first positional argument
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: etcd_snapshot_to_json <snapshot-file>")
		os.Exit(1)
	}
	snapshotFile := flag.Arg(0)

	// Open the snapshot file with bbolt
	db, err := bolt.Open(snapshotFile, 0600, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to open snapshot file.")
		os.Exit(1)
	}
	defer db.Close()

	// Slice to store key-value pairs as JSON objects
	var kvPairs []KeyValueJSON

	// Traverse the snapshot and decode each key-value pair
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			return b.ForEach(func(k, v []byte) error {
				// Decode using protobuf's mvccpb.KeyValue structure
				kv := &mvccpb.KeyValue{}
				if err := proto.Unmarshal(v, kv); err != nil {
					// Skip entries that fail to decode without logging to stdout
					return nil
				}

				// Decode the key as a string
				keyStr := decodeOrReturnBase64(kv.Key)

				// Process the value, conditionally decoding if valid Base64
				valueStr := decodeOrReturnBase64(kv.Value)

				// Append to the JSON output slice
				kvPairs = append(kvPairs, KeyValueJSON{
					Key:   keyStr,
					Value: valueStr,
				})
				return nil
			})
		})
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to read snapshot.")
		os.Exit(1)
	}

	// Marshal the entire list of key-value pairs as JSON and print to stdout
	jsonOutput, err := json.MarshalIndent(kvPairs, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to marshal JSON.")
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}

// decodeOrReturnBase64 attempts to decode a string if it's valid Base64; otherwise, returns it as-is
func decodeOrReturnBase64(data []byte) string {
	// Check if the data is valid UTF-8; if it is, return it as-is
	if utf8.Valid(data) {
		return string(data)
	}

	// Attempt Base64 decoding, return as-is if not successful
	if decoded, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		return string(decoded)
	}
	return string(data)
}

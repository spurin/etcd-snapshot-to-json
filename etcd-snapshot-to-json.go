package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "unicode/utf8"

    bolt "go.etcd.io/bbolt"
    "go.etcd.io/etcd/api/v3/mvccpb"
    "github.com/golang/protobuf/proto"
    "github.com/spf13/cobra"
)

// KeyValueJSON represents a key-value pair in JSON format
type KeyValueJSON struct {
    Key            string `json:"key"`
    Value          string `json:"value"`
    CreateRevision int64  `json:"create_revision"`
    ModRevision    int64  `json:"mod_revision"`
    Version        int64  `json:"version"`
}

var (
    rootCmd    = &cobra.Command{
        Use:   "etcd_snapshot_to_json [snapshot-file]",
        Short: "Converts an ETCD snapshot file to a JSON representation of key-value pairs",
        Args:  cobra.ExactArgs(1),
        Run:   run,
    }
    keysFilter string
    latest     bool
)

func init() {
    rootCmd.PersistentFlags().StringVar(&keysFilter, "keys", "", "Comma-separated list of keys to include in the output")
    rootCmd.PersistentFlags().BoolVar(&latest, "latest", false, "Include only the latest version of each key")
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func run(cmd *cobra.Command, args []string) {
    snapshotFile := args[0]
    db, err := bolt.Open(snapshotFile, 0600, nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to open snapshot file: %v\n", err)
        os.Exit(1)
    }
    defer db.Close()

    var kvPairs []KeyValueJSON
    keysToInclude := make(map[string]bool)
    highestVersion := make(map[string]KeyValueJSON)

    if keysFilter != "" {
        for _, key := range strings.Split(keysFilter, ",") {
            keysToInclude[key] = true
        }
    }

    err = db.View(func(tx *bolt.Tx) error {
        return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
            return b.ForEach(func(k, v []byte) error {
                kv := &mvccpb.KeyValue{}
                if err := proto.Unmarshal(v, kv); err != nil {
                    return nil
                }
                keyStr := decodeOrReturnBase64(kv.Key)
                if _, ok := keysToInclude[keyStr]; ok || len(keysToInclude) == 0 {
                    valueStr := decodeOrReturnBase64(kv.Value)
                    keyValue := KeyValueJSON{
                        Key:            keyStr,
                        Value:          valueStr,
                        CreateRevision: kv.CreateRevision,
                        ModRevision:    kv.ModRevision,
                        Version:        kv.Version,
                    }

                    if latest {
                        // Only keep the highest version of each key
                        if existing, exists := highestVersion[keyStr]; !exists || existing.Version < kv.Version {
                            highestVersion[keyStr] = keyValue
                        }
                    } else {
                        // Append all versions if latest is not specified
                        kvPairs = append(kvPairs, keyValue)
                    }
                }
                return nil
            })
        })
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to read snapshot: %v\n", err)
        os.Exit(1)
    }

    // If latest flag is set, convert the map to a slice
    if latest {
        for _, kv := range highestVersion {
            kvPairs = append(kvPairs, kv)
        }
    }

    jsonOutput, err := json.MarshalIndent(kvPairs, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
        os.Exit(1)
    }
    fmt.Println(string(jsonOutput))
}

func decodeOrReturnBase64(data []byte) string {
    if utf8.Valid(data) {
        return string(data)
    }
    if decoded, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
        return string(decoded)
    }
    return string(data)
}

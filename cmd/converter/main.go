package main

import (
	"fmt"
)

type parquetTick struct {
	Timestamp int64   `parquet:"name=timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"` // Store as milliseconds since epoch
	Bid       float64 `parquet:"name=bid, type=DOUBLE"`
	Ask       float64 `parquet:"name=ask, type=DOUBLE"`
}

func main() {
	// Convert all csv files where the target does not exist
	err := convertMissingParquetFiles()
	if err != nil {
		fmt.Printf("❌ Conversion failed: %v\n", err)
	}
}

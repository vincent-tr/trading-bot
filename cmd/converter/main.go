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
	fmt.Println("🔄 Starting data conversion...")

	// Convert HistData ZIP files
	fmt.Println("\n📁 Processing HistData files...")
	err := convertMissingParquetFiles()
	if err != nil {
		fmt.Printf("❌ HistData conversion failed: %v\n", err)
	}

	// Convert Dukascopy CSV files
	fmt.Println("\n📁 Processing Dukascopy files...")
	err = convertDukascopyCsvFiles()
	if err != nil {
		fmt.Printf("❌ Dukascopy conversion failed: %v\n", err)
	}

	fmt.Println("\n✅ Conversion complete!")
}

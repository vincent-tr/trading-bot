package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const dukascopyPath = dataPath + "/dukascopy"

// Dukascopy CSV format:
// timestamp,askPrice,bidPrice
// Example: 1704146412108,1.10481,1.10427
// Timestamp is Unix milliseconds (GMT)

func loadDukascopyCsv(csvFile string) ([]parquetTick, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %v", csvFile, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	ticks := make([]parquetTick, 0)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %v", err)
		}
		if len(row) < 3 {
			return nil, fmt.Errorf("expected 3 columns in CSV row, got %d: %v", len(row), row)
		}

		// Parse timestamp (already in Unix milliseconds)
		timestamp, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp '%s': %v", row[0], err)
		}

		ask, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ask price '%s': %v", row[1], err)
		}

		bid, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bid price '%s': %v", row[2], err)
		}

		tick := parquetTick{
			Timestamp: timestamp,
			Bid:       bid,
			Ask:       ask,
		}
		ticks = append(ticks, tick)
	}

	return ticks, nil
}

func convertDukascopyCsvFiles() error {
	// Look for CSV files in dukascopy directory
	pattern := filepath.Join(dukascopyPath, "*.csv")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list CSV files: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("ℹ️  No Dukascopy CSV files found in brokers/backtesting/data/dukascopy/")
		return nil
	}

	for _, csvFile := range files {
		base := filepath.Base(csvFile)
		// Parse: eurusd-tick-2024-01-01-2024-01-31.csv → EURUSD_202401.parquet
		parts := strings.Split(strings.TrimSuffix(base, ".csv"), "-")
		if len(parts) < 4 {
			fmt.Printf("⚠️  Skipping file with unexpected name format: %s\n", base)
			continue
		}

		instrument := strings.ToUpper(parts[0]) // EURUSD
		year := parts[2]                        // 2024
		month := parts[3]                       // 01
		parquetName := fmt.Sprintf("%s_%s%s.parquet", instrument, year, month)
		parquetPath := filepath.Join(dukascopyPath, parquetName)

		if _, err := os.Stat(parquetPath); err == nil {
			fmt.Printf("✅ Parquet exists: %s (skipping)\n", parquetName)
			continue
		}

		fmt.Printf("📦 Converting: %s → %s\n", base, parquetName)

		ticks, err := loadDukascopyCsv(csvFile)
		if err != nil {
			return fmt.Errorf("failed to load CSV: %v", err)
		}

		if err := writeParquet(parquetPath, ticks); err != nil {
			return fmt.Errorf("failed to write parquet: %v", err)
		}

		// Delete source CSV file after successful conversion
		if err := os.Remove(csvFile); err != nil {
			fmt.Printf("⚠️  Warning: failed to delete source file %s: %v\n", base, err)
		} else {
			fmt.Printf("🗑️  Deleted source file: %s\n", base)
		}
	}

	return nil
}

# trading-bot

## Data Sources

### HistData

1. Download historical data from [HistData](https://www.histdata.com/)
   - Example: [EUR/USD Tick Data Quotes](https://www.histdata.com/download-free-forex-historical-data/?/ascii/tick-data-quotes/EURUSD)
2. Place ZIP files in `brokers/backtesting/data/histdata/`
3. Convert to Parquet format (see Data Conversion below)

### Dukascopy

**Automated Download (Recommended)**

Use the download script to automatically download all historical data:

```bash
# Download all data from 2009-12 to current month (default)
./download-dukascopy.sh

# Or specify a custom date range (start_year start_month end_year end_month)
./download-dukascopy.sh 2022 1 2024 12

# Run in background
nohup ./download-dukascopy.sh > download.log 2>&1 &
```

The script will:
- Skip files that already exist in `brokers/backtesting/data/dukascopy/`
- Clean up partial downloads before starting
- Move downloaded CSV files to the correct directory automatically

**Manual Download**

Download tick data manually using the [dukascopy-node](https://github.com/Leo4815162342/dukascopy-node) tool:
```bash
npx dukascopy-node -i eurusd -from 2024-01-01 -to 2024-01-31 -t tick -f csv
```
See available instruments: [FX Majors](https://github.com/Leo4815162342/dukascopy-node?tab=readme-ov-file#fx_majors)

Then place CSV files in `brokers/backtesting/data/dukascopy/` and convert to Parquet format (see Data Conversion below)

## Data Conversion

The converter processes both HistData ZIP files and Dukascopy CSV files, converting them to a unified Parquet format:

```bash
make convert
# or directly: go run ./cmd/converter
```

Output format:
- Timestamp: Unix milliseconds (UTC/GMT)
- Bid: Double
- Ask: Double


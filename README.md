# trading-bot

## Data Sources

### HistData

1. Download historical data from [HistData](https://www.histdata.com/)
   - Example: [EUR/USD Tick Data Quotes](https://www.histdata.com/download-free-forex-historical-data/?/ascii/tick-data-quotes/EURUSD)
2. Convert the downloaded ZIP files to Parquet format:
   ```bash
   go run cmd/converter/main.go
   ```

### Dukascopy

Download tick data using the [dukascopy-node](https://github.com/Leo4815162342/dukascopy-node) tool:

```bash
npx dukascopy-node -i eurusd -from 2024-01-01 -to 2024-01-31 -t tick -f csv
```

See available instruments: [FX Majors](https://github.com/Leo4815162342/dukascopy-node?tab=readme-ov-file#fx_majors)


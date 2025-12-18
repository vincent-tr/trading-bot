#!/bin/bash

# Download all Dukascopy EURUSD tick data
# Usage: ./download-dukascopy.sh [start_year] [start_month] [end_year] [end_month]
# Example: ./download-dukascopy.sh 2009 12 2025 6
# Or just run without args to download from 2009-12 to current date

set -e

# Directories
DOWNLOAD_DIR="download"
TARGET_DIR="brokers/backtesting/data/dukascopy"

# Default values
START_YEAR=${1:-2009}
START_MONTH=${2:-12}
END_YEAR=${3:-$(date +%Y)}
END_MONTH=${4:-$(date +%m)}

# Remove leading zeros for arithmetic
START_MONTH=$((10#$START_MONTH))
END_MONTH=$((10#$END_MONTH))

# Create target directory if it doesn't exist
mkdir -p "$TARGET_DIR"

# Clean up partially downloaded CSV files in download directory
echo "🧹 Cleaning up download directory..."
if [ -d "$DOWNLOAD_DIR" ]; then
    rm -f "$DOWNLOAD_DIR"/*.csv
    echo "✅ Cleaned $DOWNLOAD_DIR"
else
    mkdir -p "$DOWNLOAD_DIR"
    echo "✅ Created $DOWNLOAD_DIR"
fi
echo ""

echo "📥 Starting Dukascopy data download..."
echo "📅 Date range: ${START_YEAR}-$(printf "%02d" $START_MONTH) to ${END_YEAR}-$(printf "%02d" $END_MONTH)"
echo ""

YEAR=$START_YEAR
MONTH=$START_MONTH
TOTAL=0
SUCCESS=0
FAILED=0
SKIPPED=0

while [ $YEAR -lt $END_YEAR ] || { [ $YEAR -eq $END_YEAR ] && [ $MONTH -le $END_MONTH ]; }; do
    # Calculate last day of month
    if [ $MONTH -eq 12 ]; then
        NEXT_MONTH=1
        NEXT_YEAR=$((YEAR + 1))
    else
        NEXT_MONTH=$((MONTH + 1))
        NEXT_YEAR=$YEAR
    fi
    
    # Get last day of current month (day before first day of next month)
    LAST_DAY=$(date -d "${NEXT_YEAR}-$(printf "%02d" $NEXT_MONTH)-01 -1 day" +%d)
    
    # Format dates
    FROM_DATE=$(printf "%04d-%02d-01" $YEAR $MONTH)
    TO_DATE=$(printf "%04d-%02d-%s" $YEAR $MONTH $LAST_DAY)
    
    # Expected filename
    CSV_FILE="eurusd-tick-${FROM_DATE}-${TO_DATE}.csv"
    TARGET_FILE="$TARGET_DIR/$CSV_FILE"
    
    # Check if file already exists in target directory
    if [ -f "$TARGET_FILE" ]; then
        echo "⏭️  Skipping $FROM_DATE (already exists)"
        SKIPPED=$((SKIPPED + 1))
        TOTAL=$((TOTAL + 1))
    else
        echo "⏳ Downloading $FROM_DATE to $TO_DATE..."
        TOTAL=$((TOTAL + 1))
        
        # Run npx command
        if npx dukascopy-node -i eurusd -from "$FROM_DATE" -to "$TO_DATE" -t tick -f csv; then
            # Move downloaded file to target directory
            if [ -f "$DOWNLOAD_DIR/$CSV_FILE" ]; then
                mv "$DOWNLOAD_DIR/$CSV_FILE" "$TARGET_FILE"
                echo "✅ Downloaded and moved $CSV_FILE"
                SUCCESS=$((SUCCESS + 1))
            else
                echo "⚠️  Downloaded but file not found: $CSV_FILE"
                FAILED=$((FAILED + 1))
            fi
        else
            echo "❌ Failed to download $FROM_DATE"
            FAILED=$((FAILED + 1))
        fi
        echo ""
    fi
    
    # Move to next month
    if [ $MONTH -eq 12 ]; then
        MONTH=1
        YEAR=$((YEAR + 1))
    else
        MONTH=$((MONTH + 1))
    fi
done

echo "=========================================="
echo "✅ Download complete!"
echo "📊 Total: $TOTAL | Success: $SUCCESS | Skipped: $SKIPPED | Failed: $FAILED"
echo "=========================================="

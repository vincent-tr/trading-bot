.PHONY: convert download

# Run the data converter
convert:
	@echo "🔄 Running data converter..."
	go run ./cmd/converter

# Download Dukascopy data
download-dukascopy:
	@echo "📥 Downloading Dukascopy data..."
	./download-dukascopy.sh

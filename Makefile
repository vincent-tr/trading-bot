.PHONY: convert download-dukascopy oneshot

# Run the data converter
convert:
	@echo "🔄 Running data converter..."
	go run ./cmd/converter

# Download Dukascopy data
download-dukascopy:
	@echo "📥 Downloading Dukascopy data..."
	./download-dukascopy.sh

# Run oneshot command
oneshot:
	@echo "🚀 Running oneshot..."
	go run ./cmd/oneshot

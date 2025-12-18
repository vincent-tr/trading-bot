.PHONY: convert

# Run the data converter
convert:
	@echo "🔄 Running data converter..."
	go run ./cmd/converter

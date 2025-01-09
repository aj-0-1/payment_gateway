.PHONY: unit-test api-test start-server

start-server:
	@echo "Starting payment gateway server..."
	@go run cmd/server/main.go &
	@sleep 1

unit-test:
	@echo "Running unit tests..."
	@go test -v ./...

api-test: start-server
	@echo "\nTesting payment processing..."
	@curl -s -X POST http://localhost:8080/payments \
		-H "Content-Type: application/json" \
		-d '{"merchant_id":"merchant1","card_details":{"number":"1111111111111111","expiry_month":11,"expiry_year":2026,"cvv":"123"},"amount":{"value":100,"currency":"GBP"}}' \
		| tee /tmp/payment_response.json

	@echo "\n\nTesting payment retrieval..."
	@PAYMENT_ID=$$(jq -r '.payment_id' /tmp/payment_response.json); \
	curl -s http://localhost:8080/payments/$$PAYMENT_ID


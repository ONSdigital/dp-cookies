.PHONY: debug
debug:
	go run main.go

.PHONY: test
test: 
	go test -race -cover ./...

.PHONY: audit
audit:
	go list -m all | nancy sleuth
	
.PHONY: build	
build:
	go build ./...
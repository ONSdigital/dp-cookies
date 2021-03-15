.PHONY: debug
debug:
	go run main.go

.PHONY: test
test: 
	go test -race -cover ./...

.PHONY: audit
audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore
	
.PHONY: build	
build:
	go build ./...
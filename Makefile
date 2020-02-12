.PHONY: debug
debug:
	go run main.go

.PHONY: test
test: 
	go test -race -cover ./...
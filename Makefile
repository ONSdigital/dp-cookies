export LIBRARY_TEST:=TRUE  #Used by dp-cookies/cookies/cookies.go L37 to identify whether its running a test.

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
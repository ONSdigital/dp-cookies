.PHONY: debug
debug:
	go run main.go

.PHONY: test
test: 
	#LIBRARY_TEST=TRUE is used by dp-cookies/cookies/cookies.go @ L37 to identify whether its running in a test/locally, as we don't have the means to test secure cookies.
	LIBRARY_TEST=TRUE go test -race -cover ./...

.PHONY: audit
audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore
	
.PHONY: build	
build:
	go build ./...
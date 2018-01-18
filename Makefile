
FILES = $(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: test

gofmt:
	@gofmt -w $(FILES)
	@gofmt -r '&α{} -> new(α)' -w $(FILES)

test:
	go install ./protoc-gen-king_browser
	
	protoc --king_browser_out=.  -I . ./test/example/example.proto
	protoc --go_out=.  -I . ./test/example/example.proto
	
	protoc --go_out=.  -I . ./test/common/common.proto

.PHONY: test
test: build
	@go test -shuffle on $(shell go list ./... | grep -v sample)

.PHONY: test-verbose
test-verbose:
	go test -v -shuffle on ./...

build:
	go build

dot:
	go test ./automata/ -run TestDot
#	go test tlex.go > ex.dot; dot -Kdot -Tpng ex.dot -oex.png

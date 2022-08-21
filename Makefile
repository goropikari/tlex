dot:
	go run golex.go > ex.dot; dot -Kdot -Tpng ex.dot -oex.png

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

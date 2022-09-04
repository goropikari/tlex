build:
	../tlex -src sample.l -o main.go

test: build
	go test -shuffle on

pprof:
	go tool pprof -http=":8081" cpu.pprof

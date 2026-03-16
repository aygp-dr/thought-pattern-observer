.PHONY: build run test clean

build:
	go build -o bin/thought-pattern-observer .

run: build
	./bin/thought-pattern-observer

test:
	go test ./...

clean:
	rm -rf bin/

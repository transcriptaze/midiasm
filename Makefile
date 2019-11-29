all: test      \
	 benchmark \
     coverage

format: 
	go fmt ./...

build: format
	go build -o bin ./... 

test: build
	go clean -testcache
	go test ./...

benchmark: build
	go test -bench ./...

coverage: build
	go clean -testcache
	go test -cover ./...

clean:
	go clean
	rm -rf bin/*

debug: build
	./bin/midiasm       examples/example-01.mid
	./bin/midiasm notes examples/example-01.mid

example: build
	mkdir -p tmp
	rm -f tmp/example.out
	./bin/midiasm --out tmp/example.out examples/example-01.mid
	cat tmp/example.out

entangled: build
	./bin/midiasm       examples/entangled.mid
	./bin/midiasm notes examples/entangled.mid



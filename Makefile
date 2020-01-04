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
#	./bin/midiasm notes examples/example-01.mid

example: build
	mkdir -p tmp
	rm -f tmp/example.*
	./bin/midiasm       --debug --verbose --out tmp/example.txt examples/example-01.mid
	./bin/midiasm notes --debug --verbose --out tmp/example.notes examples/example-01.mid
	cat tmp/example.txt
	cat tmp/example.notes

split: build
	rm -f tmp/example-01.*
	./bin/midiasm --split --out tmp examples/example-01.mid
	cat tmp/example-01.MThd
	cat tmp/example-01-0.MTrk
	cat tmp/example-01-1.MTrk

entangled: build
	./bin/midiasm       examples/entangled.mid
	./bin/midiasm notes examples/entangled.mid



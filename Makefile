DIST ?= development
CMD   = ./bin/midiasm

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin/*

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

build-all: build test
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	env GOOS=linux   GOARCH=amd64 GOWORK=off go build -trimpath -o dist/$(DIST)/linux   ./...
	env GOOS=darwin  GOARCH=amd64 GOWORK=off go build -trimpath -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64 GOWORK=off go build -trimpath -o dist/$(DIST)/windows ./...

debug: build
	go test -v ./midi/encoding/midifile
	go test -v ./midi -run TestUnmarshalNoteAlias
#	$(CMD) --debug --templates examples/example-01.templates examples/example-01.mid

example: build
	mkdir -p tmp
	rm -f tmp/example.*
	$(CMD)       --debug --verbose --out tmp/example.txt examples/example-01.mid
	$(CMD) notes --debug --verbose --out tmp/example.notes examples/example-01.mid
	cat tmp/example.txt
	cat tmp/example.notes

split: build
	rm -f tmp/example-01.*
	$(CMD) --split --out tmp examples/example-01.mid
	cat tmp/example-01.MThd
	cat tmp/example-01-0.MTrk
	cat tmp/example-01-1.MTrk

entangled: build
	$(CMD)       examples/entangled.mid
	$(CMD) notes examples/entangled.mid

greensleeves: build
	$(CMD) notes --transpose +12 examples/greensleeves.mid

greensleeves2: build
	$(CMD) notes --json --transpose +12 examples/greensleeves-simple.mid \
	| jq .notes \
	| jq 'map({ note: .note, velocity: .velocity, start: .start, end: .end })' 

click: build
	$(CMD) click --debug examples/interstellar.mid

export: build
	$(CMD) export --debug examples/reference-01.mid

transpose: build
	rm -f ./tmp/greensleeves+12.mid
# 	$(CMD)           --debug --verbose examples/greensleeves.mid
	$(CMD) transpose --debug --transpose +0 -out ./tmp/greensleeves+12.mid examples/greensleeves.mid
# 	$(CMD)           --debug --verbose ./tmp/greensleeves+12.mid 
# 	diff examples/greensleeves.mid ./tmp/greensleeves+12.mid 



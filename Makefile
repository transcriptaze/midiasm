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
	mkdir -p bin
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
	mkdir -p dist/linux/$(DIST)
	mkdir -p dist/darwin/$(DIST)
	mkdir -p dist/windows/$(DIST)
	env GOOS=linux   GOARCH=amd64 GOWORK=off go build -trimpath -o dist/linux/$(DIST)   ./...
	env GOOS=darwin  GOARCH=amd64 GOWORK=off go build -trimpath -o dist/darwin/$(DIST)  ./...
	env GOOS=windows GOARCH=amd64 GOWORK=off go build -trimpath -o dist/windows/$(DIST) ./...

release: build-all
	tar --directory=dist/linux   --exclude=".DS_Store" -cvzf dist/$(DIST)-linux.tar.gz   $(DIST)
	tar --directory=dist/darwin  --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin.tar.gz  $(DIST)
	tar --directory=dist/windows --exclude=".DS_Store" -cvzf dist/$(DIST)-windows.tar.gz $(DIST)

debug: build
# 	$(CMD) --debug --verbose --C4 examples/reference-01.mid
	go test ./ops/assemble/... -run TestTextReference

delve: build
# 	dlv test github.com/transcriptaze/midiasm/midi -- run TestMTrkMarshalTrack0
	dlv debug github.com/transcriptaze/midiasm/cmd/midiasm -- assemble --debug --verbose --out tmp/example.mid examples/example.txt

help: build
	$(CMD) help
	$(CMD) help disassemble
	$(CMD) help assemble
	$(CMD) help export
	$(CMD) help notes
	$(CMD) help click
	$(CMD) help transpose

version: build
	$(CMD) version

disassemble: build
	rm -f tmp/example.txt
	$(CMD) --debug --verbose      examples/example-01.mid
	$(CMD) --debug --verbose --C4 examples/example-01.mid
	$(CMD) disassemble --debug --verbose --out tmp/example.txt examples/example-01.mid
	cat tmp/example.txt

assemble: build
	$(CMD) assemble --debug --verbose --out tmp/example.mid examples/example.txt

notes: build
	$(CMD) notes --debug --verbose --transpose +5 --out tmp/example.notes examples/example-01.mid
	cat tmp/example.notes

example: build
	mkdir -p tmp
	rm -f tmp/example.*
	$(CMD)       --debug --verbose --out tmp/example.txt examples/example-01.mid
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
	$(CMD) transpose --debug --semitones +1 -out ./tmp/greensleeves+12.mid examples/greensleeves.mid


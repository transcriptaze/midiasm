DIST ?= development_v0.3.x
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
	go build -o bin ./cmd/...

test: build
	go clean -testcache
	go test ./...

benchmark: build
	go test -bench ./...

coverage: build
	go clean -testcache
	go test -cover ./...

build-all: build test
	mkdir -p dist/$(DIST)/linux/midiasm
	mkdir -p dist/$(DIST)/arm/midiasm
	mkdir -p dist/$(DIST)/arm7/midiasm
	mkdir -p dist/$(DIST)/darwin-x64/midiasm
	mkdir -p dist/$(DIST)/darwin-arm64/midiasm
	mkdir -p dist/$(DIST)/windows/midiasm

	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux/midiasm        ./cmd/...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm/midiasm          ./cmd/...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7/midiasm         ./cmd/...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin-x64/midiasm   ./cmd/...
	env GOOS=darwin  GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin-arm64/midiasm ./cmd/...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows/midiasm      ./cmd/...

release: build-all
	tar --directory=dist/$(DIST)/linux        --exclude=".DS_Store" -cvzf dist/$(DIST)-linux.tar.gz        midiasm
	tar --directory=dist/$(DIST)/arm          --exclude=".DS_Store" -cvzf dist/$(DIST)-arm.tar.gz          midiasm
	tar --directory=dist/$(DIST)/arm7         --exclude=".DS_Store" -cvzf dist/$(DIST)-arm7.tar.gz         midiasm
	tar --directory=dist/$(DIST)/darwin-x64   --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin-x64.tar.gz   midiasm
	tar --directory=dist/$(DIST)/darwin-arm64 --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin-arm64.tar.gz midiasm
	cd dist/$(DIST)/windows; zip --recurse-paths ../../$(DIST)-windows.zip midiasm

debug: build
	go test -v ./ops/notes/... -run TestExtractNotesWithMissingNoteOff

delve: build
	dlv test github.com/transcriptaze/midiasm/ops/notes -- run TestExtractNotes
# 	dlv debug github.com/transcriptaze/midiasm/cmd/midiasm -- assemble --debug --verbose --out tmp/example.mid examples/example.txt

help: build
	$(CMD) help
	$(CMD) help commands
	$(CMD) help disassemble
	$(CMD) help assemble
	$(CMD) help export
	$(CMD) help notes
	$(CMD) help click
	$(CMD) help transpose
	$(CMD) help tsv

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
	$(CMD) transpose --debug --semitones +1 -out ./tmp/greensleeves+1.mid examples/greensleeves.mid
	$(CMD) transpose --debug --semitones +12 -out ./tmp/greensleeves+12.mid examples/greensleeves.mid

tsv: build
	rm -f ./tmp/reference.tsv
	$(CMD) tsv --debug examples/reference.mid
	$(CMD) tsv --debug --out ./tmp/reference.tsv examples/reference.mid
# 	$(CMD) tsv --debug --tabular       --out ./tmp/reference.tsv examples/reference.mid
# 	$(CMD) tsv --debug --delimiter '|' --out ./tmp/reference.tsv examples/reference.mid
	cat ./tmp/reference.tsv


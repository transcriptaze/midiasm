.PHONY: build
.PHONY: test

update:
	npm update	
	npm audit fix

build:
	npx eslint --fix *.js  
	npx eslint --fix src/**/*.js  

test: build
	npx eslint --fix test/**/*.js  
	npm test

integration-tests: test

jsdoc: build
	npx jsdoc src --destination jsdoc --package package.json --readme README.md

build-all: build

debug: build
	npx mocha 'test/**/ipx_spec.js'

disassemble: build
	node examples/disassemble.js examples/reference.mid

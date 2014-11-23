NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' ./envtmpl/envtmpl.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' ./envtmpl/envtmpl.go)
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps build

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

build:
	@mkdir -p bin/
	(cd envtmpl && go build -o ../bin/$(NAME))
	./gen_readme.sh

test: deps
	go list ./... | xargs -n1 go test -timeout=3s

xcompile: deps test
	@rm -rf build/
	@mkdir -p build
	(cd envtmpl && gox \
		-output="../build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)")

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -zcvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done

clean:
	@rm -rf bin/
	@rm -rf build/

.PHONY: all deps build test xcompile package clean

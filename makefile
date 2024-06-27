LDFLAGS += -s -w -X "main.BuildDate=$(shell date -u "+%Y-%m-%dT%H:%M:%S")"
LDFLAGS += -s -w -X "main.BuildVersion=$(shell cat version)"

msg?=
.ONESHELL:
gitcheck:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi

.ONESHELL:
pkg: gitcheck test
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	git commit -am "$(msg)"
	#jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v" && git push origin HEAD

test:
	go test -failfast -v .

init:
	go mod tidy

.PHONY: build
build: 
	go build -ldflags '$(LDFLAGS)'

install: init
	go install -ldflags='$(LDFLAGS)' .


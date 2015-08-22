VERSION = $(shell git rev-parse --short HEAD)

default: container

build:
	docker run -v $(CURDIR):/src -e LDFLAGS='-X main.version $(VERSION)' centurylink/golang-builder:latest

container: build ca-certificates.pem
	docker build -t registry.luzifer.io/rootcastore .
	docker push registry.luzifer.io/rootcastore

ca-certificates.pem:
	curl -ssLo certdata.txt https://hg.mozilla.org/mozilla-central/raw-file/tip/security/nss/lib/ckfw/builtins/certdata.txt
	extract-nss-root-certs > ca-certificates.pem
	rm certdata.txt

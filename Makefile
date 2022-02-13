GO ?= go

test:
	$(GO) test ./bart/... $(ARGS)

vet:
	$(GO) vet ./bart/...

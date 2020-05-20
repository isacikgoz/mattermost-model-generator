.PHONY: generate

generate:
	go generate
	gofmt -s -w ./output/.

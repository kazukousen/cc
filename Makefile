SRCS=$(wildcard *.go)

cc: $(SRCS)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $^

test: cc
	docker run --rm -v $(shell pwd):/cc -w /cc compilerbook bash -c \
		'./test.sh'

clean:
	rm -f cc *.o *~ tmp*

.PHONY: test clean

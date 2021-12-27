
cc: main.go codegen.go parse.go tokenize.go
	# docker run --rm -w /tmp/cc -v $(shell pwd):/tmp/cc tinygo/tinygo:0.21.0 /bin/bash -c 'tinygo build -o cc'
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cc main.go codegen.go parse.go tokenize.go

test: cc
	docker run --rm -v $(shell pwd):/cc -w /cc compilerbook bash -c \
		'./test.sh'

clean:
	rm -f cc *.o *~ tmp*

.PHONY: test clean
